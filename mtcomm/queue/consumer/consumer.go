package consumer

import (
	"fmt"
	monitor "monitor/client"
	"mtcomm/db/redis"
	logger "mtcomm/log"
	"runtime"
	"time"

	rd "github.com/gomodule/redigo/redis"
)

type Consumer interface {
	ConsumeData(queueName string, exec ExecFunc)
}

type consumer struct {
	Client redis.RedisClient
	log    *logger.Logger
	m      monitor.MonitorCaller
}

func NewConsumer(c redis.RedisClient) Consumer {
	log := logger.GetDefaultLogger()
	m := monitor.NewMonitorCaller(c)
	return &consumer{
		Client: c,
		log:    log,
		m:      m,
	}
}

type ExecFunc func(value string) error

func (c *consumer) ConsumeData(queueName string, exec ExecFunc) {
	for {
		value, err := c.Client.Lpop("queue:" + queueName)
		if err != nil {
			if err != rd.ErrNil {
				c.log.Error("msg", err)
			}
			time.Sleep(5 * time.Second)
			continue
		}
		// 处理
		doTask(c.log, queueName, exec, value, c.m)
	}
}

func doTask(log *logger.Logger, queueName string, exec ExecFunc, value string, m monitor.MonitorCaller) {
	defer func() {
		if x := recover(); x != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			msg := fmt.Sprintf("Panic Msg: %v\n%s", x, buf)

			log := logger.GetDefaultLogger()
			log.Error("panic", "true", "msg", msg)

			//monitor
			go func() {
				req := &monitor.MonitorRequest{
					ServiceName: queueName + "_consumer",
					MethodName:  "ConsumeData",
					Param:       value,
					ErrorMsg:    msg,
					ErrorTime:   time.Now().String(),
				}
				m.PushErrorLog(req)
			}()
			// 等待10秒
			time.Sleep(10 * time.Second)
		}
	}()

	err := exec(value)
	if err != nil {
		// 错误处理
		log.Error(queueName+"_consumer", "fail", "msg", err)
		//monitor
		go func() {
			req := &monitor.MonitorRequest{
				ServiceName: queueName + "_consumer",
				MethodName:  "exec",
				Param:       value,
				ErrorMsg:    err.Error(),
				ErrorTime:   time.Now().String(),
			}
			m.PushErrorLog(req)
		}()
	}
}

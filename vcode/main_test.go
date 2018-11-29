package main

import (
	"flag"
	"testing"
	//	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	//	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"mtcomm/db/redis"
	logger "mtcomm/log"

	"github.com/tjz101/goprop"
	"golang.org/x/net/context"
)

func TestCode(t *testing.T) {
	/* init log */
	logger.SetDefaultLogLevel(logger.LevelDebug)
	logger.With("serviceName", serviceName)
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop := goprop.NewProp()
	prop.Read(*propFile)
	/* init mysql */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:       context.TODO(),
		Logger:    logger.GetDefaultLogger(),
		RedisHost: prop.Get("redis_host"),
	})

	//	//init
	//	fieldKeys := []string{"method", "error"}
	//	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	//		Namespace: namespace,
	//		Subsystem: serviceName,
	//		Name:      "request_count",
	//		Help:      "Number of requests received.",
	//	}, fieldKeys)
	//	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	//		Namespace: namespace,
	//		Subsystem: serviceName,
	//		Name:      "request_latency_microseconds",
	//		Help:      "Total duration of requests in microseconds.",
	//	}, fieldKeys)

	defer redisClient.Del("vcode:18370381046:2")

	/* create service */
	var svc CodeService
	svc = codeService{}
	//	svc = loggingMiddleware{svc}
	//	svc = instrumentingMiddleware{requestCount, requestLatency, svc}
	//错误的手机号码
	err := svc.GetCode(Code{PhoneNo: "1837038104", Type: "0"})
	if err == nil {
		t.Error("手机号码校验错误")
		return
	}
	//错误的验证码类型
	err = svc.GetCode(Code{PhoneNo: "18370381046", Type: "5"})
	if err == nil {
		t.Error("验证码类型校验错误")
		return
	}
	//错误的手机号码和验证码类型
	err = svc.GetCode(Code{PhoneNo: "1837038104", Type: "5"})
	if err == nil {
		t.Error("校验错误")
		return
	}
	//取验证码
	err = svc.GetCode(Code{PhoneNo: "18370381046", Type: "2"})
	if err != nil {
		t.Error("获取验证码错误")
		return
	}
	//57s内重复取验证码
	err = svc.GetCode(Code{PhoneNo: "18370381046", Type: "2"})
	if err == nil {
		t.Error("57s内重复获取证码错误")
		return
	}
	vcode, _ := redisClient.Get("vcode:18370381046:2")
	//校验验证码手机号码错误
	err = svc.CheckCode(Code{PhoneNo: "1837038104", Type: "2", Vcode: vcode})
	if err == nil {
		t.Error("手机号码错误校验失败")
		return
	}
	//校验验证码类型错误
	err = svc.CheckCode(Code{PhoneNo: "18370381046", Type: "1", Vcode: vcode})
	if err == nil {
		t.Error("错误类型校验失败")
		return
	}
	//验证码错误
	err = svc.CheckCode(Code{PhoneNo: "18370381046", Type: "2", Vcode: "0000"})
	if err == nil {
		t.Error("验证码校验失败")
		return
	}
	//正确的手机号码类型验证码
	err = svc.CheckCode(Code{PhoneNo: "18370381046", Type: "2", Vcode: vcode})
	if err != nil {
		t.Error("验证码校验失败")
		return
	}
}

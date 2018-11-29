package main

import (
	"context"
	"encoding/json"
	"flag"
	mail "mail/client"
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
	"mtcomm/k8s"
	logger "mtcomm/log"
	csr "mtcomm/queue/consumer"
	"os"
	sms "sms/client"
	"strconv"

	"github.com/bluele/gcache"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/tjz101/goprop"
)

const (
	listenPort  = ":8888"
	serviceName = "monitor"
)

var (
	namespace   string
	tracer      stdopentracing.Tracer
	prop        *goprop.Prop
	log         *logger.Logger
	svc         MonitorService
	k8sClient   k8s.K8sClient
	cacheClient gcache.Cache
	redisClient redis.RedisClient
	mysqlClient mysql.MysqlClient
	mailClient  mail.MailCaller
	smsClient   sms.SmsCaller
)

func init() {
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop = goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace

	/* init log */
	logLevel, _ := strconv.Atoi(prop.Get("log_level"))
	logger.SetDefaultLogLevel(logLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	/* init k8s */
	k8sClient = k8s.NewK8sClient()

	/* init gcache */
	cacheClient = gcache.New(10).LRU().Build()

	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	/* init mysql */
	max, _ := strconv.Atoi(prop.Get("mysql_mail_maxidleconn"))
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     prop.Get("mysql_monitor_username"),
		Password:     prop.Get("mysql_monitor_password"),
		IP:           prop.Get("mysql_monitor_host"),
		Port:         prop.Get("mysql_monitor_port"),
		DatabaseName: prop.Get("mysql_db_monitor"),
		Logger:       logger.GetDefaultLogger(),
		MaxIdleConns: max,
	})

	/* init mail */
	mailClient = mail.NewMailCaller(redisClient)

	/* init sms */
	smsClient = sms.NewSmsCaller(redisClient)
}

func main() {
	// init
	zipkinAddr := prop.Get("zipkinAddr")

	// init tracing domain.
	{
		if zipkinAddr != "" {
			log.Info("tracer", "Zipkin", "zipkinAddr", zipkinAddr)
			collector, err := zipkin.NewHTTPCollector(zipkinAddr, zipkin.HTTPBatchSize(1))
			if err != nil {
				log.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, listenPort, serviceName),
			)
			if err != nil {
				log.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
		} else {
			log.Info("tracer", "none")
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	/* create service */
	svc = monitorService{}
	svc = loggingMiddleware{svc}

	callback := func(data string) error {
		//service
		req := &monitorRequest{}
		err := json.Unmarshal([]byte(data), req)
		if err != nil {
			return err
		}
		svc.monitor(req)
		return nil
	}

	//create consumer
	cc := csr.NewConsumer(redisClient)
	cc.ConsumeData("monitor", callback)

	done := make(chan struct{})
	<-done
}

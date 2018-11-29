package main

import (
	"context"
	"encoding/json"
	"flag"
	idgen "idgen/client"
	monitor "monitor/client"
	http3 "mtcomm/caller/http3part"
	"mtcomm/db/redis"
	"mtcomm/k8s"
	logger "mtcomm/log"
	csr "mtcomm/queue/consumer"
	"os"
	"strconv"

	"github.com/bluele/gcache"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/tjz101/goprop"
)

const (
	listenPort  = ":8888"
	serviceName = "sms"
)

var (
	namespace     string
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	svc           SmsService
	k8sClient     k8s.K8sClient
	idClient      idgen.IdGenClient
	cacheClient   gcache.Cache
	monitorClient monitor.MonitorCaller
	http3Client   http3.CallProxy
	redisClient   redis.RedisClient
)

func init() {
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop = goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace

	/* init log */
	LogLevel, _ := strconv.Atoi(prop.Get("LogLevel"))
	logger.SetDefaultLogLevel(LogLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	/* init k8s */
	k8sClient = k8s.NewK8sClient()

	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	/* init local cache */
	cacheClient = gcache.New(10).LRU().Build()

	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)
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
	svc = smsService{}
	svc = loggingMiddleware{svc}

	callback := func(data string) error {
		//service
		param := &SmsInfo{}
		err := json.Unmarshal([]byte(data), param)
		if err != nil {
			return err
		}
		_, err = svc.Sms(param)
		if err != nil {
			return err
		}
		return nil
	}

	//create consumer
	cc := csr.NewConsumer(redisClient)
	cc.ConsumeData("sms", callback)

	done := make(chan struct{})
	<-done
}

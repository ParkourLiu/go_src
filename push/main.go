package main

import (
	"context"
	"encoding/json"
	"flag"
	"mtcomm/db/redis"
	"mtcomm/k8s"
	logger "mtcomm/log"
	csr "mtcomm/queue/consumer"
	"os"

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
	svc         PushService
	k8sClient   k8s.K8sClient
	cacheClient gcache.Cache
	redisClient redis.RedisClient

	appid         string
	appkey        string
	mastersecret  string
	appsecret     string
	authtoken_URL string
	List_push_URL string
	Push_All_Url  string
)

func init() {
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop = goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace

	appid = prop.Get("push_appid")
	appkey = prop.Get("push_appkey")
	mastersecret = prop.Get("push_mastersecret")
	appsecret = prop.Get("push_appsecret")
	authtoken_URL = "https://restapi.getui.com/v1/" + appid + "/auth_sign"
	List_push_URL = "https://restapi.getui.com/v1/" + appid + "/push_single"
	Push_All_Url = "https://restapi.getui.com/v1/" + appid + "/push_app"
	/* init log */
	logger.SetDefaultLogLevel(logger.LevelDebug)
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
	svc = pushService{}
	svc = loggingMiddleware{svc}

	callback := func(data string) error {
		//service
		param := &pushInfoList{}
		err := json.Unmarshal([]byte(data), param)
		if err != nil {
			return err
		}

		_, err = svc.Push(param)
		if err != nil {
			return err
		}
		return nil
	}

	//create consumer
	cc := csr.NewConsumer(redisClient)
	cc.ConsumeData("push", callback)

	done := make(chan struct{})
	<-done
}

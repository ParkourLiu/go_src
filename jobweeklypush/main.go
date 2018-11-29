package main

import (
	"flag"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	"os"
	recommend "recommend/client"
	"strconv"
	//"time"
	push "push/client"

	"context"
	"mtcomm/db/redis"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/tjz101/goprop"
)

const (
	listenPort  = ":8888"
	serviceName = "jobweeklypush"
)

var (
	namespace       string
	tracer          stdopentracing.Tracer
	prop            *goprop.Prop
	log             *logger.Logger
	recommendClient recommend.RecommendClient
	k8sClient       k8s.K8sClient
	pushClient      push.PushCaller
	redisClient     redis.RedisClient
)

func init() {
	/* init k8s */
	k8sClient = k8s.NewK8sClient()

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

	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	//初始化推送
	pushClient = push.NewPushCaller(redisClient)
}
func main() {
	log.Debug("▃█▃█▃█▃█▃█▃█savehotdynamic█▃█▃█▃█▃█▃█▃", "开始")
	// init
	zipkinAddr := prop.Get("zipkinAddr")

	// init tracing domain.
	{
		if zipkinAddr != "" {
			logger.Info("tracer", "Zipkin", "zipkinAddr", zipkinAddr)
			collector, err := zipkin.NewHTTPCollector(zipkinAddr, zipkin.HTTPBatchSize(1))
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, listenPort, serviceName),
			)
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
		} else {
			logger.Info("tracer", "none")
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}
	list := []push.PushInfo{}
	pushInfo0 := push.PushInfo{
		Title:    "不少有趣的人分享了动态，快去打开看看",
		Text:     " ",
		JsonInfo: "{\"alert\": \"不少有趣的人分享了动态，快去打开看看\",\"extras\": {\"type\": \"PA1\"}}",
		Alias:    " ",
	}
	list = append(list, pushInfo0)
	req := &push.PushRequest{
		PushType:  "all_0", //0代表推送消息，1代表透传消息,all_0代表给所有用户推送通知
		PushInfos: list,
	}
	err := pushClient.PushErrorLog(req)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	log.Debug("▃█▃█▃█▃█▃█▃█PushFavour█▃█▃█▃█▃█▃█▃", "结束")
}

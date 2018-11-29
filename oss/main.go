package main

import (
	"flag"
	k8s "mtcomm/k8s"
	"net/http"
	"os"
	"strconv"

	monitor "monitor/client"
	"mtcomm/helper/health"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"

	"github.com/afex/hystrix-go/hystrix"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tjz101/goprop"

	"encoding/base64"

	idgen "idgen/client"
	"mtcomm/db/redis"

	"context"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	listenPort = ":8888"

	serviceName = "oss"
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var (
	namespace     string
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	k8sClient     k8s.K8sClient
	monitorClient monitor.MonitorCaller
	redisClient   redis.RedisClient
	idGenClient   idgen.IdGenClient

	//oss
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string

	//ossInfo
	coder           *base64.Encoding
	accessKeyId     string
	accessKeySecret string
	host            string
	expire_time     int64
	upload_dir      string
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
	//OSS
	AccessKeyID = prop.Get("aliyun_oss_access_key_id")
	AccessKeySecret = prop.Get("aliyun_oss_access_key_secret")
	Endpoint = prop.Get("aliyun_oss_endpoint")
	BucketName = prop.Get("aliyun_oss_bucket_name_mtalk")
	//OssInfo
	coder = base64.NewEncoding(base64Table)
	accessKeyId = AccessKeyID
	accessKeySecret = AccessKeySecret
	host = "http://" + BucketName + "." + Endpoint
	expire_time = 60
	upload_dir = ""

	/* init log */
	LogLevel, _ := strconv.Atoi(prop.Get("LogLevel"))
	logger.SetDefaultLogLevel(LogLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()
	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})
	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)
}
func main() {
	// init
	zipkinAddr := prop.Get("zipkinAddr")

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(*logger.GetDefaultLogger().GetDefaultKitLogger()),
	}

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

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: serviceName,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: serviceName,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	//初始化id生成器
	idGenClient = idgen.NewIdGenClient(k8sClient, tracer, namespace)
	/* create service */
	var svc Oss
	svc = oss{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	//1
	GetOssTokenForWebHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetOssTokenForWeb")(makeGetOssTokenForWebEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetOssTokenForWeb", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetOssTokenForWebHandler := mkhttpserver.NewHttpMtalkServer(GetOssTokenForWebHandler, serviceName, "GetOssTokenForWeb", monitorClient, &Response{Err: "panic"})
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/getOssTokenForWeb", mtalkGetOssTokenForWebHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

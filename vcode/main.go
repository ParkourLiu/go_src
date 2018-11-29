package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	sms "sms/client"
	UserClient "user/client"

	"github.com/afex/hystrix-go/hystrix"

	"golang.org/x/net/context"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tjz101/goprop"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	listenPort  = ":8888"
	serviceName = "vcode"

	//type值错误
	flag1 = serviceName + "101"

	//该手机号已注册
	flag2 = serviceName + "102"

	//该用户不存在
	flag3 = serviceName + "103"

	///该手机号与用户绑定手机号不一致
	flag4 = serviceName + "104"

	//请间隔1分钟再试
	flag5 = serviceName + "105"

	//验证码错误
	flag6 = serviceName + "106"
)

var (
	namespace     string
	redisClient   redis.RedisClient
	tracer        stdopentracing.Tracer
	log           *logger.Logger
	smsClient     sms.SmsCaller
	userClient    UserClient.UserClient
	k8sClient     k8s.K8sClient
	monitorClient monitor.MonitorCaller
)

func main() {
	/* init k8s */
	k8sClient = k8s.NewK8sClient()
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop := goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace
	zipkinAddr := prop.Get("zipkinAddr")
	/* init log */
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

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(*logger.GetDefaultLogger().GetDefaultKitLogger()),
	}

	//初始化短信模块
	smsClient = sms.NewSmsCaller(redisClient)

	// Tracing domain.

	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)

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
	//初始化user
	userClient = UserClient.NewUserClient(k8sClient, tracer, namespace, serviceName, "user", "8888")

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

	/* create service */
	var svc CodeService
	svc = codeService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	validateAuthCodeHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "validateAuthCode")(validateAuthCodeEndpoint(svc)),
		decodeValidateAuthCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "validateAuthCode", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkvalidateAuthCodeHandler := mkhttpserver.NewHttpMtalkServer(validateAuthCodeHandler, serviceName, "validateAuthCode", monitorClient, &validateAuthCodeResponse{Err: "panic"})

	getAuthCodeHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "getAuthCode")(getAuthCodeEndpoint(svc)),
		decodeGetCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "getAuthCode", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkgetAuthCodeHandler := mkhttpserver.NewHttpMtalkServer(getAuthCodeHandler, serviceName, "getAuthCode", monitorClient, &getCodeResponse{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/getAuthCode", mtalkgetAuthCodeHandler)
	http.Handle("/validateAuthCode", mtalkvalidateAuthCodeHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{RedisClient: redisClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

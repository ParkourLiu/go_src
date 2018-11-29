package main

import (
	"context"
	"flag"
	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	"net/http"
	"os"
	push "push/client"
	"strconv"
	UserClient "user/client"

	"github.com/afex/hystrix-go/hystrix"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tjz101/goprop"
)

const (
	listenPort  = ":8888"
	serviceName = "token"
)

var (
	namespace     string
	redisClient   redis.RedisClient
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	monitorClient monitor.MonitorCaller
	userClient    UserClient.UserClient
	pushClient    push.PushCaller
	k8sClient     k8s.K8sClient
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
	logLevel, _ := strconv.Atoi(prop.Get("log_level"))
	logger.SetDefaultLogLevel(logLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	/* init reids */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)

	pushClient = push.NewPushCaller(redisClient)
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

	//初始化user
	userClient = UserClient.NewUserClient(k8sClient, tracer, namespace, serviceName, "user", "8888")
	/* create service */
	var svc TokenService
	svc = tokenService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	createTokenHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "CreateToken")(makeCreateTokenEndpoint(svc)),
		decodeCreateTokenRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CreateToken", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)

	deleteTokenHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "DeleteToken")(makeDeleteTokenEndpoint(svc)),
		decodeDeleteTokenRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "DeleteToken", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkCreateTokenHandler := mkhttpserver.NewHttpMtalkServer(createTokenHandler, serviceName, "CreateToken", monitorClient, &CreateTokenResponse{Token: "", Err: "panic", Code: ""})
	mtalkDeleteTokenHandler := mkhttpserver.NewHttpMtalkServer(deleteTokenHandler, serviceName, "DeleteToken", monitorClient, &DeleteTokenResponse{Err: "panic", Code: ""})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/createToken", mtalkCreateTokenHandler)
	http.Handle("/deleteToken", mtalkDeleteTokenHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{RedisClient: redisClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

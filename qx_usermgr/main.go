package main

import (
	"context"
	"flag"
	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"
	"mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	"net/http"
	"os"
	"strconv"
	UserClient "qx_user/client"
	"github.com/afex/hystrix-go/hystrix"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tjz101/goprop"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	//登录方式
	loginFlag_Phone  = "0"
	loginFlag_QQ     = "1"
	loginFlag_Wechat = "2"
	loginFlag_Sina   = "3"

	listenPort = ":8888"

	serviceName = "qx_usermgr"

	code1 = serviceName + "101" //该手机号未注册
	code2 = serviceName + "102" //该平台下此用户已绑定过了
	code3 = serviceName + "103" //不可解绑以当前平台登录的平台
)

var (
	namespace     string
	redisClient   redis.RedisClient
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	userClient    UserClient.UserClient
	k8sClient     k8s.K8sClient
	monitorClient monitor.MonitorCaller
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
	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)

	//初始化user
	userClient = UserClient.NewUserClient(k8sClient, tracer, namespace)

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

	/* create service */
	var svc UserService
	svc = userService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	RegAndLoginHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "RegAndLogin")(makeRegAndLoginEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "RegAndLogin", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkRegAndLoginHandler := mkhttpserver.NewHttpMtalkServer(RegAndLoginHandler, serviceName, "RegAndLogin", monitorClient, &Response{Err: "panic"})

	OtherLoginHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "OtherLogin")(makeOtherLoginEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "OtherLogin", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkOtherLoginHandler := mkhttpserver.NewHttpMtalkServer(OtherLoginHandler, serviceName, "OtherLogin", monitorClient, &Response{Err: "panic"})

	SearchUserHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SearchUser")(makeSearchUserEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SearchUser", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSearchUserHandler := mkhttpserver.NewHttpMtalkServer(SearchUserHandler, serviceName, "SearchUser", monitorClient, &Response{Err: "panic"})

	UpdateUserHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateUser")(makeUpdateUserEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateUser", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkUpdateUserHandler := mkhttpserver.NewHttpMtalkServer(UpdateUserHandler, serviceName, "UpdateUser", monitorClient, &Response{Err: "panic"})

	ChangeBindHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ChangeBind")(makeChangeBindEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ChangeBind", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkChangeBindHandler := mkhttpserver.NewHttpMtalkServer(ChangeBindHandler, serviceName, "ChangeBind", monitorClient, &Response{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/regAndLogin", mtalkRegAndLoginHandler)
	http.Handle("/otherLogin", mtalkOtherLoginHandler)
	http.Handle("/searchUser", mtalkSearchUserHandler)
	http.Handle("/updateUser", mtalkUpdateUserHandler)
	http.Handle("/changeBind", mtalkChangeBindHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{RedisClient: redisClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

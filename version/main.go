package main

import (
	"context"
	"flag"
	idgen "idgen/client"
	monitor "monitor/client"
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	"mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	"net/http"
	"os"
	"strconv"

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
	listenPort = ":8888"

	serviceName = "version"

	//有强制更新
	flag1 = "6666"

	//有更新
	flag2 = "9999"
)

var (
	namespace     string
	mysqlClient   mysql.MysqlClient
	redisClient   redis.RedisClient
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	idGenClient   idgen.IdGenClient
	monitorClient monitor.MonitorCaller
	k8sClient     k8s.K8sClient
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

	/* init mysql */
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     prop.Get("mysql_mtalk_username"),
		Password:     prop.Get("mysql_mtalk_password"),
		IP:           prop.Get("mysql_mtalk_host"),
		Port:         prop.Get("mysql_mtalk_port"),
		DatabaseName: prop.Get("mysql_db_mtalk"),
		Logger:       logger.GetDefaultLogger(),
	})
	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	k8sClient = k8s.NewK8sClient()

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

	/* create service */
	var svc VersionService
	svc = versionService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	VersionInfoHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "VersionInfo")(makeVersionInfoEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "VersionInfo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkVersionInfoHandler := mkhttpserver.NewHttpMtalkServer(VersionInfoHandler, serviceName, "VersionInfo", monitorClient, &Response{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)

	http.Handle("/versionInfo", mtalkVersionInfoHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

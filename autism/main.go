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

	serviceName = "autism"
	flag1       = serviceName + "01" //星星id不存在
	appid       = "wx70f1ac2a93a63410"
	secret      = "eae728f30c360581237fc8f0e7e6ddf0"
	grant_type  = "authorization_code"
	//	appid = wx70f1ac2a93a63410
	//	secret = eae728f30c360581237fc8f0e7e6ddf0
	//grant_type = authorization_code
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

	//========================
	//weixin
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
		UserName:     prop.Get("mysql_autism_username"),
		Password:     prop.Get("mysql_autism_password"),
		IP:           prop.Get("mysql_autism_host"),
		Port:         prop.Get("mysql_autism_port"),
		DatabaseName: prop.Get("mysql_db_autism"),
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

	//初始化id生成器
	idGenClient = idgen.NewIdGenClient(k8sClient, tracer, namespace)
	/* create service */
	var svc AutismService
	svc = autismService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	StarDetailsHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "StarDetails")(makeStarDetailsEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "StarDetails", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkStarDetailsHandler := mkhttpserver.NewHttpMtalkServer(StarDetailsHandler, serviceName, "StarDetails", monitorClient, &StarDetailsResponse{Err: "panic"})

	SaveCommentHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SaveComment")(makeSaveCommentEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SaveComment", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSaveCommentHandler := mkhttpserver.NewHttpMtalkServer(SaveCommentHandler, serviceName, "SaveComment", monitorClient, &StarDetailsResponse{Err: "panic"})
	//==================================================================================================================
	StarListHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "StarList")(autismStarListEndpoint(svc)),
		decodeStarListRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "StarList", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	autismStarListHandler := mkhttpserver.NewHttpMtalkServer(StarListHandler, serviceName, "StarList", monitorClient, &autismResponse{Err: "panic"})

	LikesHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "Likes")(autismLikesEndpoint(svc)),
		decodeLikesRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "Likes", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	autismLikesHandler := mkhttpserver.NewHttpMtalkServer(LikesHandler, serviceName, "Likes", monitorClient, &autismResponse{Err: "panic"})

	GetUnionidHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetUnionid")(autismGetUnionidEndpoint(svc)),
		decodeGetUnionidRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetUnionid", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	autismGetUnionidHandler := mkhttpserver.NewHttpMtalkServer(GetUnionidHandler, serviceName, "GetUnionid", monitorClient, &autismResponse{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/starDetails", mtalkStarDetailsHandler)
	http.Handle("/saveComment", mtalkSaveCommentHandler)

	http.Handle("/getUnionid", autismGetUnionidHandler)
	http.Handle("/starList", autismStarListHandler)
	http.Handle("/likes", autismLikesHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

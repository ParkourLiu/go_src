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

	serviceName = "mschool"

	//该学校本学年已创建过
	code1 = serviceName + "101"
	//无权操作
	code2 = serviceName + "102"
	//id不存在
	code3 = serviceName + "103"
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
	InviteUrl     string

	//oss
	key_id            string
	key_secret        string
	bucket_name_mtalk string
	bucket_name_face  string
	oss_endpoint      string
	domain_app        string
)

func init() {
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop = goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace
	//oss
	key_id = prop.Get("aliyun_oss_access_key_id")
	key_secret = prop.Get("aliyun_oss_access_key_secret")
	bucket_name_mtalk = prop.Get("aliyun_oss_bucket_name_mtalk")
	bucket_name_face = prop.Get("aliyun_oss_bucket_name_face")
	oss_endpoint = prop.Get("aliyun_oss_endpoint")
	domain_app = prop.Get("domain_app")
	InviteUrl = prop.Get("InviteUrl")
	/* init log */
	LogLevel, _ := strconv.Atoi(prop.Get("LogLevel"))
	logger.SetDefaultLogLevel(LogLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	/* init mysql */
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     prop.Get("mysql_mschool_username"),
		Password:     prop.Get("mysql_mschool_password"),
		IP:           prop.Get("mysql_mschool_host"),
		Port:         prop.Get("mysql_mschool_port"),
		DatabaseName: prop.Get("mysql_db_mschool"),
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
	var svc MschoolService
	svc = mschoolService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	CreateSchoolYearHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "CreateSchoolYear")(makeCreateSchoolYearEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CreateSchoolYear", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkCreateSchoolYear := mkhttpserver.NewHttpMtalkServer(CreateSchoolYearHandler, serviceName, "CreateSchoolYear", monitorClient, &Response{Err: "panic"})

	SearchSchoolHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SearchSchool")(makeSearchSchoolEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SearchSchool", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSearchSchool := mkhttpserver.NewHttpMtalkServer(SearchSchoolHandler, serviceName, "SearchSchool", monitorClient, &Response{Err: "panic"})

	MySchoolHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "MySchool")(makeMySchoolEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "MySchool", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkMySchool := mkhttpserver.NewHttpMtalkServer(MySchoolHandler, serviceName, "MySchool", monitorClient, &Response{Err: "panic"})

	SetWorkDayHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SetWorkDay")(makeSetWorkDayEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SetWorkDay", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSetWorkDay := mkhttpserver.NewHttpMtalkServer(SetWorkDayHandler, serviceName, "SetWorkDay", monitorClient, &Response{Err: "panic"})

	LookWorkDayHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "LookWorkDay")(makeLookWorkDayEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "LookWorkDay", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkLookWorkDay := mkhttpserver.NewHttpMtalkServer(LookWorkDayHandler, serviceName, "LookWorkDay", monitorClient, &Response{Err: "panic"})

	WorkDayHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "WorkDay")(makeWorkDayEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "WorkDay", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkWorkDay := mkhttpserver.NewHttpMtalkServer(WorkDayHandler, serviceName, "WorkDay", monitorClient, &Response{Err: "panic"})

	UpSchoolYearHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpSchoolYear")(makeUpSchoolYearEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpSchoolYear", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkUpSchoolYear := mkhttpserver.NewHttpMtalkServer(UpSchoolYearHandler, serviceName, "UpSchoolYear", monitorClient, &Response{Err: "panic"})

	FaceDataGatherHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "FaceDataGather")(makeFaceDataGatherEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "FaceDataGather", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkFaceDataGather := mkhttpserver.NewHttpMtalkServer(FaceDataGatherHandler, serviceName, "FaceDataGather", monitorClient, &Response{Err: "panic"})

	LabelGuidDataGatherHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "LabelGuidDataGather")(makeLabelGuidDataGatherEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "LabelGuidDataGather", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkLabelGuidDataGather := mkhttpserver.NewHttpMtalkServer(LabelGuidDataGatherHandler, serviceName, "LabelGuidDataGather", monitorClient, &Response{Err: "panic"})

	GetDataFileUrlHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetDataFileUrl")(makeGetDataFileUrlEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetDataFileUrl", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetDataFileUrl := mkhttpserver.NewHttpMtalkServer(GetDataFileUrlHandler, serviceName, "GetDataFileUrl", monitorClient, &Response{Err: "panic"})

	DelDataFileHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "DelDataFile")(makeDelDataFileEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "DelDataFile", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkDelDataFile := mkhttpserver.NewHttpMtalkServer(DelDataFileHandler, serviceName, "DelDataFile", monitorClient, &Response{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/createSchoolYear", mtalkCreateSchoolYear)
	http.Handle("/searchSchool", mtalkSearchSchool)
	http.Handle("/mySchool", mtalkMySchool)
	http.Handle("/setWorkDay", mtalkSetWorkDay)
	http.Handle("/lookWorkDay", mtalkLookWorkDay)
	http.Handle("/workDay", mtalkWorkDay)
	http.Handle("/upSchoolYear", mtalkUpSchoolYear)

	http.Handle("/faceDataGather", mtalkFaceDataGather)
	http.Handle("/labelGuidDataGather", mtalkLabelGuidDataGather)
	http.Handle("/getDataFileUrl", mtalkGetDataFileUrl)
	http.Handle("/delDataFile", mtalkDelDataFile)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

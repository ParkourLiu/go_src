package main

import (
	activity "activity/client"
	"context"
	"flag"
	grpdynamic "grpdynamic/client"
	mactivity "mactivity/client"
	mgrpmgr "mgrpmgr/client"
	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"
	"mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	"net/http"
	"os"
	push "push/client"
	"strconv"
	user "user/client"

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
	listenPort = ":8888"

	serviceName = "recommend"

	///该id不存在
	flag1 = serviceName + "101"

	fanUId = "UFans"
)

var (
	namespace        string
	redisClient      redis.RedisClient
	tracer           stdopentracing.Tracer
	prop             *goprop.Prop
	log              *logger.Logger
	k8sClient        k8s.K8sClient
	monitorClient    monitor.MonitorCaller
	grpdynamicClient grpdynamic.DynamicClient
	userClient       user.UserClient
	pushClient       push.PushCaller
	activityClient   activity.ActivityClient
	mgrpmgrClient    mgrpmgr.MgrpClient
	mactivityClient  mactivity.MactivityClient
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
	//初始化推送
	pushClient = push.NewPushCaller(redisClient)
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

	//初始化动态
	grpdynamicClient = grpdynamic.NewDynamicClient(k8sClient, tracer, namespace, serviceName, "grpdynamic", "8888")
	//初始化用户
	userClient = user.NewUserClient(k8sClient, tracer, namespace, serviceName, "user", "8888")
	//活动
	activityClient = activity.NewActivityClient(k8sClient, tracer, namespace, serviceName, "activity", "8888")
	//社群服务层
	mgrpmgrClient = mgrpmgr.NewMgrpClient(k8sClient, tracer, namespace, serviceName, "mgrpmgr", "8888")
	//活动服务层
	mactivityClient = mactivity.NewMactivityClient(k8sClient, tracer, namespace, serviceName, "mactivity", "8888")

	/* create service */
	var svc Recommend
	svc = recommend{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}
	//1
	PopuserUsersHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "PopuserUsers")(makePopuserUsersEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "PopuserUsers", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkPopuserUsersHandler := mkhttpserver.NewHttpMtalkServer(PopuserUsersHandler, serviceName, "PopuserUsers", monitorClient, &Response{Err: "panic"})
	//2
	RmduserUsersHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "RmduserUsers")(makeRmduserUsersEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "RmduserUsers", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkRmduserUsersHandler := mkhttpserver.NewHttpMtalkServer(RmduserUsersHandler, serviceName, "RmduserUsers", monitorClient, &Response{Err: "panic"})
	//3
	HotDynamicHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "HotDynamic")(makeHotDynamicEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "HotDynamic", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkHotDynamicHandler := mkhttpserver.NewHttpMtalkServer(HotDynamicHandler, serviceName, "HotDynamic", monitorClient, &Response{Err: "panic"})

	//4
	NewDynamicHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "NewDynamic")(makeNewDynamicEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "NewDynamic", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkNewDynamicHandler := mkhttpserver.NewHttpMtalkServer(NewDynamicHandler, serviceName, "NewDynamic", monitorClient, &Response{Err: "panic"})

	//3
	ReverseHotDynamicHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ReverseHotDynamic")(makeReverseHotDynamicEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ReverseHotDynamic", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkReverseHotDynamicHandler := mkhttpserver.NewHttpMtalkServer(ReverseHotDynamicHandler, serviceName, "ReverseHotDynamic", monitorClient, &Response{Err: "panic"})

	//4
	ReverseNewDynamicHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ReverseNewDynamic")(makeReverseNewDynamicEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ReverseNewDynamic", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkReverseNewDynamicHandler := mkhttpserver.NewHttpMtalkServer(ReverseNewDynamicHandler, serviceName, "ReverseNewDynamic", monitorClient, &Response{Err: "panic"})

	SavePopuserUsersHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SavePopuserUsers")(makeSavePopuserUsersEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SavePopuserUsers", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSavePopuserUsersHandler := mkhttpserver.NewHttpMtalkServer(SavePopuserUsersHandler, serviceName, "SavePopuserUsers", monitorClient, &Response{Err: "panic"})

	SaveRmduserUsersHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SaveRmduserUsers")(makeSaveRmduserUsersEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SaveRmduserUsers", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSaveRmduserUsersHandler := mkhttpserver.NewHttpMtalkServer(SaveRmduserUsersHandler, serviceName, "SaveRmduserUsers", monitorClient, &Response{Err: "panic"})

	SaveHotDynamicHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SaveHotDynamic")(makeSaveHotDynamicEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SaveHotDynamic", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSaveSaveHotDynamicHandler := mkhttpserver.NewHttpMtalkServer(SaveHotDynamicHandler, serviceName, "SaveHotDynamic", monitorClient, &Response{Err: "panic"})

	FriendRecommendHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "FriendRecommend")(makeFriendRecommendEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "FriendRecommend", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkFriendRecommendHandler := mkhttpserver.NewHttpMtalkServer(FriendRecommendHandler, serviceName, "FriendRecommend", monitorClient, &Response{Err: "panic"})

	SaveFriendRecommendHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SaveFriendRecommend")(makeSaveFriendRecommendEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SaveFriendRecommend", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSaveFriendRecommendHandler := mkhttpserver.NewHttpMtalkServer(SaveFriendRecommendHandler, serviceName, "SaveFriendRecommend", monitorClient, &Response{Err: "panic"})

	SearchRecommendHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SearchRecommend")(makeSearchRecommendEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SearchRecommend", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSearchRecommendHandler := mkhttpserver.NewHttpMtalkServer(SearchRecommendHandler, serviceName, "SearchRecommend", monitorClient, &Response{Err: "panic"})

	PushFavourHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "PushFavour")(makePushFavourEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "PushFavour", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkPushFavourHandler := mkhttpserver.NewHttpMtalkServer(PushFavourHandler, serviceName, "PushFavour", monitorClient, &Response{Err: "panic"})

	PushStartActivityHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "PushStartActivity")(makePushStartActivityEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "PushStartActivity", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkPushStartActivityHandler := mkhttpserver.NewHttpMtalkServer(PushStartActivityHandler, serviceName, "PushStartActivity", monitorClient, &Response{Err: "panic"})

	SaveHomePageCacheHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SaveHomePageCache")(makeSaveHomePageCacheEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SaveHomePageCache", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSaveHomePageCacheHandler := mkhttpserver.NewHttpMtalkServer(SaveHomePageCacheHandler, serviceName, "SaveHomePageCache", monitorClient, &Response{Err: "panic"})

	AddDynamicFansHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "AddDynamicFans")(makeAddDynamicFansEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "AddDynamicFans", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkAddDynamicFansHandler := mkhttpserver.NewHttpMtalkServer(AddDynamicFansHandler, serviceName, "AddDynamicFans", monitorClient, &Response{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/popuserUsers", mtalkPopuserUsersHandler)
	http.Handle("/rmduserUsers", mtalkRmduserUsersHandler)
	http.Handle("/hotDynamic", mtalkHotDynamicHandler)
	http.Handle("/newDynamic", mtalkNewDynamicHandler)
	http.Handle("/reverseHotDynamic", mtalkReverseHotDynamicHandler)
	http.Handle("/reverseNewDynamic", mtalkReverseNewDynamicHandler)
	http.Handle("/friendRecommend", mtalkFriendRecommendHandler)
	http.Handle("/searchRecommend", mtalkSearchRecommendHandler)

	http.Handle("/savePopuserUsers", mtalkSavePopuserUsersHandler)
	http.Handle("/saveRmduserUsers", mtalkSaveRmduserUsersHandler)
	http.Handle("/saveHotDynamic", mtalkSaveSaveHotDynamicHandler)
	http.Handle("/saveFriendRecommend", mtalkSaveFriendRecommendHandler)
	http.Handle("/saveHomePageCache", mtalkSaveHomePageCacheHandler)

	http.Handle("/pushFavour", mtalkPushFavourHandler)
	http.Handle("/pushStartActivity", mtalkPushStartActivityHandler)

	http.Handle("/addDynamicFans", mtalkAddDynamicFansHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{RedisClient: redisClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

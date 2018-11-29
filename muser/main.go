package main

import (
	"context"
	"flag"
	idgen "idgen/client"
	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	"net/http"
	"os"
	search "search/client"
	"strconv"
	UserClient "user/client"
	vcode "vcode/client"

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

	serviceName = "muser"

	///账号错误
	flag1 = serviceName + "101"

	//账号不可以为空
	flag2 = serviceName + "102"

	//数据异常
	flag3 = serviceName + "103"

	//OtherLogienType值错误
	flag4 = serviceName + "104"

	//该手机号已注册
	flag5 = serviceName + "105"

	//该平台下此用户已绑定过了
	flag6 = serviceName + "106"

	//不可解绑登录状态的平台
	flag7 = serviceName + "107"

	//非法用户名
	flag8 = serviceName + "108"

	QinJiaName = "钦家,钦家官方" //此处字段代表注册用户不可以使用此字段中名字，新增请用英文逗号隔开即可实现
)

var (
	namespace     string
	redisClient   redis.RedisClient
	tracer        stdopentracing.Tracer
	idGenClient   idgen.IdGenClient
	prop          *goprop.Prop
	log           *logger.Logger
	userClient    UserClient.UserClient
	k8sClient     k8s.K8sClient
	vcodeClient   vcode.VcodeClient
	monitorClient monitor.MonitorCaller
	searchClient  search.SearchCaller
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
	/* init searchClient */
	searchClient = search.NewSearchCaller(redisClient)

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
	//初始化验证码
	vcodeClient = vcode.NewVcodeClient(k8sClient, tracer, namespace, serviceName, "vcode", "8888")
	//初始化id生成器
	idGenClient = idgen.NewIdGenClient(k8sClient, tracer, namespace)

	/* create service */
	var svc UserService
	svc = userService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	RegHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "Reg")(makeRegEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "Reg", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkRegHandler := mkhttpserver.NewHttpMtalkServer(RegHandler, serviceName, "Reg", monitorClient, &userResponse{Err: "panic"})

	LoginHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "Login")(makeLoginEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "Login", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkLoginHandler := mkhttpserver.NewHttpMtalkServer(LoginHandler, serviceName, "Login", monitorClient, &userResponse{Err: "panic"})

	ShortcutLoginHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ShortcutLogin")(makeShortcutLoginEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ShortcutLogin", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkShortcutLoginHandler := mkhttpserver.NewHttpMtalkServer(ShortcutLoginHandler, serviceName, "ShortcutLogin", monitorClient, &userResponse{Err: "panic"})

	FindPasswordHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "FindPassword")(makeFindPasswordEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "FindPassword", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkFindPasswordHandler := mkhttpserver.NewHttpMtalkServer(FindPasswordHandler, serviceName, "FindPassword", monitorClient, &userResponse{Err: "panic"})

	ChangePhoneNoHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ChangePhoneNo")(makeChangePhoneNoEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ChangePhoneNo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkChangePhoneNoHandler := mkhttpserver.NewHttpMtalkServer(ChangePhoneNoHandler, serviceName, "ChangePhoneNo", monitorClient, &userResponse{Err: "panic"})

	OtherLoginHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "OtherLogin")(makeOtherLoginEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "OtherLogin", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkOtherLoginHandler := mkhttpserver.NewHttpMtalkServer(OtherLoginHandler, serviceName, "OtherLogin", monitorClient, &userResponse{Err: "panic"})

	UpdateUserHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateUser")(makeUpdateUserEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateUser", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkUpdateUserHandler := mkhttpserver.NewHttpMtalkServer(UpdateUserHandler, serviceName, "UpdateUser", monitorClient, &userResponse{Err: "panic"})

	MyHomeHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "MyHome")(makeMyHomeEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "MyHome", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkMyHomeHandler := mkhttpserver.NewHttpMtalkServer(MyHomeHandler, serviceName, "MyHome", monitorClient, &userResponse{Err: "panic"})

	HomePageSloganHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "HomePageSlogan")(makeHomePageSloganEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "HomePageSlogan", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkHomePageSloganHandler := mkhttpserver.NewHttpMtalkServer(HomePageSloganHandler, serviceName, "HomePageSlogan", monitorClient, &userResponse{Err: "panic"})

	CheckPhoneBookHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "CheckPhoneBook")(makeCheckPhoneBookEndpoint(svc)),
		decodePhonnesBookRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CheckPhoneBook", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkCheckPhoneBookHandler := mkhttpserver.NewHttpMtalkServer(CheckPhoneBookHandler, serviceName, "CheckPhoneBook", monitorClient, &userResponse{Err: "panic"})

	PhoneBookUserHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "PhoneBookUser")(makePhoneBookUserEndpoint(svc)),
		decodePhonnesBookRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "PhoneBookUser", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkPhoneBookUserHandler := mkhttpserver.NewHttpMtalkServer(PhoneBookUserHandler, serviceName, "PhoneBookUser", monitorClient, &userResponse{Err: "panic"})

	ActiveUserHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ActiveUser")(makeActiveUserEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ActiveUser", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkActiveUserHandler := mkhttpserver.NewHttpMtalkServer(ActiveUserHandler, serviceName, "ActiveUser", monitorClient, &userResponse{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/reg", mtalkRegHandler)
	http.Handle("/login", mtalkLoginHandler)
	http.Handle("/shortcutLogin", mtalkShortcutLoginHandler)
	http.Handle("/findPassword", mtalkFindPasswordHandler)
	http.Handle("/changePhoneNo", mtalkChangePhoneNoHandler)
	http.Handle("/otherLogin", mtalkOtherLoginHandler)
	http.Handle("/updateUser", mtalkUpdateUserHandler)
	http.Handle("/myHome", mtalkMyHomeHandler)
	http.Handle("/homePageSlogan", mtalkHomePageSloganHandler)
	http.Handle("/checkPhoneBook", mtalkCheckPhoneBookHandler)
	http.Handle("/phoneBookUser", mtalkPhoneBookUserHandler)
	http.Handle("/activeUser", mtalkActiveUserHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{RedisClient: redisClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

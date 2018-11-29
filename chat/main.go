package main

import (
	sdk "chat/rcserversdk"
	"context"
	"flag"
	idgen "idgen/client"
	monitor "monitor/client"
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"
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
	listenPort  = ":8888"
	serviceName = "chat"
	RyUserId    = "chat:rytoken:keyIsUserId:"
	UserIdRy    = "chat:userId:keyIsryToken:"

	UserId       = "U15263480470032"
	AppKeyStr    = "uwd1c0sxuppq1"
	AppSecretStr = "dcBkLsLoKgl4"

	defaultNickName  = "钦家用户"
	defaultImageName = " "
	groupChatNum     = 500

	//redis 用户id
	userId                  = "U:"
	tokenStr                = "token:"
	groupChatInfo           = "chat:groupChatInfo:"
	groupChatMemberList     = "chat:groupChatMemberList:"
	memberJoinGroupChatList = "chat:memberJoinGroupChatList:"

	Count      = "3"
	rytoken101 = "chat101" //必要参数为空
	rytoken102 = "chat102" //Status不等于字符串0 或1
	rytoken103 = "chat103" //改用户没注册
	rytoken104 = "chat104" //取不到该用户的信息
	rytoken105 = "chat105" //redis取不到 userId ， nickName ，imageName其中某一个
	rytoken106 = "chat106" //  根据已有参数得不到Userid（该用户没登录）
	rytoken107 = "chat107" //根据条件查不到数据
	chat108    = "chat108" //群聊成员数量为0，请排查后台数据
	chat109    = "chat109" //此次操作群聊成员数量
	chat110    = "chat110" // GroupChatUrl， GroupChatName，  GroupChatNotice，请确保有一个值不为空
	chat111    = "chat111" //redis里没有该群聊的信息
	chat112    = "chat112" //该用户不是群主，无权操作
	chat113    = "chat113" //该活动无群聊
	chat114    = "chat114" //群聊不存在
)

var (
	namespace      string
	idGenClient    idgen.IdGenClient
	mysqlClient    mysql.MysqlClient
	tracer         stdopentracing.Tracer
	redisClient    redis.RedisClient
	prop           *goprop.Prop
	log            *logger.Logger
	monitorClient  monitor.MonitorCaller
	fieldKeys      []string
	options        []httptransport.ServerOption
	zipkinAddr     string
	requestLatency *kitprometheus.Summary
	requestCount   *kitprometheus.Counter
	group          *sdk.Group
	message        *sdk.Message
	user           *sdk.User
	dao            *Dao
)

func init() {
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop = goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace
	zipkinAddr = prop.Get("zipkinAddr")

	/* init log */
	logger.SetDefaultLogLevel(logger.LevelDebug)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	max, _ := strconv.Atoi(prop.Get("mysql_chat_maxidleconn"))
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     prop.Get("mysql_chat_username"),
		Password:     prop.Get("mysql_chat_password"),
		IP:           prop.Get("mysql_chat_host"),
		Port:         prop.Get("mysql_chat_port"),
		DatabaseName: prop.Get("mysql_db_chat"),
		Logger:       logger.GetDefaultLogger(),
		MaxIdleConns: max,
	})

	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"), //127.0.0.1:6379   192.168.10.163:8888
		RedisPassword: prop.Get("redis_password"),
	})

	options = []httptransport.ServerOption{
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
	fieldKeys = []string{"method", "error"}
	requestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: serviceName,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: serviceName,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)
	k8sClient := k8s.NewK8sClient()
	idGenClient = idgen.NewIdGenClient(k8sClient, tracer, namespace)

	group = &sdk.Group{AppKey: AppKeyStr, AppSecret: AppSecretStr}

	message = &sdk.Message{AppKey: AppKeyStr, AppSecret: AppSecretStr}
	dao = &Dao{}
	user = &sdk.User{AppKey: AppKeyStr, AppSecret: AppSecretStr}
}
func main() {

	/* create service */
	var svc RYToKenService
	svc = &rYToKenService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	//parkour======================================================start
	AddOfficialMSG := httptransport.NewServer(
		opentracing.TraceServer(tracer, "AddOfficialMSG")(makeAddOfficialMSGEndpoint(svc)),
		decodeOfficialMSGRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "AddOfficialMSG", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkAddOfficialMSGHandler := mkhttpserver.NewHttpMtalkServer(AddOfficialMSG, serviceName, "AddOfficialMSG", monitorClient, &panicResponse{Err: "panic"})

	LookOfficialMSG := httptransport.NewServer(
		opentracing.TraceServer(tracer, "LookOfficialMSG")(makeLookOfficialMSGEndpoint(svc)),
		decodeOfficialMSGRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "LookOfficialMSG", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkLookOfficialMSGHandler := mkhttpserver.NewHttpMtalkServer(LookOfficialMSG, serviceName, "LookOfficialMSG", monitorClient, &panicResponse{Err: "panic"})

	InformChat := httptransport.NewServer(
		opentracing.TraceServer(tracer, "InformChat")(makeAddOfficialMSGEndpoint(svc)),
		decodeOfficialMSGRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "InformChat", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkInformChatHandler := mkhttpserver.NewHttpMtalkServer(InformChat, serviceName, "InformChat", monitorClient, &panicResponse{Err: "panic"})

	//parkour======================================================end

	getRYToken := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetRYToken")(makeGetRYTokenEndpoint(svc)),
		decodeDefaRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetRYToken", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetRYTokenHandler := mkhttpserver.NewHttpMtalkServer(getRYToken, serviceName, "GetRYToken", monitorClient, &panicResponse{Err: "panic"})

	getUserInfo := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetUserInfo")(makeGetUserInfoEndpoint(svc)),
		decodeDefaRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetUserInfo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetUserInfoHandler := mkhttpserver.NewHttpMtalkServer(getUserInfo, serviceName, "GetUserInfo", monitorClient, &panicResponse{Err: "panic"})

	createGroupChat := httptransport.NewServer(
		opentracing.TraceServer(tracer, "CreateGroupChat")(makeCreateGroupChatEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CreateGroupChat", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	//	mtalkCreateGroupChatHandler := mkhttpserver.NewHttpMtalkServer(createGroupChat, serviceName, "CreateGroupChat", monitorClient, &panicResponse{Err: "panic"})

	joinGroupChat := httptransport.NewServer(
		opentracing.TraceServer(tracer, "JoinGroupChat")(makeJoinGroupChatEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "JoinGroupChat", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkJoinGroupChatHandler := mkhttpserver.NewHttpMtalkServer(joinGroupChat, serviceName, "JoinGroupChat", monitorClient, &panicResponse{Err: "panic"})

	QuitGroupChat := httptransport.NewServer(
		opentracing.TraceServer(tracer, "QuitGroupChat")(makeQuitGroupChatEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "QuitGroupChat", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkQuitGroupChatHandler := mkhttpserver.NewHttpMtalkServer(QuitGroupChat, serviceName, "QuitGroupChat", monitorClient, &panicResponse{Err: "panic"})

	queryGroupChatMemberList := httptransport.NewServer(
		opentracing.TraceServer(tracer, "QueryGroupChatMemberList")(makeQueryGroupChatMemberListEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "QueryGroupChatMemberList", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkQueryGroupChatMemberListHandler := mkhttpserver.NewHttpMtalkServer(queryGroupChatMemberList, serviceName, "QueryGroupChatMemberList", monitorClient, &panicResponse{Err: "panic"})

	getArrayGroupInfo := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetArrayGroupInfo")(makeGetArrayGroupInfoEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetArrayGroupInfo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetArrayGroupInfoHandler := mkhttpserver.NewHttpMtalkServer(getArrayGroupInfo, serviceName, "GetArrayGroupInfo", monitorClient, &panicResponse{Err: "panic"})

	getMyGroupChat := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetMyGroupChat")(makeGetMyGroupChatEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetMyGroupChat", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetMyGroupChatHandler := mkhttpserver.NewHttpMtalkServer(getMyGroupChat, serviceName, "GetMyGroupChat", monitorClient, &panicResponse{Err: "panic"})

	updateGroupChat := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateGroupChat")(makeUpdateGroupChatEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateGroupChat", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkUpdateGroupChatHandler := mkhttpserver.NewHttpMtalkServer(updateGroupChat, serviceName, "UpdateGroupChat", monitorClient, &panicResponse{Err: "panic"})

	dismiss := httptransport.NewServer(
		opentracing.TraceServer(tracer, "Dismiss")(makeDismissEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "Dismiss", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkDismissHandler := mkhttpserver.NewHttpMtalkServer(dismiss, serviceName, "Dismiss", monitorClient, &panicResponse{Err: "panic"})

	searchChatInfo := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SearchChatInfo")(makeSearchChatInfoEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SearchChatInfo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkSearchChatInfoHandler := mkhttpserver.NewHttpMtalkServer(searchChatInfo, serviceName, "SearchChatInfo", monitorClient, &panicResponse{Err: "panic"})

	getClassId := httptransport.NewServer(
		opentracing.TraceServer(tracer, "GetClassId")(makeGetClassIdEndpoint(svc)),
		decodeGroupChatInfoquest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetClassId", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	mtalkGetClassIdHandler := mkhttpserver.NewHttpMtalkServer(getClassId, serviceName, "GetClassId", monitorClient, &panicResponse{Err: "panic"})
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	http.Handle("/", hystrixStreamHandler)

	http.Handle("/getRYToken", mtalkGetRYTokenHandler)                             //mtalkGetRYTokenHandler getRYToken
	http.Handle("/getUserInfo", mtalkGetUserInfoHandler)                           //getUserInfo  mtalkGetUserInfoHandler
	http.Handle("/createGroupChat", createGroupChat)                               //createGroupChat  mtalkCreateGroupChatHandler
	http.Handle("/joinGroupChat", mtalkJoinGroupChatHandler)                       //joinGroupChat  mtalkJoinGroupChatHandler
	http.Handle("/quitGroupChat", mtalkQuitGroupChatHandler)                       //joinGroupChat  mtalkQuitGroupChatHandler
	http.Handle("/queryGroupChatMemberList", mtalkQueryGroupChatMemberListHandler) //queryGroupChatMemberList  mtalkQueryGroupChatMemberListHandler
	http.Handle("/getArrayGroupInfo", mtalkGetArrayGroupInfoHandler)               //mtalkGetArrayGroupInfoHandler
	http.Handle("/getMyGroupChat", mtalkGetMyGroupChatHandler)                     //    getMyGroupChat  mtalkGetMyGroupChatHandler
	http.Handle("/updateGroupChat", mtalkUpdateGroupChatHandler)                   //updateGroupChat mtalkUpdateGroupChatHandler
	http.Handle("/dismiss", mtalkDismissHandler)                                   //dismiss  mtalkDismissHandler

	//parkour======================================================start
	http.Handle("/addOfficialMSG", mtalkAddOfficialMSGHandler)
	http.Handle("/lookOfficialMSG", mtalkLookOfficialMSGHandler)
	http.Handle("/informChat", mtalkInformChatHandler)
	http.Handle("/getClassId", mtalkGetClassIdHandler)
	//parkour======================================================end

	http.Handle("/searchChatInfo", mtalkSearchChatInfoHandler)
	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{RedisClient: redisClient, MySqlTableName: "groupchat", MysqlClient: mysqlClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

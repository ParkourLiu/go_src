package main

import (
	chat "chat/client"
	idgen "idgen/client"
	monitor "monitor/client"
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	monitorhelper "mtcomm/helper/monitor"

	"context"
	"flag"
	"mtcomm/k8s"
	logger "mtcomm/log"
	mkhttpserver "mtcomm/middleware/handler/httpserver"
	"net/http"
	"os"
	push "push/client"
	"strconv"

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
	serviceName = "mclass"
	flag1       = serviceName + "101" //邀请码不存在
	flag2       = serviceName + "102" //重复加入班级
	flag3       = serviceName + "103" //无权操作
	flag4       = serviceName + "104" //班级id不存在
	flag10      = serviceName + "110" //数据异常
	flag11      = serviceName + "111" //正在审核中
	flag12      = serviceName + "112" //Id不存在

)

var (
	namespace     string
	redisClient   redis.RedisClient
	mysqlClient   mysql.MysqlClient
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	monitorClient monitor.MonitorCaller
	k8sClient     k8s.K8sClient
	idGenClient   idgen.IdGenClient
	//oss
	key_id                string
	key_secret            string
	bucket_name_mtalk     string
	oss_endpoint          string
	domain_app            string
	mgtalk2Url            string //mgtalk2请求url
	chatClient            chat.ChatClient
	pushClient            push.PushCaller
	InviteUrl             string
	baidu_face_API_Key    string //人脸识别库API_Key
	baidu_face_Secret_Key string //人脸识别库Secret_Key
	bucket_name_face      string
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
	oss_endpoint = prop.Get("aliyun_oss_endpoint")
	domain_app = prop.Get("domain_app")
	mgtalk2Url = prop.Get("mgtalk2Url")
	InviteUrl = prop.Get("InviteUrl")
	baidu_face_API_Key = prop.Get("baidu_face_API_Key")
	baidu_face_Secret_Key = prop.Get("baidu_face_Secret_Key")
	bucket_name_face = prop.Get("aliyun_oss_bucket_name_face")
	/* init log */
	logLevel, _ := strconv.Atoi(prop.Get("log_level"))
	logger.SetDefaultLogLevel(logLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	/* init mysql */
	//max, _ := strconv.Atoi(prop.Get("mysql_grp_member_maxidleconn"))
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     prop.Get("mysql_mschool_username"),
		Password:     prop.Get("mysql_mschool_password"),
		IP:           prop.Get("mysql_mschool_host"),
		Port:         prop.Get("mysql_mschool_port"),
		DatabaseName: prop.Get("mysql_db_mschool"),
		Logger:       logger.GetDefaultLogger(),
		//	MaxIdleConns: max,
	})
	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})
	k8sClient = k8s.NewK8sClient()
	zipkinAddr := prop.Get("zipkinAddr")

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
	chatClient = chat.NewChatClient(k8sClient, tracer, namespace, serviceName, "chat", "8888")
	pushClient = push.NewPushCaller(redisClient)
}
func main() {
	// init

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(*logger.GetDefaultLogger().GetDefaultKitLogger()),
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
	var svc McreateClassService
	svc = createClassService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	SelectByScInviteCodeHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "SelectByScInviteCode")(makeSelectByScInviteCodeEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "SelectByScInviteCode", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)

	CreateClassHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "CreateClass")(makeAddClassEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "AddClass", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)

	TeacherJoinClassHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "TeacherJoinClass")(makeTeacherJoinClassEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "TeacherJoinClass", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)

	FamilyJoinClassHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "FamilyJoinClass")(makeFamilyJoinClassEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "FamilyJoinClass", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	NewMemberHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "LeaveMessage")(makeNewMemberEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "NewMember", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	ApproveMembersHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ApproveMembers")(makeApproveMembersEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ApproveMembers", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	ManagerMemberHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ManagerMember")(makeManagerMemberEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "makeManagerMemberEndpoint", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	OperateMemberHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "OperateMember")(makeOperateMemberEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "makeOperateMemberEndpoint", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	ClassQrCodeHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ClassQrCode")(makeClassQrCodeEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "makeClassQrCodeEndpoint", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	UpdateStudentHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateStudent")(makeUpdateStudentEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateStudent", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	FindAllMemberHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "FindAllMember")(makeFindAllMemberEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "FindAllMember", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	FindTeacherMemberHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "FindTeacherMember")(makeFindTeacherMemberEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "FindTeacherMember", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	UpdateGroupChatInfoHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateGroupChatInfo")(makeUpdateGroupChatInfoEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateGroupChatInfo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)
	UpdateTeachInfoHandler := httptransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateGroupChatInfo")(makeUpdateTeachInfoEndpoint(svc)),
		decodeSelectByScInviteCodeRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateTeachInfo", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)

	mtalkSelectByScInviteCodeHandler := mkhttpserver.NewHttpMtalkServer(SelectByScInviteCodeHandler, serviceName, "SelectByScInviteCode", monitorClient, &SchoolYearResponse{Err: "panic"})
	mtalkCreateClassHandler := mkhttpserver.NewHttpMtalkServer(CreateClassHandler, serviceName, "CreateClass", monitorClient, &AddClassResponse{Err: "panic"})
	mtalkTeacherJoinClassHandler := mkhttpserver.NewHttpMtalkServer(TeacherJoinClassHandler, serviceName, "TeacherJoinClass", monitorClient, &Response{Err: "panic"})
	mtalkFamilyJoinClassClassHandler := mkhttpserver.NewHttpMtalkServer(FamilyJoinClassHandler, serviceName, "FamilyJoinClass", monitorClient, &Response{Err: "panic"})
	mtalkNewMemberHandler := mkhttpserver.NewHttpMtalkServer(NewMemberHandler, serviceName, "NewMember", monitorClient, &GroupResponse{Err: "panic"})
	mtalkApproveMembersHandler := mkhttpserver.NewHttpMtalkServer(ApproveMembersHandler, serviceName, "ApproveMembers", monitorClient, &Response{Err: "panic"})
	mtalkManagerMemberHandler := mkhttpserver.NewHttpMtalkServer(ManagerMemberHandler, serviceName, "ManagerMember", monitorClient, &GroupResponse{Err: "panic"})
	mtalkOperateMemberHandler := mkhttpserver.NewHttpMtalkServer(OperateMemberHandler, serviceName, "OperateMember", monitorClient, &Response{Err: "panic"})
	mtalkClassQrCodeHandler := mkhttpserver.NewHttpMtalkServer(ClassQrCodeHandler, serviceName, "ClassQrCode", monitorClient, &AddClassResponse{Err: "panic"})
	mtalkUpdateStudentHandler := mkhttpserver.NewHttpMtalkServer(UpdateStudentHandler, serviceName, "UpdateStudent", monitorClient, &Response{Err: "panic"})
	mtalkFindAllMemberHandler := mkhttpserver.NewHttpMtalkServer(FindAllMemberHandler, serviceName, "FindAllMember", monitorClient, &Response{Err: "panic"})
	mtalkFindTeacherMemberHandler := mkhttpserver.NewHttpMtalkServer(FindTeacherMemberHandler, serviceName, "FindTeacherMember", monitorClient, &Response{Err: "panic"})
	mtalkUpdateGroupChatInfoHandler := mkhttpserver.NewHttpMtalkServer(UpdateGroupChatInfoHandler, serviceName, "UpdateGroupChatInfo", monitorClient, &Response{Err: "panic"})
	mtalkUpdateTeachInfoHandler := mkhttpserver.NewHttpMtalkServer(UpdateTeachInfoHandler, serviceName, "UpdateTeachInfo", monitorClient, &Response{Err: "panic"})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	http.Handle("/", hystrixStreamHandler)
	http.Handle("/selectByScInviteCode", mtalkSelectByScInviteCodeHandler)
	http.Handle("/createClass", mtalkCreateClassHandler)
	http.Handle("/teacherJoinClass", mtalkTeacherJoinClassHandler)
	http.Handle("/familyJoinClass", mtalkFamilyJoinClassClassHandler)
	http.Handle("/newMember", mtalkNewMemberHandler)
	http.Handle("/approveMembers", mtalkApproveMembersHandler)
	http.Handle("/managerMember", mtalkManagerMemberHandler)
	http.Handle("/operateMember", mtalkOperateMemberHandler)
	http.Handle("/classQrCode", mtalkClassQrCodeHandler)
	http.Handle("/updateStudent", mtalkUpdateStudentHandler)
	http.Handle("/findAllMember", mtalkFindAllMemberHandler)
	http.Handle("/findTeacherMember", mtalkFindTeacherMemberHandler)
	http.Handle("/updateGroupChatInfo", mtalkUpdateGroupChatInfoHandler)
	http.Handle("/updateTeachInfo", mtalkUpdateTeachInfoHandler)

	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/monitor", &monitorhelper.MonitorHandler{MySqlTableName: "schoolyear", MysqlClient: mysqlClient, RedisClient: redisClient})
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

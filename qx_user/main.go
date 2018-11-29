package main

import (
	"flag"
	"net"
	"os"
	"strconv"
	"mtcomm/k8s"
	logger "mtcomm/log"
	mkgrpcserver "mtcomm/middleware/handler/grpcserver"

	"golang.org/x/net/context"

	"qx_user/pb"
	idgen "qx_idgen/client"
	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/db/mysql"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkinopt "github.com/openzipkin/zipkin-go-opentracing"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tjz101/goprop"
	"google.golang.org/grpc"
)

const (
	listenPort  = ":8888"
	serviceName = "qx_user"
)

var (
	namespace     string
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	monitorClient monitor.MonitorCaller
	redisClient   redis.RedisClient
	mysqlClient   mysql.MysqlClient
	idGenClient   idgen.IdGenClient
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
	logLevel, _ := strconv.Atoi(prop.Get("log_level"))
	logger.SetDefaultLogLevel(logLevel)
	logger.With("serviceName", serviceName)
	log = logger.GetDefaultLogger()

	/* init mysql */
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     prop.Get("mysql_user_username"),
		Password:     prop.Get("mysql_user_password"),
		IP:           prop.Get("mysql_user_host"),
		Port:         prop.Get("mysql_user_port"),
		DatabaseName: prop.Get("mysql_db_user"),
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

	//初始化id生成器
	idGenClient = idgen.NewIdGenClient(k8sClient, tracer, namespace)
}

func main() {
	// init
	zipkinAddr := prop.Get("zipkinAddr")

	// init tracing domain.
	{
		if zipkinAddr != "" {
			logger.Info("tracer", "Zipkin", "zipkinAddr", zipkinAddr)
			collector, err := zipkinopt.NewHTTPCollector(zipkinAddr, zipkinopt.HTTPBatchSize(10))
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
			tracer, err = zipkinopt.NewTracer(
				zipkinopt.NewRecorder(collector, false, listenPort, serviceName),
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

	var zipkinTracer *zipkingo.Tracer
	{
		var (
			err           error
			useNoopTracer = (zipkinAddr == "")
			reporter      = zipkinhttp.NewReporter(zipkinAddr)
		)
		defer reporter.Close()
		zEP, _ := zipkingo.NewEndpoint(serviceName, zipkinAddr)
		zipkinTracer, err = zipkingo.NewTracer(
			reporter, zipkingo.WithLocalEndpoint(zEP), zipkingo.WithNoopTracer(useNoopTracer),
		)
		if err != nil {
			log.Error("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			log.Info("tracer", "Zipkin", "type", "Native", "URL", zipkinAddr)
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

	zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(*log.GetDefaultKitLogger()),
		zipkinServer,
	}

	SearchUserByIdHandler := grpctransport.NewServer(
		opentracing.TraceServer(tracer, "SearchUserById")(makeSearchUserByIdEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "SearchUserById", *log.GetDefaultKitLogger())))...,
	)
	mtalkSearchUserByIdHandler := mkgrpcserver.NewGrpcMtalkServer(SearchUserByIdHandler, serviceName, "SearchUserById", monitorClient)

	SearchUsersHandler := grpctransport.NewServer(
		opentracing.TraceServer(tracer, "SearchUsers")(makeSearchUsersEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "SearchUsers", *log.GetDefaultKitLogger())))...,
	)
	mtalkSearchUsersHandler := mkgrpcserver.NewGrpcMtalkServer(SearchUsersHandler, serviceName, "SearchUsers", monitorClient)

	AddUserHandler := grpctransport.NewServer(
		opentracing.TraceServer(tracer, "AddUser")(makeAddUserEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "AddUser", *log.GetDefaultKitLogger())))...,
	)
	mtalkAddUserHandler := mkgrpcserver.NewGrpcMtalkServer(AddUserHandler, serviceName, "AddUser", monitorClient)

	UpdateUserHandler := grpctransport.NewServer(
		opentracing.TraceServer(tracer, "UpdateUser")(makeUpdateUserEndpoint(svc)),
		decodeRequest,
		encodeResponse,
		append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "UpdateUser", *log.GetDefaultKitLogger())))...,
	)
	mtalkUpdateUserHandler := mkgrpcserver.NewGrpcMtalkServer(UpdateUserHandler, serviceName, "UpdateUser", monitorClient)

	/* grpc server */
	grpcSvr := &grpcServer{
		searchUserById: mtalkSearchUserByIdHandler,
		searchUsers:    mtalkSearchUsersHandler,
		addUser:        mtalkAddUserHandler,
		updateUser:     mtalkUpdateUserHandler,
	}
	// The gRPC listener mounts the Go kit gRPC server we created.
	grpcListener, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Error("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	log.Error("transport", "gRPC", "addr", listenPort)
	baseServer := grpc.NewServer()
	pb.RegisterUserServiceServer(baseServer, grpcSvr)
	baseServer.Serve(grpcListener)
}

type grpcServer struct {
	searchUserById grpctransport.Handler
	searchUsers    grpctransport.Handler
	addUser        grpctransport.Handler
	updateUser     grpctransport.Handler
}

func (s *grpcServer) SearchUserById(ctx context.Context, req *pb.UserRequest) (*pb.UserReply, error) {
	_, rep, err := s.searchUserById.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserReply), nil
}
func (s *grpcServer) SearchUsers(ctx context.Context, req *pb.UserRequest) (*pb.UsersReply, error) {
	_, rep, err := s.searchUsers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UsersReply), nil
}
func (s *grpcServer) AddUser(ctx context.Context, req *pb.UserRequest) (*pb.UserReply, error) {
	_, rep, err := s.addUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserReply), nil
}
func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserReply, error) {
	_, rep, err := s.updateUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserReply), nil
}

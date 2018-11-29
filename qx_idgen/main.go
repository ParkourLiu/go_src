package main

import (
	"flag"
	"net"
	"os"
	"strconv"
	"strings"

	logger "mtcomm/log"
	mkgrpcserver "mtcomm/middleware/handler/grpcserver"

	"golang.org/x/net/context"

	"idgen/pb"
	monitor "monitor/client"
	"mtcomm/db/redis"

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
	serviceName = "qx_idgen"
)

var (
	namespace     string
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	svc           IdGeneraterService
	monitorClient monitor.MonitorCaller
	redisClient   redis.RedisClient
	seqNo         string //5bit
	nodeNo        int64  // datacenter (5bit) + nodeSeq (5bit)
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

	/* get client seq no */
	hostName := os.Getenv("HOST_NAME")
	if hostName == "" {
		panic("Don't set HOST_NAME env!")
	}
	log.Debug("hostName", hostName)
	array := strings.Split(hostName, "-")
	seqNo = array[len(array)-1]
	_, err := strconv.Atoi(seqNo)
	if err != nil {
		panic("HOST_INDEX is error!")
	}
	log.Debug("HOST_INDEX", seqNo)

	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

	/* init monitorClient */
	monitorClient = monitor.NewMonitorCaller(redisClient)

	/* init node no for id gen */
	var dn, hn int64
	dn, err = strconv.ParseInt(prop.Get("snowflake_datacenter"), 10, 64)
	if err != nil {
		panic("snowflake_datacenter is wrong")
	}
	hn, err = strconv.ParseInt(seqNo, 10, 64)
	if err != nil {
		panic("snowflake_host_seqNo is wrong")
	}
	nodeNo = int64(dn<<5 | hn)
	log.Debug("nodeNo", nodeNo)
}

func main() {
	// init
	zipkinAddr := prop.Get("zipkinAddr")

	// init tracing domain.
	{
		if zipkinAddr != "" {
			logger.Info("tracer", "Zipkin", "zipkinAddr", zipkinAddr)
			collector, err := zipkinopt.NewHTTPCollector(zipkinAddr, zipkinopt.HTTPBatchSize(1))
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
	svc = idGeneraterService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(*log.GetDefaultKitLogger()),
		zipkinServer,
	}

	generateUniqueIdV1Handler := grpctransport.NewServer(
		opentracing.TraceServer(tracer, "GenerateUniqueIdV1")(makeGenerateUniqueIdV1Endpoint(svc)),
		decodeGenerateUniqueIdV1Request,
		encodeGenerateUniqueIdV1Response,
		append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GenerateUniqueIdV1", *log.GetDefaultKitLogger())))...,
	)
	mtalkGenerateUniqueIdV1Handler := mkgrpcserver.NewGrpcMtalkServer(generateUniqueIdV1Handler, serviceName, "GenerateUniqueIdV1", monitorClient)

	/* grpc server */
	grpcSvr := &grpcServer{
		generateUniqueIdV1: mtalkGenerateUniqueIdV1Handler,
	}
	// The gRPC listener mounts the Go kit gRPC server we created.
	grpcListener, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Error("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	log.Error("transport", "gRPC", "addr", listenPort)
	baseServer := grpc.NewServer()
	pb.RegisterIdGeneraterServer(baseServer, grpcSvr)
	baseServer.Serve(grpcListener)
}

type grpcServer struct {
	generateUniqueIdV1 grpctransport.Handler
}

func (s *grpcServer) GenerateUniqueIdV1(ctx context.Context, req *pb.GenerateUniqueIdV1Request) (*pb.GenerateUniqueIdV1Reply, error) {
	_, rep, err := s.generateUniqueIdV1.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GenerateUniqueIdV1Reply), nil
}

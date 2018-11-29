package main

import (
	"context"
	"flag"
	monitor "monitor/client"
	"mtcomm/db/redis"
	"mtcomm/helper/health"
	"mtcomm/k8s"
	logger "mtcomm/log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/tjz101/goprop"
)

const (
	listenPort  = ":8888"
	serviceName = "gateway"
)

var (
	namespace     string
	tracer        stdopentracing.Tracer
	prop          *goprop.Prop
	log           *logger.Logger
	k8sClient     k8s.K8sClient
	monitorClient monitor.MonitorCaller
	redisClient   redis.RedisClient
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

	k8sClient = k8s.NewK8sClient()

	/* init redis */
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     prop.Get("redis_host"),
		RedisPassword: prop.Get("redis_password"),
	})

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

	/* create service */
	var svc GatewayService
	svc = gatewayService{}
	svc = loggingMiddleware{svc}

	handlerFunc := httptransport.NewServer(
		opentracing.TraceServer(tracer, "ForwardRequest")(makeGatewayEndpoint(svc)),
		decodeForwardRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "ForwardRequest", *logger.GetDefaultLogger().GetDefaultKitLogger())))...,
	)

	r := mux.NewRouter()
	r.PathPrefix("/mtalk").Handler(handlerFunc).Methods("POST")
	http.HandleFunc("/health", health.HealthHandler)
	http.Handle("/", r)
	logger.Info("msg", "HTTP", "addr", listenPort)
	logger.Info("err", http.ListenAndServe(listenPort, nil))
}

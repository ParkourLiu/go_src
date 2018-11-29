package main

import (
	"testing"
	logger "mtcomm/log"
	"os"
	zipkinopt "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdopentracing "github.com/opentracing/opentracing-go"
)


var _  = os.Setenv("HOST_NAME", "lio-1")


func TestSaveUser_1(t *testing.T) {
	// init
	zipkinAddr := prop.Get("zipkinAddr")
	listenPort := prop.Get("listenPort")

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
	
	a1, err := svc.GenerateUniqueIdV1(0)
	if err != nil {
		t.Error(err)
	}
	if len(a1) != 0 {
		t.Error("error")
	}
	
	a2, err := svc.GenerateUniqueIdV1(5)
	if err != nil {
		t.Error(err)
	}
	if len(a2) != 5 {
		t.Error("error")
	}	
}

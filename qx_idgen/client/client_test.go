package client_test

import (
	"fmt"
	"idgen/client"
	logger "mtcomm/log"
	"testing"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinopt "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) GetHost(namespace string, serviceName string) (string, error) {
	args := m.Called(namespace, serviceName)
	return args.String(0), nil
}

func (m *MyMockedObject) IsClusterEnv() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MyMockedObject) GetFirstPort(namespace string, serviceName string) (string, error) {
	return "", nil
}

var (
	tracer  stdopentracing.Tracer
	testObj *MyMockedObject
	c       client.IdGenClient
)

func init() {
	// create an instance of our test object
	testObj = new(MyMockedObject)

	// setup expectations
	testObj.On("GetHost", "default", "idgen-service").Return("localhost", nil)
	testObj.On("IsClusterEnv").Return(true)

	serviceName := "test"
	listenPort := ":8888"
	/* init log */
	logger.SetDefaultLogLevel(logger.LevelDebug)
	log := logger.GetDefaultLogger()
	zipkinAddr := "http://192.168.10.164:9411/api/v1/spans"
	// init tracing domain.
	{
		if zipkinAddr != "" {
			log.Info("tracer", "Zipkin", "zipkinAddr", zipkinAddr)
			collector, err := zipkinopt.NewHTTPCollector(zipkinAddr, zipkinopt.HTTPBatchSize(1))
			if err != nil {
				log.Error("tracer", "Zipkin", "err", err)
			}
			tracer, err = zipkinopt.NewTracer(
				zipkinopt.NewRecorder(collector, false, listenPort, serviceName),
			)
			if err != nil {
				log.Error("tracer", "Zipkin", "err", err)
			}
		} else {
			log.Info("tracer", "none")
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	c = client.NewIdGenClient(testObj, tracer, "default")
}
func TestGetId(t *testing.T) {
	c := client.NewIdGenClient(testObj, tracer, "default")
	for i := 0; i < 10; i++ {
		fmt.Println(c.GetUniqueId())
	}
}

func BenchmarkGetTrendIncrId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.GetUniqueId()
	}
}

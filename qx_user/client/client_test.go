package client_test

import (
	"qx_user/client"
	logger "mtcomm/log"
	"testing"
	"qx_user/caller"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinopt "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/stretchr/testify/mock"
	"github.com/golang/go/src/pkg/fmt"
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
	c       client.UserClient
)

func init() {
	// create an instance of our test object
	testObj = new(MyMockedObject)

	// setup expectations
	testObj.On("GetHost", "default", "qx_user-service").Return("localhost", nil)
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

	c = client.NewUserClient(testObj, tracer, "default")
}

//运行4个后，查看结果
func TestSearchUserById(t *testing.T) {
	req, err := c.SearchUserById(&caller.UserRequest{UserId: "U01jkstv72oag3r"})
	fmt.Println("-------------req:", req)
	fmt.Println("-------------err:", err)
}

func TestSearchUsers(t *testing.T) {
	req, err := c.SearchUsers(&caller.UserRequest{UserId: "U01jkstv72oag3r"})
	fmt.Println("-------------req:", req)
	fmt.Println("-------------err:", err)
}

func TestAddUser(t *testing.T) {
	req, err := c.AddUser(&caller.UserRequest{PhoneNo: "12341234123"})
	fmt.Println("-------------req:", req)
	fmt.Println("-------------err:", err)
}

func TestUpdateUser(t *testing.T) {
	req, err := c.UpdateUser(&caller.UserRequest{UserId: "U01jkstv72oag3r", NickName: "aaaaaaaaa"})
	fmt.Println("-------------req:", req)
	fmt.Println("-------------err:", err)
}

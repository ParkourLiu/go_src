package grpc

import (
	"errors"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var commands map[string]bool //hystrix command

func init() {
	commands = make(map[string]bool)
}

type CallerParameter struct {
	ErrorPercentThreshold int
	Timeout               int
	MaxThread             int
	PbMethod              string //pb的方法名
	Enc                   grpctransport.EncodeRequestFunc
	Dec                   grpctransport.DecodeResponseFunc
	GrpcReply             interface{}
}

type CallProxy interface {
	MakeRemoteEndpoint(p *CallerParameter) (endpoint.Endpoint, error)
}

type CallProxyStruction struct {
	K8sClient         k8s.K8sClient
	Tracer            stdopentracing.Tracer
	CallerServiceName string //调用源的服务名
	Namespace         string //kubernetes namesapce
	PbServiceName     string //pb的服务名
	CalledServiceName string //被调用服务名
	CalledServicePort string //kubernetes service port
}

func initBreaker(b *CallerParameter, command string) {
	if b.Timeout == 0 {
		b.Timeout = 5000
	}
	if b.MaxThread == 0 {
		b.MaxThread = 100
	}
	if b.ErrorPercentThreshold == 0 {
		b.ErrorPercentThreshold = 50
	}

	hystrix.ConfigureCommand(command, hystrix.CommandConfig{
		Timeout:               b.Timeout,
		MaxConcurrentRequests: b.MaxThread,
		ErrorPercentThreshold: b.ErrorPercentThreshold,
	})
}

// Http Method: GET, POST etc.
func (cp *CallProxyStruction) MakeRemoteEndpoint(p *CallerParameter) (endpoint.Endpoint, error) {
	log := logger.GetDefaultLogger()
	log.Debug("start", "MakeRemoteEndpoint")
	if !cp.K8sClient.IsClusterEnv() {
		return nil, errors.New("It is not a k8s cluster env.")
	}

	command := cp.CallerServiceName + "_" + cp.CalledServiceName + "_" + p.PbMethod
	_, ok := commands[command]
	if !ok {
		//未配置hystrix
		log.Debug("command", command)
		initBreaker(p, command)
		commands[command] = true
	}

	host, err1 := cp.K8sClient.GetHost(cp.Namespace, cp.CalledServiceName+"-service")
	if err1 != nil {
		return nil, err1
	}

	instance := host + ":" + cp.CalledServicePort

	conn, err := grpc.Dial(instance, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, err
	}

	e := grpctransport.NewClient(
		conn,
		cp.PbServiceName,
		p.PbMethod,
		p.Enc,
		p.Dec,
		p.GrpcReply,
		grpctransport.ClientBefore(opentracing.ContextToGRPC(cp.Tracer, *logger.GetDefaultLogger().GetDefaultKitLogger())),
	).Endpoint()
	e = opentracing.TraceClient(cp.Tracer, p.PbMethod)(e)
	e = circuitbreaker.Hystrix(command)(e)
	return e, nil
}

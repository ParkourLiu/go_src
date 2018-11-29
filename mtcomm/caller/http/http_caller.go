package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"context"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	"mtcomm/middleware/retry"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
)

var commands map[string]bool //hystrix command

func init() {
	commands = make(map[string]bool)
}

type CallerParameter struct {
	ErrorPercentThreshold int
	Timeout               int
	MaxThread             int
	MethodName            string //http path的名字，不含斜杠
	HttpMethod            string
	EncodeRequestFunc     httptransport.EncodeRequestFunc
	DecodeResponseFunc    httptransport.DecodeResponseFunc
}

type CallProxy interface {
	MakeRemoteEndpoint(p *CallerParameter) (endpoint.Endpoint, error)
}

type CallProxyStruction struct {
	K8sClient         k8s.K8sClient
	Tracer            stdopentracing.Tracer
	CallerServiceName string //调用源的服务名
	Namespace         string //kubernetes namesapce
	CalledServiceName string //被调用服务名
	CalledServicePort string //kubernetes service port
}

func initBreaker(b *CallerParameter, command string) {
	if b.Timeout == 0 {
		b.Timeout = 5000
	}
	if b.MaxThread == 0 {
		b.MaxThread = 50
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
	//	if !cp.K8sClient.IsClusterEnv() {
	//		return nil, errors.New("It is not a k8s cluster env.")
	//	}
	command := cp.CallerServiceName + "_" + cp.CalledServiceName + "_" + p.MethodName
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

	instance := host + ":" + fmt.Sprint(cp.CalledServicePort)
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err2 := url.Parse(instance)
	if err2 != nil {
		return nil, err2
	}
	u.Path = "/" + p.MethodName

	var enc httptransport.EncodeRequestFunc
	if p.EncodeRequestFunc == nil {
		enc = encodeRequest
	} else {
		enc = p.EncodeRequestFunc
	}

	log.Debug("request_uri", instance+u.Path)
	e := httptransport.NewClient(
		p.HttpMethod,
		u,
		enc,
		p.DecodeResponseFunc,
		httptransport.SetClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					//					Timeout:   time.Duration(int64(p.Timeout/1000 + 1)),
					KeepAlive: 15 * time.Minute,
				}).DialContext,
				MaxIdleConns:          p.MaxThread + 10,
				IdleConnTimeout:       15 * time.Minute,
				TLSHandshakeTimeout:   0 * time.Second,
				ExpectContinueTimeout: 0 * time.Second,
				MaxIdleConnsPerHost:   p.MaxThread + 10,
			},
		}),
		httptransport.ClientBefore(opentracing.ContextToHTTP(cp.Tracer, *logger.GetDefaultLogger().GetDefaultKitLogger())),
	).Endpoint()
	e = opentracing.TraceClient(cp.Tracer, p.MethodName)(e)
	e = circuitbreaker.Hystrix(command)(e)
	e = retry.Retry()(e)
	return e, nil
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

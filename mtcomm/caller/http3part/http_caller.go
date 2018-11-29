package http3part

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"context"
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
	HttpMethod            string
	EncodeRequestFunc     httptransport.EncodeRequestFunc
	DecodeResponseFunc    httptransport.DecodeResponseFunc
}

type CallProxy interface {
	MakeRemoteEndpoint(p *CallerParameter) (endpoint.Endpoint, error)
}

type CallProxyStruction struct {
	Tracer      stdopentracing.Tracer
	CommandName string //格式为：调用源的服务名_调用目标服务名_调用方法名
	CalledUrl   string
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

	command := cp.CommandName
	_, ok := commands[command]
	if !ok {
		//未配置hystrix
		log.Debug("command", command)
		initBreaker(p, command)
		commands[command] = true
	}

	u, err2 := url.Parse(cp.CalledUrl)
	if err2 != nil {
		return nil, err2
	}

	var enc httptransport.EncodeRequestFunc
	if p.EncodeRequestFunc == nil {
		enc = encodeRequest
	} else {
		enc = p.EncodeRequestFunc
	}
	e := httptransport.NewClient(
		p.HttpMethod,
		u,
		enc,
		p.DecodeResponseFunc,
		httptransport.ClientBefore(opentracing.ContextToHTTP(cp.Tracer, *logger.GetDefaultLogger().GetDefaultKitLogger())),
	).Endpoint()
	e = opentracing.TraceClient(cp.Tracer, command)(e)
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

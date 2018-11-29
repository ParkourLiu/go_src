package main

import (
	"bytes"
	"context"
	"io/ioutil"
	ch "mtcomm/caller/http"
	k8s "mtcomm/k8s"
	"net/http"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type ProxyCaller interface {
	Forward(jsonData *reqData) (string, error)
}

type proxyCaller struct {
	caller ch.CallProxyStruction
}

func NewProxyCaller(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) ProxyCaller {
	return &proxyCaller{
		caller: ch.CallProxyStruction{
			K8sClient:         k8sClient,
			Tracer:            tracer,
			Namespace:         namespace,
			CallerServiceName: callerServiceName,
			CalledServiceName: calledServiceName,
			CalledServicePort: calledServicePort,
		},
	}
}

func (caller *proxyCaller) Forward(jsonData *reqData) (string, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               20000,
		MaxThread:             1000,
		MethodName:            jsonData.Method,
		HttpMethod:            jsonData.HttpMethod,
		EncodeRequestFunc: func(_ context.Context, r *http.Request, request interface{}) error {
			req := request.(string)
			r.Body = ioutil.NopCloser(bytes.NewReader([]byte(req)))
			return nil
		},
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			return string(b), nil
		},
	}

	e, err1 := caller.caller.MakeRemoteEndpoint(p)
	if err1 != nil {
		return "", err1
	}

	resp, err2 := e(context.TODO(), jsonData.Data)
	if err2 != nil {
		return "", err2
	}

	return resp.(string), nil
}

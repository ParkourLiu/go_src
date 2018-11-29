package client

import (
	"context"
	"encoding/json"
	"errors"
	ch "mtcomm/caller/http"
	k8s "mtcomm/k8s"
	"net/http"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type VcodeClient interface {
	CheckCode(*VcodeRequest) (*VcodeResponse, error)
	GetCode(*VcodeRequest) (*VcodeResponse, error)
}

type vcodeClient struct {
	caller ch.CallProxyStruction
}

func NewVcodeClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) VcodeClient {
	return &vcodeClient{
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

type VcodeRequest struct {
	UserId  string `json:"userId"`
	PhoneNo string `json:"phoneNo"`
	Type    string `json:"type"`
	Vcode   string `json:"vcode"`
}

type VcodeResponse struct {
	Err string `json:"err"`
}

func (caller *vcodeClient) GetCode(req *VcodeRequest) (*VcodeResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "getAuthCode",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response VcodeResponse
			if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
				return nil, err
			}
			return response, nil
		},
	}

	e, err1 := caller.caller.MakeRemoteEndpoint(p)
	if err1 != nil {
		return nil, err1
	}

	resp, err2 := e(context.TODO(), req)
	if err2 != nil {
		return nil, err2
	}

	result, _ := resp.(VcodeResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *vcodeClient) CheckCode(req *VcodeRequest) (*VcodeResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "validateAuthCode",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response VcodeResponse
			if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
				return nil, err
			}
			return response, nil
		},
	}

	e, err1 := caller.caller.MakeRemoteEndpoint(p)
	if err1 != nil {
		return nil, err1
	}

	resp, err2 := e(context.TODO(), req)
	if err2 != nil {
		return nil, err2
	}

	result, _ := resp.(VcodeResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

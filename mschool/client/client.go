package client

import (
	"context"
	"encoding/json"
	"errors"
	ch "mtcomm/caller/http"
	k8s "mtcomm/k8s"
	"net/http"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/golang/go/src/pkg/fmt"
)

type MschoolClient interface {
	FaceDataGather() (*Response, error)      //存储活跃用户（周被赞最多）
	LabelGuidDataGather() (*Response, error) //存储分享中用户(周发布个人动态+被赞数最多)
}

type mschoolClient struct {
	caller ch.CallProxyStruction
}

func NewMschoolClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) MschoolClient {
	fmt.Println("callerServiceName:", callerServiceName, "calledServiceName:", calledServiceName)
	return &mschoolClient{
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

type Response struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}

func (caller *mschoolClient) FaceDataGather() (*Response, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "faceDataGather",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response Response
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

	resp, err2 := e(context.TODO(), nil)
	if err2 != nil {
		return nil, err2
	}

	result, _ := resp.(Response)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *mschoolClient) LabelGuidDataGather() (*Response, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "labelGuidDataGather",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response Response
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

	resp, err2 := e(context.TODO(), nil)
	if err2 != nil {
		return nil, err2
	}

	result, _ := resp.(Response)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

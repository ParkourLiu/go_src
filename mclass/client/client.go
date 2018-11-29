package client

import (
	"context"
	"encoding/json"
	"errors"
	ch "mtcomm/caller/http"
	"mtcomm/k8s"
	"net/http"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type MclassClient interface {
	FindTeacherMember(*MclassRequest) (*MteacherResponse, error)
}

type mclassClient struct {
	caller ch.CallProxyStruction
}

func NewLeaveClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) MclassClient {
	return &mclassClient{
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

type MclassRequest struct {
	Tmember  []map[string]string `json:"tmember"`
	ClId     string              `json:"clId"`
	UserId   string              `json:"userId"`
	Relation string              `json:"relation"`
}

type MteacherResponse struct {
	Data []map[string]string `json:"data"`
	Err  string              `json:"err"`
	Code string              `json:"code"`
}

func (caller *mclassClient) FindTeacherMember(req *MclassRequest) (*MteacherResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "findTeacherMember",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response MteacherResponse
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

	result, _ := resp.(MteacherResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

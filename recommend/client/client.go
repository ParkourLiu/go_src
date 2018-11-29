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

type RecommendClient interface {
	//存数据
	SavePopuserUsers() (*recommendResponse, error)    //存储活跃用户（周被赞最多）
	SaveRmduserUsers() (*recommendResponse, error)    //存储分享中用户(周发布个人动态+被赞数最多)
	SaveHotDynamic() (*recommendResponse, error)      //存储最热动态（赞最多）
	SaveFriendRecommend() (*recommendResponse, error) //好友推荐（周发布动态最多）
	PushFavour() (*recommendResponse, error)          //给发动态今天有点赞的用户发推送
	PushStartActivity() (*recommendResponse, error)   //给明天要开始的活动参与者推送通知
	SaveHomePageCache() (*recommendResponse, error)   //首页缓存
	AddDynamicFans() (*recommendResponse, error)      //增加动态的点赞数
}

type recommendClient struct {
	caller ch.CallProxyStruction
}

func NewRecommendClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) RecommendClient {
	return &recommendClient{
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

type recommendResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}

func (caller *recommendClient) SavePopuserUsers() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "savePopuserUsers",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) SaveRmduserUsers() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "saveRmduserUsers",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) SaveHotDynamic() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "saveHotDynamic",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) SaveFriendRecommend() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "saveFriendRecommend",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) PushFavour() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "pushFavour",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) PushStartActivity() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "pushStartActivity",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) SaveHomePageCache() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "saveHomePageCache",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *recommendClient) AddDynamicFans() (*recommendResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "addDynamicFans",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response recommendResponse
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

	result, _ := resp.(recommendResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

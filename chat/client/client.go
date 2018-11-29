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

type ChatClient interface {
	//GetRYToken(info *chat.ToKenInfo) (*DefaultDataResponse, error)
	//GetUserInfo(info *chat.ToKenInfo) (*DefaultDataResponse, error)
	CreateGroupChat(info *GroupChatInfo) (*DefaultResponse, error)
	JoinGroupChat(info *GroupChatInfo) (*DefaultResponse, error)
	QuitGroupChat(info *GroupChatInfo) (*DefaultResponse, error)
	//QueryGroupChatMemberList(info *chat.GroupChatInfo) (*ChatResponse, error)
	//GetArrayGroupInfo(info *chat.GroupChatInfo) (*DefaultDataResponse, error)
	SearchChatInfo(info *GroupChatInfo) (*ChatInfoResponse, error)
	UpdateGroupChat(info *GroupChatInfo)(*DefaultDataResponse,error)
}

type chatClient struct {
	caller ch.CallProxyStruction
}

func NewChatClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) ChatClient {
	return &chatClient{
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

type DefaultResponse struct {
	Err  string `json:"err"`
	Code string `json:"code"`
}
type DefaultDataResponse struct {
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
	Code string      `json:"code"`
}
type ChatResponse struct {
	Err  string                 `json:"err"`
	Data map[string]interface{} `json:"data"`
	Code string                 `json:"code"`
}
type ChatInfoResponse struct {
	Data map[string]string `json:"data"`
	Err  string            `json:"err"`
	Code string            `json:"code"`
}

type GroupChatInfo struct {
	Type             string   `json:"type"`
	UuId             string   `json:"uuId"`
	Token            string   `json:"token"`
	CreateUserId     string   `json:"createUserId"`     //创建者的userid
	UserId           string   `json:"UserId"`           //用户id
	GroupChatId      string   `json:"groupChatId"`      // 群聊的id
	GroupChatIdArray []string `json:"groupChatIdArray"` // 群聊的id数组
	GroupChatName    string   `json:"groupChatName"`    //（活动的名字）
	GroupChatUrl     string   `json:"groupChatUrl"`     //群聊的头像  默认是活动的头像
	GroupChatNotice  string   `json:"groupChatNotice"`  //群聊公告
	UserIdArray      []string `json:"userIdArray"`      // 批量加入群聊的用户id数组
	LastUserId       string   `json:"lastUserId"`       //最后一条 用户id
	AcId             string   `json:"acId"`             //活动的id
	GId              string   `json:"gId"`              //社群id
}

func (caller *chatClient) CreateGroupChat(req *GroupChatInfo) (*DefaultResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "createGroupChat",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response DefaultResponse
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

	result, _ := resp.(DefaultResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *chatClient) QuitGroupChat(req *GroupChatInfo) (*DefaultResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "quitGroupChat",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response DefaultResponse
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

	result, _ := resp.(DefaultResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *chatClient) JoinGroupChat(req *GroupChatInfo) (*DefaultResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "joinGroupChat",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response DefaultResponse
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

	result, _ := resp.(DefaultResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

func (caller *chatClient) SearchChatInfo(req *GroupChatInfo) (*ChatInfoResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "searchChatInfo",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response ChatInfoResponse
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

	result, _ := resp.(ChatInfoResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}
func (caller *chatClient) UpdateGroupChat(req *GroupChatInfo) (*DefaultDataResponse, error) {
	p := &ch.CallerParameter{
		ErrorPercentThreshold: 50,
		Timeout:               5000,
		MaxThread:             500,
		MethodName:            "updateGroupChat",
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			var response DefaultDataResponse
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

	result, _ := resp.(DefaultDataResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}
	return &result, nil
}

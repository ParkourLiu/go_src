package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

//parkour======================================================start
func makeAddOfficialMSGEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(Chat)
		code, err := svc.AddOfficialMSG(&req)
		if err != nil {
			return OfficialMSGResponse{Err: err.Error(), Code: code}, nil
		}
		return OfficialMSGResponse{Err: "", Code: code}, nil
	}
}
func makeLookOfficialMSGEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		data, code, err := svc.LookOfficialMSG()
		if err != nil {
			return OfficialMSGResponse{Err: err.Error(), Code: code}, nil
		}
		return OfficialMSGResponse{Err: "", Code: code, Data: data}, nil
	}
}
func makeInformChatEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(Chat)
		code, err := svc.InformChat(&req)
		if err != nil {
			return OfficialMSGResponse{Err: err.Error(), Code: code}, nil
		}
		return OfficialMSGResponse{Err: "", Code: code}, nil
	}
}

//parkour======================================================end

func makeGetRYTokenEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(ToKenInfo)
		data, code, err := svc.GetRYToken(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error(), Code: code}, nil
		}
		return DefaultDataResponse{Err: "", Code: code, Data: data}, nil
	}
}
func makeGetUserInfoEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(ToKenInfo)
		data, code, err := svc.GetUserInfo(&req)
		if err != nil {
			return DefaultDataResponse{Err: err.Error(), Data: "", Code: code}, nil
		}
		return DefaultDataResponse{Err: "", Code: code, Data: data}, nil
	}
}
func makeGetArrayGroupInfoEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.GetArrayGroupInfo(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error()}, nil
		}
		return DefaultDataResponse{Err: "", Code: code, Data: data}, nil

	}
}

func makeGetMyGroupChatEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.GetMyGroupChat(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error()}, nil
		}
		return DefaultDataResponse{Err: "", Code: code, Data: data}, nil

	}
}
func makeUpdateGroupChatEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.UpdateGroupChat(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error()}, nil
		}
		return DefaultDataResponse{Err: "", Code: code, Data: data}, nil

	}
}
func makeCreateGroupChatEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.CreateGroupChat(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error(), Code: code}, nil
		}
		return DefaultDataResponse{Err: "", Data: data, Code: code}, nil
	}
}
func makeJoinGroupChatEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		code, err := svc.JoinGroupChat(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error(), Code: code}, nil
		}
		return DefaultResponse{Err: "", Code: code}, nil
	}
}
func makeQuitGroupChatEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		code, err := svc.QuitGroupChat(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error(), Code: code}, nil
		}
		return DefaultResponse{Err: "", Code: code}, nil
	}
}
func makeDismissEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		code, err := svc.Dismiss(&req)
		if err != nil {
			return DefaultResponse{Err: err.Error(), Code: code}, nil
		}
		return DefaultResponse{Err: "", Code: code}, nil
	}
}

func makeQueryGroupChatMemberListEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.QueryGroupChatMemberList(&req)
		if err != nil {
			return GroupChatInfoResponse{Err: err.Error(), Data: data, Code: code}, nil
		}
		return GroupChatInfoResponse{Err: "", Data: data, Code: code}, nil
	}
}
func makeSearchChatInfoEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.SearchChatInfo(&req)
		if err != nil {
			return ClassChatInfoReponse{data, err.Error(), code}, nil
		}
		return ClassChatInfoReponse{data, "", code}, nil
	}
}
func makeGetClassIdEndpoint(svc RYToKenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//类型断言
		req := request.(GroupChatInfo)
		data, code, err := svc.GetClassId(&req)
		if err != nil {
			return GetClassIdResponse{data, code, err.Error()}, nil
		}
		return GetClassIdResponse{data, code, ""}, nil
	}
}

//解码 把json字符串为结构体转
func decodeGroupChatInfoquest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GroupChatInfo
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return DefaultResponse{Err: err.Error()}, nil
	}
	return request, nil
}

//解码 把json字符串为结构体转
func decodeDefaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ToKenInfo
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return DefaultResponse{Err: err.Error()}, nil
	}
	return request, nil
}

//编码 把结构体转为json字符串
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type ToKenInfo struct {
	UserId      string   `json:"userId"`      //用户id
	Status      string   `json:"status"`      //token失效状态，  字符串的0表示有效，1表示无效
	UserIdArray []string `json:"userIdArray"` //userid列表
}
type GroupChatInfo struct {
	Type             string   `json:"type"`
	UuId             string   `json:"uuId"`
	Token            string   `json:"token"`
	CreateUserId     string   `json:"createUserId"`     //创建者的userid
	UserId           string   `json:"userId"`           //用户id
	GroupChatId      string   `json:"groupChatId"`      // 群聊的id
	GroupChatIdArray []string `json:"groupChatIdArray"` // 群聊的id数组
	GroupChatName    string   `json:"groupChatName"`    //（活动的名字）
	GroupChatUrl     string   `json:"groupChatUrl"`     //群聊的头像  默认是活动的头像
	GroupChatNotice  string   `json:"groupChatNotice"`  //群聊公告
	UserIdArray      []string `json:"userIdArray"`      // 批量加入群聊的用户id数组
	LastUserId       string   `json:"lastUserId"`       //最后一条 用户id
	AcId             string   `json:"acId"`             //活动的id
	GId              string   `json:"gId"`              //社群id
	ClId             string   `json:"clId"`             //班级id
}
type DefaultResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}
type DefaultDataResponse struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Err  string      `json:"err"`
}
type GroupChatInfoResponse struct {
	Code string                 `json:"code"`
	Data map[string]interface{} `json:"data"`
	Err  string                 `json:"err"`
}
type panicResponse struct {
	Err string   `json:"err"`
	Ids []string `json:"ids"`
}
type ClassChatInfoReponse struct {
	Data map[string]string `json:"data"`
	Err  string            `json:"err"`
	Code string            `json:"code"`
}
type GetClassIdResponse struct {
	ClId string `json:"clId"`
	Code string `json:"code"`
	Err  string `json:err`
}

//parkour======================================================start

//解码 把json字符串为结构体转
func decodeOfficialMSGRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Chat
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return DefaultResponse{Err: err.Error()}, nil
	}
	return request, nil
}

type Chat struct {
	//官方消息
	OfficialMSGId string `json:"officialMSGId"` //官方消息id
	Title         string `json:"title"`         //官方消息标题
	Content       string `json:"content"`       //官方消息内容
	Url           string `json:"url"`           //官方消息url
	//举报群聊
	InformChatId  string `json:"informChatId"`  //举报ID
	GroupChatId   string `json:"groupChatId"`   // 群聊的id
	InformUserId  string `json:"informUserId"`  // 举报群聊的用户ID
	InformExplain string `json:"informExplain"` //举报说明
	InformImg     string `json:"informImg"`     //举报图片
}

type OfficialMSGResponse struct {
	Code string              `json:"code"`
	Err  string              `json:"err"`
	Data []map[string]string `json:"data"`
}

//parkour======================================================end

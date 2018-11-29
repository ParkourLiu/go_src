package main

import (
	"encoding/json"
	"net/http"

	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeStarDetailsEndpoint(svc AutismService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Autism)
		returnMap, code, err := svc.StarDetails(&Autism{LastCoid: req.LastCoid, ClickUserId: req.ClickUserId, St_id: req.St_id})
		if err != nil {
			return StarDetailsResponse{returnMap, code, err.Error()}, nil
		}
		return StarDetailsResponse{returnMap, code, ""}, nil
	}
}

func makeSaveCommentEndpoint(svc AutismService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Autism)
		code, err := svc.SaveComment(&Autism{ClickUserId: req.ClickUserId, St_id: req.St_id, Comment: req.Comment})
		if err != nil {
			return StarDetailsResponse{nil, code, err.Error()}, nil
		}
		return StarDetailsResponse{nil, code, ""}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Autism
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type StarDetailsResponse struct {
	ReturnMap map[string]interface{} `json:"returnMap"`
	Code      string                 `json:"code"`
	Err       string                 `json:"err"`
}

//========================================================================================================
func autismStarListEndpoint(svc AutismService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		u := &Autism{}
		brightCount, likeCount, rolls, starList, err := svc.StarList(u)
		if err != nil {
			return StarListResponse{"101", err.Error(), "0", "0", nil, nil}, nil
		}
		return StarListResponse{"100", "", brightCount, likeCount, rolls, starList}, nil
	}
}

func autismLikesEndpoint(svc AutismService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(likeRequest)
		u := &Autism{St_id: req.St_id, User_id: req.User_id}
		err := svc.Likes(u)
		if err != nil {
			return autismResponse{"101", err.Error()}, nil
		}
		return autismResponse{"100", ""}, nil
	}
}

func autismGetUnionidEndpoint(svc AutismService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUnionid)
		u := &Autism{Code: req.Code}
		data, err := svc.GetUnionid(u)
		if err != nil {
			return GetUniconidResponse{Code:"101",Err:err.Error()}, nil
		}
		return GetUniconidResponse{Code:"100",Err:"",Unionid:data.Unionid,Gender:data.Gender,Icon_url:data.Headimgurl,Name:data.Nickname}, nil
	}
}

func decodeStarListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request starListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func decodeLikesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request likeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func decodeGetUnionidRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getUnionid
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type autismResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}
type likeRequest struct {
	St_id   string `json:"st_id"`
	User_id string `json:"user_id"`
}

type getUnionid struct {
	Code string `json:"code"`
}

type starListRequest struct {
}

/*brightCount,likeCount,rolls,starList, err*/
type StarListResponse struct {
	Code        string              `json:"code"`
	Err         string              `json:"err"`
	BrightCount string              `json:"brightCount"`
	LikeCount   string              `json:"likeCount"`
	Rolls       []map[string]string `json:"rolls"`
	StarList    []map[string]string `json:"starList"`
}
type GetUniconidResponse struct {
	Code    string `json:"code"`
	Err     string `json:"err"`
	Unionid string `json:"unionid"`
	Gender  string `json:"gender"`
	Icon_url string `json:"icon_url"`
	Name   string `json:"name"`
}

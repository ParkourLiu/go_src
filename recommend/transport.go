package main

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makePopuserUsersEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PageInfo)
		u := &PageInfo{LastId: req.LastId, PageSize: req.PageSize, LookUserId: req.LookUserId}
		data, code, err := svc.PopuserUsers(u)
		if err != nil {
			return UserResponse{"", err.Error(), "", nil}, nil
		}
		return UserResponse{code, "", "", data}, nil
	}
}

func makeRmduserUsersEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, code, err := svc.RmduserUsers()
		if err != nil {
			return UserResponse{"", err.Error(), "", nil}, nil
		}
		return UserResponse{code, "", "", data}, nil
	}
}

func makeHotDynamicEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(PageInfo)
		u := &PageInfo{LastId: req.LastId, PageSize: req.PageSize, LookUserId: req.LookUserId}
		data, pageFlag, code, err := svc.HotDynamic(u)
		if err != nil {
			return Response{"", err.Error(), pageFlag, nil}, nil
		}
		return Response{code, "", pageFlag, data}, nil
	}
}

func makeNewDynamicEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(PageInfo)
		u := &PageInfo{LastId: req.LastId, PageSize: req.PageSize, LookUserId: req.LookUserId}
		data, pageFlag, code, err := svc.NewDynamic(u)
		if err != nil {
			return Response{"", err.Error(), pageFlag, nil}, nil
		}
		return Response{code, "", pageFlag, data}, nil
	}
}

func makeReverseHotDynamicEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(PageInfo)
		u := &PageInfo{LastId: req.LastId, PageSize: req.PageSize, LookUserId: req.LookUserId}
		data, pageFlag, code, err := svc.ReverseHotDynamic(u)
		if err != nil {
			return Response{"", err.Error(), pageFlag, nil}, nil
		}
		return Response{code, "", pageFlag, data}, nil
	}
}

func makeReverseNewDynamicEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(PageInfo)
		u := &PageInfo{LastId: req.LastId, PageSize: req.PageSize, LookUserId: req.LookUserId}
		data, pageFlag, code, err := svc.ReverseNewDynamic(u)
		if err != nil {
			return Response{"", err.Error(), pageFlag, nil}, nil
		}
		return Response{code, "", pageFlag, data}, nil
	}
}

func makeSavePopuserUsersEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.SavePopuserUsers()
		if err != nil {
			return Response{"", err.Error(), "", nil}, nil
		}
		return Response{code, "", "", nil}, nil
	}
}

func makeSaveRmduserUsersEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.SaveRmduserUsers()
		if err != nil {
			return Response{"", err.Error(), "", nil}, nil
		}
		return Response{code, "", "", nil}, nil
	}
}

func makeSaveHotDynamicEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.SaveHotDynamic()
		if err != nil {
			return Response{"", err.Error(), "", nil}, nil
		}
		return Response{code, "", "", nil}, nil
	}
}

func makeFriendRecommendEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PageInfo)
		u := &PageInfo{LookUserId: req.LookUserId}
		data, code, err := svc.FriendRecommend(u)
		if err != nil {
			return Response{"", err.Error(), "", nil}, nil
		}
		return Response{code, "", "", data}, nil
	}
}

func makeSaveFriendRecommendEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.SaveFriendRecommend()
		if err != nil {
			return Response{"", err.Error(), "", nil}, nil
		}
		return Response{code, "", "", nil}, nil
	}
}

func makeSearchRecommendEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		data, code, err := svc.SearchRecommend()
		if err != nil {
			return SearchResponse{"", err.Error(), data}, nil
		}
		return SearchResponse{code, "", data}, nil
	}
}

func makePushFavourEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.PushFavour()
		if err != nil {
			return SearchResponse{"", err.Error(), nil}, nil
		}
		return SearchResponse{code, "", nil}, nil
	}
}

func makePushStartActivityEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.PushStartActivity()
		if err != nil {
			return SearchResponse{"", err.Error(), nil}, nil
		}
		return SearchResponse{code, "", nil}, nil
	}
}

func makeSaveHomePageCacheEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.SaveHomePageCache()
		if err != nil {
			return SearchResponse{"", err.Error(), nil}, nil
		}
		return SearchResponse{code, "", nil}, nil
	}
}

func makeAddDynamicFansEndpoint(svc Recommend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		code, err := svc.AddDynamicFans()
		if err != nil {
			return SearchResponse{"", err.Error(), nil}, nil
		}
		return SearchResponse{code, "", nil}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request PageInfo
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type UserResponse struct {
	Code     string              `json:"code"`
	Err      string              `json:"err"`
	PageFlag string              `json:"pageFlag"`
	Data     []map[string]string `json:"data"`
}

type Response struct {
	Code     string                   `json:"code"`
	Err      string                   `json:"err"`
	PageFlag string                   `json:"pageFlag"`
	Data     []map[string]interface{} `json:"data"`
}

type SearchResponse struct {
	Code string   `json:"code"`
	Err  string   `json:"err"`
	Data []string `json:"data"`
}

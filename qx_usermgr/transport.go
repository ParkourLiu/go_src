package main

import (
	"encoding/json"
	"net/http"
	"context"

	U "qx_user/caller"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeRegAndLoginEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(U.UserRequest)
		data, code, err := svc.RegAndLogin(&req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code, Err: ""}, nil
	}
}
func makeOtherLoginEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(U.UserRequest)
		data, code, err := svc.OtherLogin(&req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code, Err: ""}, nil
	}
}
func makeSearchUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(U.UserRequest)
		data, code, err := svc.SearchUser(&req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code, Err: ""}, nil
	}
}
func makeUpdateUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(U.UserRequest)
		data, code, err := svc.UpdateUser(&req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code, Err: ""}, nil
	}
}
func makeChangeBindEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(U.UserRequest)
		data, code, err := svc.ChangeBind(&req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code, Err: ""}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request U.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type Response struct {
	Code string                 `json:"code"`
	Err  string                 `json:"err"`
	Data map[string]interface{} `json:"data"`
}

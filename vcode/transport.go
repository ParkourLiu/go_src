package main

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func validateAuthCodeEndpoint(svc CodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateAuthCodeRequest)
		code, err := svc.CheckCode(Code{UserId: req.UserId, PhoneNo: req.PhoneNo, Type: req.Type, Vcode: req.Vcode})
		if err != nil {
			return validateAuthCodeResponse{code, err.Error()}, nil
		}
		return validateAuthCodeResponse{code, ""}, nil
	}
}
func getAuthCodeEndpoint(svc CodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAuthCodeRequest)
		code, err := svc.GetCode(Code{UserId: req.UserId, PhoneNo: req.PhoneNo, Type: req.Type, Vcode: ""})
		if err != nil {
			return getCodeResponse{code, err.Error()}, nil
		}
		return getCodeResponse{code, ""}, nil
	}
}
func decodeGetCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getAuthCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}
func decodeValidateAuthCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request validateAuthCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type validateAuthCodeRequest struct {
	UserId  string `json:"userId"`
	PhoneNo string `json:"phoneNo"`
	Type    string `json:"type"`
	Vcode   string `json:"vcode"`
}
type getAuthCodeRequest struct {
	UserId  string `json:"userId"`
	PhoneNo string `json:"phoneNo"`
	Type    string `json:"type"`
}

type getCodeResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}
type validateAuthCodeResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}

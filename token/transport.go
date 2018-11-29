package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeCreateTokenEndpoint(svc TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(tokenRequest)
		tk, err := svc.CreateToken(&Token{PlatForm: req.PlatForm, UserId: req.UserId, UUID: req.UUID})
		if err != nil {
			if err.Error() == "Parameter Check Error" {
				return CreateTokenResponse{tk, "", serviceName + "102"}, nil
			}
			return CreateTokenResponse{"", err.Error(), ""}, nil
		}
		return CreateTokenResponse{tk, "", "100"}, nil
	}
}

func makeDeleteTokenEndpoint(svc TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(tokenRequest)
		err := svc.DeleteToken(&Token{PlatForm: req.PlatForm, UUID: req.UUID, Token: req.Token})
		if err != nil {
			if err.Error() == "Parameter Check Error" {
				return DeleteTokenResponse{"", serviceName + "102"}, nil
			}
			return DeleteTokenResponse{err.Error(), ""}, nil
		}
		return DeleteTokenResponse{"", "100"}, nil
	}
}

func decodeCreateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDeleteTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type tokenRequest struct {
	PlatForm string `json:"platForm"`
	UserId   string `json:"userId"`
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
}

type CreateTokenResponse struct {
	Token string `json:"token"`
	Err   string `json:"err"`
	Code  string `json:"code"`
}

type DeleteTokenResponse struct {
	Err  string `json:"err"`
	Code string `json:"code"`
}

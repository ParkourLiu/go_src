package main

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeVersionInfoEndpoint(svc VersionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Version)
		requestList, code, err := svc.VersionInfo(&Version{Flag: req.Flag, NewVersion: req.NewVersion})
		if err != nil {
			msg := err.Error()
			return Response{code, string(msg), nil}, nil
		}
		return Response{code, "", requestList}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Version
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type Response struct {
	Code string            `json:"code"`
	Err  string            `json:"err"`
	Data map[string]string `json:"data"`
}

package main

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeGetOssTokenForWebEndpoint(svc Oss) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ossToken, err := svc.GetOssTokenForWeb()
		if err != nil {
			return Response{"", err.Error(), nil}, nil
		}
		return Response{"100", "", ossToken}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Response
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type Response struct {
	Code     string                 `json:"code"`
	Err      string                 `json:"err"`
	OssToken map[string]interface{} `json:"ossToken"`
}

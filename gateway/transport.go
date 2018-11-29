package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeGatewayEndpoint(svc GatewayService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*reqData)
		return svc.Forward(req)
	}
}

func decodeForwardRequest(_ context.Context, req *http.Request) (interface{}, error) {
	array := strings.Split(req.RequestURI, "/")
	if len(array) != 4 {
		return "", errors.New("URL is error.")
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return &reqData{
		Data:              string(b),
		HttpMethod:        req.Method,
		CalledServiceName: array[2],
		Method:            array[3],
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	resp := []byte(response.(string))
	w.Write(resp)
	return nil
}

type reqData struct {
	Data              string
	Method            string
	CalledServiceName string
	HttpMethod        string
}

func (r *reqData) String() string {
	var buff bytes.Buffer
	buff.WriteString("json: ")
	buff.WriteString(r.Data)
	buff.WriteString(", ")
	buff.WriteString("method: ")
	buff.WriteString(r.HttpMethod)
	buff.WriteString(", ")
	buff.WriteString("callMethod: ")
	buff.WriteString(r.Method)
	buff.WriteString(", ")
	buff.WriteString("callService: ")
	buff.WriteString(r.CalledServiceName)
	return buff.String()
}

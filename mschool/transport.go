package main

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeCreateSchoolYearEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		requestList, code, err := svc.CreateSchoolYear(&req) //传入pbid,lastPid,Type,PageSize,LookUserId
		if err != nil {
			msg := err.Error()
			return SchoolResponse{code, string(msg), nil}, nil
		}
		return SchoolResponse{code, "", requestList}, nil
	}
}

func makeSearchSchoolEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		requestList, code, err := svc.SearchSchool(&req) //传入pbid,lastPid,Type,PageSize,LookUserId
		if err != nil {
			msg := err.Error()
			return SchoolResponse{code, string(msg), nil}, nil
		}
		return SchoolResponse{code, "", requestList}, nil
	}
}

func makeMySchoolEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		requestList, code, err := svc.MySchool(&req) //传入pbid,lastPid,Type,PageSize,LookUserId
		if err != nil {
			msg := err.Error()
			return SchoolResponse1{code, string(msg), nil}, nil
		}
		return SchoolResponse1{code, "", requestList}, nil
	}
}

func makeSetWorkDayEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		code, err := svc.SetWorkDay(&req) //传入pbid,lastPid,Type,PageSize,LookUserId
		if err != nil {
			msg := err.Error()
			return Response{code, msg}, nil
		}
		return Response{code, ""}, nil
	}
}
func makeLookWorkDayEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		requestList, code, err := svc.LookWorkDay(&req) //传入pbid,lastPid,Type,PageSize,LookUserId
		if err != nil {
			msg := err.Error()
			return SchoolResponse2{code, string(msg), nil}, nil
		}
		return SchoolResponse2{code, "", requestList}, nil
	}
}

func makeWorkDayEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		requestList, code, err := svc.WorkDay(&req) //传入pbid,lastPid,Type,PageSize,LookUserId
		if err != nil {
			msg := err.Error()
			return SchoolResponse3{code, string(msg), nil}, nil
		}
		return SchoolResponse3{code, "", requestList}, nil
	}
}

func makeUpSchoolYearEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		code, err := svc.UpSchoolYear(&req)
		if err != nil {
			msg := err.Error()
			return Response{code, msg}, nil
		}
		return Response{code, ""}, nil
	}
}

func makeFaceDataGatherEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		code, err := svc.FaceDataGather()
		if err != nil {
			msg := err.Error()
			return Response{code, msg}, nil
		}
		return Response{code, ""}, nil
	}
}
func makeLabelGuidDataGatherEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		code, err := svc.LabelGuidDataGather()
		if err != nil {
			msg := err.Error()
			return Response{code, msg}, nil
		}
		return Response{code, ""}, nil
	}
}

func makeGetDataFileUrlEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		requestMap, code, err := svc.GetDataFileUrl(&req) //scid
		if err != nil {
			msg := err.Error()
			return SchoolResponse4{code, string(msg), nil}, nil
		}
		return SchoolResponse4{code, "", requestMap}, nil
	}
}
func makeDelDataFileEndpoint(svc MschoolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(School)
		code, err := svc.DelDataFile(&req) //scid
		if err != nil {
			msg := err.Error()
			return Response{code, string(msg)}, nil
		}
		return Response{code, ""}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request School
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type SchoolResponse struct {
	Code string                 `json:"code"`
	Err  string                 `json:"err"`
	Data map[string]interface{} `json:"data"`
}

type SchoolResponse1 struct {
	Code string                `json:"code"`
	Err  string                `json:"err"`
	Data [][]map[string]string `json:"data"`
}

type SchoolResponse2 struct {
	Code string                   `json:"code"`
	Err  string                   `json:"err"`
	Data []map[string]interface{} `json:"data"`
}

type SchoolResponse3 struct {
	Code string              `json:"code"`
	Err  string              `json:"err"`
	Data []map[string]string `json:"data"`
}
type SchoolResponse4 struct {
	Code string            `json:"code"`
	Err  string            `json:"err"`
	Data map[string]string `json:"data"`
}

type Response struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}

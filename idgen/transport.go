package main

import (
	"context"
	"idgen/pb"

	"github.com/go-kit/kit/endpoint"
)

func makeGenerateUniqueIdV1Endpoint(svc IdGeneraterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenerateUniqueIdV1Request)
		ids, err := svc.GenerateUniqueIdV1(req.Count)
		if err != nil {
			return GenerateUniqueIdV1Response{Err: err.Error(), Ids: []string{}}, nil
		}
		return GenerateUniqueIdV1Response{Err: "", Ids: ids}, nil
	}
}

func decodeGenerateUniqueIdV1Request(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GenerateUniqueIdV1Request)
	return GenerateUniqueIdV1Request{Count: req.Count}, nil
}

func encodeGenerateUniqueIdV1Response(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GenerateUniqueIdV1Response)
	return &pb.GenerateUniqueIdV1Reply{Err: resp.Err, Ids: resp.Ids}, nil
}

type GenerateUniqueIdV1Request struct {
	Count uint32
}

type GenerateUniqueIdV1Response struct {
	Err string
	Ids []string
}

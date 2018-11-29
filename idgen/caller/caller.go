package caller

import (
	"context"
	"errors"
	"idgen/pb"
	cg "mtcomm/caller/grpc"
	k8s "mtcomm/k8s"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type IdGeneraterCaller interface {
	GenerateUniqueIdV1(req *GenerateUniqueIdV1Request) (*GenerateUniqueIdV1Response, error)
}

type idGeneraterCaller struct {
	caller cg.CallProxyStruction
}

func NewIdGeneraterCaller(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) IdGeneraterCaller {
	return &idGeneraterCaller{
		caller: cg.CallProxyStruction{
			K8sClient:         k8sClient,
			Tracer:            tracer,
			Namespace:         namespace,
			CallerServiceName: callerServiceName,
			CalledServiceName: calledServiceName,
			CalledServicePort: calledServicePort,
			PbServiceName:     "pb.IdGenerater",
		},
	}
}

type GenerateUniqueIdV1Request struct {
	Count uint32
}

type GenerateUniqueIdV1Response struct {
	Err string
	Ids []string
}

func (caller *idGeneraterCaller) GenerateUniqueIdV1(req *GenerateUniqueIdV1Request) (*GenerateUniqueIdV1Response, error) {
	p := &cg.CallerParameter{
		PbMethod:  "GenerateUniqueIdV1",
		GrpcReply: pb.GenerateUniqueIdV1Reply{},
		Enc: func(_ context.Context, response interface{}) (interface{}, error) {
			req := response.(*GenerateUniqueIdV1Request)
			return &pb.GenerateUniqueIdV1Request{Count: req.Count}, nil
		},
		Dec: func(_ context.Context, grpcResp interface{}) (interface{}, error) {
			resp := grpcResp.(*pb.GenerateUniqueIdV1Reply)
			return &GenerateUniqueIdV1Response{Ids: resp.Ids, Err: resp.Err}, nil
		},
	}

	e, err1 := caller.caller.MakeRemoteEndpoint(p)
	if err1 != nil {
		return nil, err1
	}

	resp, err2 := e(context.TODO(), req)
	if err2 != nil {
		return nil, err2
	}

	result, _ := resp.(*GenerateUniqueIdV1Response)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}

	return result, nil
}

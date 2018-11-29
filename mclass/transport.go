package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeSelectByScInviteCodeEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		schoolYear,err, code := svc.SelectByScInviteCode(&SchoolYear{ScInviteCode: req.ScInviteCode, Flag: req.Flag,ClInviteCode:req.ClInviteCode})
		if err != nil {
			return SchoolYearResponse{nil, err.Error(), code}, nil
		}
		return SchoolYearResponse{schoolYear, "", code}, nil
	}

}
func makeAddClassEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		data, err, code := svc.CreateClass(&SchoolYear{ClId: req.ClId, SyId: req.SyId, GrId: req.GrId, ClName: req.ClName, DirectorUserId: req.DirectorUserId, ClInviteCode: req.ClInviteCode, ClInviteQRCode: req.ClInviteQRCode, Relation: req.Relation})
		if err != nil {
			return AddClassResponse{nil, err.Error(), code}, nil
		}
		return AddClassResponse{data, "", code}, nil
	}
}

func makeTeacherJoinClassEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.TeacherJoinClass(&SchoolYear{ClId:req.ClId,DirectorUserId:req.DirectorUserId,Relation:req.Relation})
		if err != nil {
			return Response{Code:code,Err:err.Error()}, nil
		}
		return Response{Code:code,Err:""}, nil
	}

}
func makeFamilyJoinClassEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.FamilyJoinClass(&req)
		if err != nil {
			return Response{Code:code,Err:err.Error()}, nil
		}
		return Response{Code:code,Err:""}, nil
	}

}
func makeNewMemberEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		data, err, code := svc.NewMember(&req)
		if err != nil {
			return GroupResponse{nil, err.Error(), code}, nil
		}
		return GroupResponse{data, "", code}, nil
	}

}
func makeApproveMembersEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.ApproveMembers(&req)
		if err != nil {
			return Response{code, err.Error()}, nil
		}
		return Response{code, ""}, nil
	}

}
func makeManagerMemberEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		data, err, code := svc.ManagerMember(&SchoolYear{ClId: req.ClId, Flag: req.Flag})
		if err != nil {
			return GroupResponse{nil, err.Error(), ""}, nil
		}
		return GroupResponse{data, "", code}, nil
	}

}
func makeOperateMemberEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.OperateMember(&SchoolYear{ClId:req.ClId, Flag: req.Flag, UserId: req.UserId, Type: req.Type, DisposeUserId: req.DisposeUserId,StId:req.StId})
		if err != nil {
			return Response{"", err.Error()}, nil
		}
		return Response{code, ""}, nil
	}

}
func makeClassQrCodeEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		data, err, code := svc.ClassQrCode(&SchoolYear{ClId: req.ClId,Flag:req.Flag})
		if err != nil {
			return AddClassResponse{nil, err.Error(), ""}, nil
		}
		return AddClassResponse{data, "", code}, nil
	}

}
func makeUpdateStudentEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.UpdateStudent(&req)
		if err != nil {
			return Response{code, err.Error()}, nil
		}
		return Response{code, ""}, nil
	}

}
func makeFindAllMemberEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		data,err, code := svc.FindAllMember(&req)
		if err != nil {
			return GroupResponse{data,err.Error(),code}, nil
		}
		return GroupResponse{data, "",code}, nil
	}

}
func makeFindTeacherMemberEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		data,err, code := svc.FindTeacherMember(&req)
		if err != nil {
			return GroupResponse{data,err.Error(),code}, nil
		}
		return GroupResponse{data, "",code}, nil
	}

}
func makeUpdateGroupChatInfoEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.UpdateGroupChatInfo(&req)
		if err != nil {
			return Response{code,err.Error()}, nil
		}
		return Response{ code,""}, nil
	}

}
func makeUpdateTeachInfoEndpoint(svc McreateClassService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SchoolYear)
		err, code := svc.UpdateTeachInfo(&req)
		if err != nil {
			return Response{code,err.Error()}, nil
		}
		return Response{ code,""}, nil
	}

}
func decodeSelectByScInviteCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SchoolYear
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type SchoolYearResponse struct {
	Data map[string]interface{} `json:"data"`
	Err  string            `json:"err"`
	Code string            `json:"code"`
}
type AddClassResponse struct {
	ClassMessage map[string]string `json:"classMessage"`
	Err          string            `json:"err"`
	Code         string            `json:"code"`
}
type GroupResponse struct {
	Data []map[string]string `json:"data"`
	Err  string              `json:"err"`
	Code string              `json:"code"`
}
type Response struct {
	Code string `json:"code"`
	Err      string            `json:"err"`
}

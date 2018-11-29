package main

import (
	"context"
	"qx_user/pb"
	"github.com/go-kit/kit/endpoint"
)

func makeSearchUserByIdEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*User)
		data, code, err := svc.SearchUserById(req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code}, nil
	}
}

func makeSearchUsersEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*User)
		data, code, err := svc.SearchUsers(req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code}, nil
	}
}

func makeAddUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*User)
		data, code, err := svc.AddUser(req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code}, nil
	}
}

func makeUpdateUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*User)
		data, code, err := svc.UpdateUser(req)
		if err != nil {
			return Response{Data: data, Code: code, Err: err.Error()}, nil
		}
		return Response{Data: data, Code: code}, nil
	}
}

func decodeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserRequest)
	user := &User{
		UserId:         req.UserId,
		PhoneNo:        req.PhoneNo,
		Password:       req.Password,
		Email:          req.Email,
		TrueName:       req.TrueName,
		NickName:       req.NickName,
		BirthDay:       req.BirthDay,
		ChineseZodiac:  req.ChineseZodiac,
		QrImageName:    req.QrImageName,
		Sex:            req.Sex,
		HomeAddress:    req.HomeAddress,
		ImageName:      req.ImageName,
		ChatName:       req.ChatName,
		ChatPwd:        req.ChatPwd,
		MtalkNo:        req.MtalkNo,
		Hometown:       req.Hometown,
		Description:    req.Description,
		PlatForm:       req.PlatForm,
		UUID:           req.UUID,
		OpenId:         req.OpenId,
		BackgroundImg:  req.BackgroundImg,
		Wechat_uid:     req.WechatUid,
		Wechat_name:    req.WechatName,
		Wechat_iconurl: req.WechatIconurl,
		Wechat_gender:  req.WechatGender,
		QQ_uid:         req.QQUid,
		QQ_name:        req.QQName,
		QQ_iconurl:     req.QQIconurl,
		QQ_gender:      req.QQGender,
		IsWater:        req.IsWater,
		Volunteer:      req.Volunteer,
	}
	return user, nil
}

func encodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(Response)

	if resp.Data["user"] != nil {
		userRequest := &pb.UserRequest{
			UserId:        resp.Data["user"].(map[string]string)["userId"],
			PhoneNo:       resp.Data["user"].(map[string]string)["phoneNo"],
			Password:      resp.Data["user"].(map[string]string)["password"],
			Email:         resp.Data["user"].(map[string]string)["email"],
			TrueName:      resp.Data["user"].(map[string]string)["trueName"],
			NickName:      resp.Data["user"].(map[string]string)["nickName"],
			BirthDay:      resp.Data["user"].(map[string]string)["birthDay"],
			ChineseZodiac: resp.Data["user"].(map[string]string)["chineseZodiac"],
			QrImageName:   resp.Data["user"].(map[string]string)["qrImageName"],
			Sex:           resp.Data["user"].(map[string]string)["sex"],
			HomeAddress:   resp.Data["user"].(map[string]string)["homeAddress"],
			ImageName:     resp.Data["user"].(map[string]string)["imageName"],
			ChatName:      resp.Data["user"].(map[string]string)["chatName"],
			ChatPwd:       resp.Data["user"].(map[string]string)["chatPwd"],
			MtalkNo:       resp.Data["user"].(map[string]string)["mtalkNo"],
			Hometown:      resp.Data["user"].(map[string]string)["hometown"],
			Description:   resp.Data["user"].(map[string]string)["description"],
			PlatForm:      resp.Data["user"].(map[string]string)["platForm"],
			UUID:          resp.Data["user"].(map[string]string)["UUID"],
			OpenId:        resp.Data["user"].(map[string]string)["openId"],
			BackgroundImg: resp.Data["user"].(map[string]string)["backgroundImg"],
			WechatUid:     resp.Data["user"].(map[string]string)["Wechat_uid"],
			WechatName:    resp.Data["user"].(map[string]string)["Wechat_name"],
			WechatIconurl: resp.Data["user"].(map[string]string)["Wechat_iconurl"],
			WechatGender:  resp.Data["user"].(map[string]string)["Wechat_gender"],
			QQUid:         resp.Data["user"].(map[string]string)["QQ_uid"],
			QQName:        resp.Data["user"].(map[string]string)["QQ_name"],
			QQIconurl:     resp.Data["user"].(map[string]string)["QQ_iconurl"],
			QQGender:      resp.Data["user"].(map[string]string)["QQ_gender"],
			IsWater:       resp.Data["user"].(map[string]string)["isWater"],
			Volunteer:     resp.Data["user"].(map[string]string)["volunteer"],
		}
		returnMap := map[string]*pb.UserRequest{}
		returnMap["user"] = userRequest
		return &pb.UserReply{Data: returnMap, Code: resp.Code, Err: resp.Err}, nil

	} else if resp.Data["users"] != nil {
		returnMap := map[string]*pb.Users{}
		Users := &pb.Users{}
		reqList := resp.Data["users"].([]map[string]string)
		for _, v := range reqList {
			UserRequest := &pb.UserRequest{
				UserId:        v["userId"],
				PhoneNo:       v["phoneNo"],
				Password:      v["password"],
				Email:         v["email"],
				TrueName:      v["trueName"],
				NickName:      v["nickName"],
				BirthDay:      v["birthDay"],
				ChineseZodiac: v["chineseZodiac"],
				QrImageName:   v["qrImageName"],
				Sex:           v["sex"],
				HomeAddress:   v["homeAddress"],
				ImageName:     v["imageName"],
				ChatName:      v["chatName"],
				ChatPwd:       v["chatPwd"],
				MtalkNo:       v["mtalkNo"],
				Hometown:      v["hometown"],
				Description:   v["description"],
				PlatForm:      v["platForm"],
				UUID:          v["UUID"],
				OpenId:        v["openId"],
				BackgroundImg: v["backgroundImg"],
				WechatUid:     v["Wechat_uid"],
				WechatName:    v["Wechat_name"],
				WechatIconurl: v["Wechat_iconurl"],
				WechatGender:  v["Wechat_gender"],
				QQUid:         v["QQ_uid"],
				QQName:        v["QQ_name"],
				QQIconurl:     v["QQ_iconurl"],
				QQGender:      v["QQ_gender"],
				IsWater:       v["isWater"],
				Volunteer:     v["volunteer"],
			}
			Users.Users = append(Users.Users, UserRequest)
		}
		returnMap["users"] = Users
		return &pb.UsersReply{Data: returnMap, Code: resp.Code, Err: resp.Err}, nil
	} else {
		return &pb.UserReply{Code: resp.Code, Err: resp.Err}, nil
	}
}

type Response struct {
	Data map[string]interface{} `json:"data"`
	Code string                 `json:"code"`
	Err  string                 `json:"err"`
}

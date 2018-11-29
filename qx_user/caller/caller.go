package caller

import (
	"context"
	"errors"
	"qx_user/pb"
	cg "mtcomm/caller/grpc"
	k8s "mtcomm/k8s"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type UserCaller interface {
	SearchUserById(req *UserRequest) (*UserResponse, error)
	SearchUsers(req *UserRequest) (*UsersResponse, error)
	AddUser(req *UserRequest) (*UserResponse, error)
	UpdateUser(req *UserRequest) (*UserResponse, error)
}

type userCaller struct {
	caller cg.CallProxyStruction
}

func NewIdUserCaller(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace, callerServiceName, calledServiceName, calledServicePort string) UserCaller {
	return &userCaller{
		caller: cg.CallProxyStruction{
			K8sClient:         k8sClient,
			Tracer:            tracer,
			Namespace:         namespace,
			CallerServiceName: callerServiceName,
			CalledServiceName: calledServiceName,
			CalledServicePort: calledServicePort,
			PbServiceName:     "pb.UserService",
		},
	}
}

type UserRequest struct {
	UserId         string `json:"userId"`         //主键
	PhoneNo        string `json:"phoneNo"`        //电话
	Password       string `json:"password"`       //密码
	Email          string `json:"email"`          //邮箱
	TrueName       string `json:"trueName"`       //真名
	NickName       string `json:"nickName"`       //昵称
	BirthDay       string `json:"birthDay"`       //生日
	ChineseZodiac  string `json:"chineseZodiac"`  //生肖
	QrImageName    string `json:"qrImageName"`    //用户二维码
	Sex            string `json:"sex"`            //性别
	HomeAddress    string `json:"homeAddress"`    //家庭住址
	ImageName      string `json:"imageName"`      //头像
	ChatName       string `json:"chatName"`       //
	ChatPwd        string `json:"chatPwd"`        //
	MtalkNo        string `json:"mtalkNo"`        //钦家号
	Hometown       string `json:"hometown"`       //
	Description    string `json:"description"`    //
	PlatForm       string `json:"platForm"`       //来源（安卓，ios）
	UUID           string `json:"UUID"`           //设备识别码
	OpenId         string `json:"openId"`         //
	BackgroundImg  string `json:"backgroundImg"`  //
	Wechat_uid     string `json:"Wechat_uid"`     //微信uid
	Wechat_name    string `json:"Wechat_name"`    //微信名
	Wechat_iconurl string `json:"Wechat_iconurl"` //微信头像
	Wechat_gender  string `json:"Wechat_gender"`  //微信性别
	QQ_uid         string `json:"QQ_uid"`         //QQid
	QQ_name        string `json:"QQ_name"`        //qq名
	QQ_iconurl     string `json:"QQ_iconurl"`     //qq头像
	QQ_gender      string `json:"QQ_gender"`      //qq性别
	IsWater        string `json:"isWater"`        //是否是水军
	Volunteer      string `json:"volunteer"`      //是否是志愿者

	Type string `json:"type"` //用于前端传参
}

type UserResponse struct {
	Data UserRequest
	Code string
	Err  string
}

type UsersResponse struct {
	Data []UserRequest
	Code string
	Err  string
}

func (caller *userCaller) SearchUserById(req *UserRequest) (*UserResponse, error) {
	p := &cg.CallerParameter{
		PbMethod:  "SearchUserById",
		GrpcReply: pb.UserReply{},
		Enc: func(_ context.Context, response interface{}) (interface{}, error) {
			req := response.(*UserRequest)
			pb := &pb.UserRequest{
				UserId:        req.UserId,
				PhoneNo:       req.PhoneNo,
				Password:      req.Password,
				Email:         req.Email,
				TrueName:      req.TrueName,
				NickName:      req.NickName,
				BirthDay:      req.BirthDay,
				ChineseZodiac: req.ChineseZodiac,
				QrImageName:   req.QrImageName,
				Sex:           req.Sex,
				HomeAddress:   req.HomeAddress,
				ImageName:     req.ImageName,
				ChatName:      req.ChatName,
				ChatPwd:       req.ChatPwd,
				MtalkNo:       req.MtalkNo,
				Hometown:      req.Hometown,
				Description:   req.Description,
				PlatForm:      req.PlatForm,
				UUID:          req.UUID,
				OpenId:        req.OpenId,
				BackgroundImg: req.BackgroundImg,
				WechatUid:     req.Wechat_uid,
				WechatName:    req.Wechat_name,
				WechatIconurl: req.Wechat_iconurl,
				WechatGender:  req.Wechat_gender,
				QQUid:         req.QQ_uid,
				QQName:        req.QQ_name,
				QQIconurl:     req.QQ_iconurl,
				QQGender:      req.QQ_gender,
				IsWater:       req.IsWater,
				Volunteer:     req.Volunteer,
			}
			return pb, nil
		},
		Dec: func(_ context.Context, grpcResp interface{}) (interface{}, error) {
			resp := grpcResp.(*pb.UserReply)
			u := UserRequest{
				UserId:         resp.Data["user"].UserId,
				PhoneNo:        resp.Data["user"].PhoneNo,
				Password:       resp.Data["user"].Password,
				Email:          resp.Data["user"].Email,
				TrueName:       resp.Data["user"].TrueName,
				NickName:       resp.Data["user"].NickName,
				BirthDay:       resp.Data["user"].BirthDay,
				ChineseZodiac:  resp.Data["user"].ChineseZodiac,
				QrImageName:    resp.Data["user"].QrImageName,
				Sex:            resp.Data["user"].Sex,
				HomeAddress:    resp.Data["user"].HomeAddress,
				ImageName:      resp.Data["user"].ImageName,
				ChatName:       resp.Data["user"].ChatName,
				ChatPwd:        resp.Data["user"].ChatPwd,
				MtalkNo:        resp.Data["user"].MtalkNo,
				Hometown:       resp.Data["user"].Hometown,
				Description:    resp.Data["user"].Description,
				PlatForm:       resp.Data["user"].PlatForm,
				UUID:           resp.Data["user"].UUID,
				OpenId:         resp.Data["user"].OpenId,
				BackgroundImg:  resp.Data["user"].BackgroundImg,
				Wechat_uid:     resp.Data["user"].WechatUid,
				Wechat_name:    resp.Data["user"].WechatName,
				Wechat_iconurl: resp.Data["user"].WechatIconurl,
				Wechat_gender:  resp.Data["user"].WechatGender,
				QQ_uid:         resp.Data["user"].QQUid,
				QQ_name:        resp.Data["user"].QQName,
				QQ_iconurl:     resp.Data["user"].QQIconurl,
				QQ_gender:      resp.Data["user"].QQGender,
				IsWater:        resp.Data["user"].IsWater,
				Volunteer:      resp.Data["user"].Volunteer,
			}
			userResponse := &UserResponse{Data: u, Code: resp.Code, Err: resp.Err}
			return userResponse, nil
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

	result, _ := resp.(*UserResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}

	return result, nil
}

func (caller *userCaller) SearchUsers(req *UserRequest) (*UsersResponse, error) {
	p := &cg.CallerParameter{
		PbMethod:  "SearchUsers",
		GrpcReply: pb.UsersReply{},
		Enc: func(_ context.Context, response interface{}) (interface{}, error) {
			req := response.(*UserRequest)
			pb := &pb.UserRequest{
				UserId:        req.UserId,
				PhoneNo:       req.PhoneNo,
				Password:      req.Password,
				Email:         req.Email,
				TrueName:      req.TrueName,
				NickName:      req.NickName,
				BirthDay:      req.BirthDay,
				ChineseZodiac: req.ChineseZodiac,
				QrImageName:   req.QrImageName,
				Sex:           req.Sex,
				HomeAddress:   req.HomeAddress,
				ImageName:     req.ImageName,
				ChatName:      req.ChatName,
				ChatPwd:       req.ChatPwd,
				MtalkNo:       req.MtalkNo,
				Hometown:      req.Hometown,
				Description:   req.Description,
				PlatForm:      req.PlatForm,
				UUID:          req.UUID,
				OpenId:        req.OpenId,
				BackgroundImg: req.BackgroundImg,
				WechatUid:     req.Wechat_uid,
				WechatName:    req.Wechat_name,
				WechatIconurl: req.Wechat_iconurl,
				WechatGender:  req.Wechat_gender,
				QQUid:         req.QQ_uid,
				QQName:        req.QQ_name,
				QQIconurl:     req.QQ_iconurl,
				QQGender:      req.QQ_gender,
				IsWater:       req.IsWater,
				Volunteer:     req.Volunteer,
			}
			return pb, nil
		},
		Dec: func(_ context.Context, grpcResp interface{}) (interface{}, error) {
			resp := grpcResp.(*pb.UsersReply)
			users := []UserRequest{}
			for _, v := range resp.Data["users"].Users {
				u := UserRequest{
					UserId:         v.UserId,
					PhoneNo:        v.PhoneNo,
					Password:       v.Password,
					Email:          v.Email,
					TrueName:       v.TrueName,
					NickName:       v.NickName,
					BirthDay:       v.BirthDay,
					ChineseZodiac:  v.ChineseZodiac,
					QrImageName:    v.QrImageName,
					Sex:            v.Sex,
					HomeAddress:    v.HomeAddress,
					ImageName:      v.ImageName,
					ChatName:       v.ChatName,
					ChatPwd:        v.ChatPwd,
					MtalkNo:        v.MtalkNo,
					Hometown:       v.Hometown,
					Description:    v.Description,
					PlatForm:       v.PlatForm,
					UUID:           v.UUID,
					OpenId:         v.OpenId,
					BackgroundImg:  v.BackgroundImg,
					Wechat_uid:     v.WechatUid,
					Wechat_name:    v.WechatName,
					Wechat_iconurl: v.WechatIconurl,
					Wechat_gender:  v.WechatGender,
					QQ_uid:         v.QQUid,
					QQ_name:        v.QQName,
					QQ_iconurl:     v.QQIconurl,
					QQ_gender:      v.QQGender,
					IsWater:        v.IsWater,
					Volunteer:      v.Volunteer,
				}
				users = append(users, u)
			}
			usersResponse := &UsersResponse{Data: users, Code: resp.Code, Err: resp.Err}
			return usersResponse, nil
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

	result, _ := resp.(*UsersResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}

	return result, nil
}

func (caller *userCaller) AddUser(req *UserRequest) (*UserResponse, error) {
	p := &cg.CallerParameter{
		PbMethod:  "AddUser",
		GrpcReply: pb.UserReply{},
		Enc: func(_ context.Context, response interface{}) (interface{}, error) {
			req := response.(*UserRequest)
			pb := &pb.UserRequest{
				UserId:        req.UserId,
				PhoneNo:       req.PhoneNo,
				Password:      req.Password,
				Email:         req.Email,
				TrueName:      req.TrueName,
				NickName:      req.NickName,
				BirthDay:      req.BirthDay,
				ChineseZodiac: req.ChineseZodiac,
				QrImageName:   req.QrImageName,
				Sex:           req.Sex,
				HomeAddress:   req.HomeAddress,
				ImageName:     req.ImageName,
				ChatName:      req.ChatName,
				ChatPwd:       req.ChatPwd,
				MtalkNo:       req.MtalkNo,
				Hometown:      req.Hometown,
				Description:   req.Description,
				PlatForm:      req.PlatForm,
				UUID:          req.UUID,
				OpenId:        req.OpenId,
				BackgroundImg: req.BackgroundImg,
				WechatUid:     req.Wechat_uid,
				WechatName:    req.Wechat_name,
				WechatIconurl: req.Wechat_iconurl,
				WechatGender:  req.Wechat_gender,
				QQUid:         req.QQ_uid,
				QQName:        req.QQ_name,
				QQIconurl:     req.QQ_iconurl,
				QQGender:      req.QQ_gender,
				IsWater:       req.IsWater,
				Volunteer:     req.Volunteer,
			}
			return pb, nil
		},
		Dec: func(_ context.Context, grpcResp interface{}) (interface{}, error) {
			resp := grpcResp.(*pb.UserReply)
			userResponse := &UserResponse{Code: resp.Code, Err: resp.Err}
			return userResponse, nil
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

	result, _ := resp.(*UserResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}

	return result, nil
}

func (caller *userCaller) UpdateUser(req *UserRequest) (*UserResponse, error) {
	p := &cg.CallerParameter{
		PbMethod:  "UpdateUser",
		GrpcReply: pb.UserReply{},
		Enc: func(_ context.Context, response interface{}) (interface{}, error) {
			req := response.(*UserRequest)
			pb := &pb.UserRequest{
				UserId:        req.UserId,
				PhoneNo:       req.PhoneNo,
				Password:      req.Password,
				Email:         req.Email,
				TrueName:      req.TrueName,
				NickName:      req.NickName,
				BirthDay:      req.BirthDay,
				ChineseZodiac: req.ChineseZodiac,
				QrImageName:   req.QrImageName,
				Sex:           req.Sex,
				HomeAddress:   req.HomeAddress,
				ImageName:     req.ImageName,
				ChatName:      req.ChatName,
				ChatPwd:       req.ChatPwd,
				MtalkNo:       req.MtalkNo,
				Hometown:      req.Hometown,
				Description:   req.Description,
				PlatForm:      req.PlatForm,
				UUID:          req.UUID,
				OpenId:        req.OpenId,
				BackgroundImg: req.BackgroundImg,
				WechatUid:     req.Wechat_uid,
				WechatName:    req.Wechat_name,
				WechatIconurl: req.Wechat_iconurl,
				WechatGender:  req.Wechat_gender,
				QQUid:         req.QQ_uid,
				QQName:        req.QQ_name,
				QQIconurl:     req.QQ_iconurl,
				QQGender:      req.QQ_gender,
				IsWater:       req.IsWater,
				Volunteer:     req.Volunteer,
			}
			return pb, nil
		},
		Dec: func(_ context.Context, grpcResp interface{}) (interface{}, error) {
			resp := grpcResp.(*pb.UserReply)
			userResponse := &UserResponse{Code: resp.Code, Err: resp.Err}
			return userResponse, nil
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

	result, _ := resp.(*UserResponse)
	if result.Err != "" {
		return nil, errors.New(result.Err)
	}

	return result, nil
}

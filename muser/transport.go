package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeHomePageSloganEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userMap, code, err := svc.HomePageSlogan()
		if err != nil {
			return SloganResponse{code, err.Error(), nil}, nil
		}
		return SloganResponse{code, "", userMap}, nil
	}
}

func makeMyHomeEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(userRequest)
		user := &User{UserId: req.UserId, LookUserId: req.LookUserId}
		log.Debug("█████进入MyHome时间：", fmt.Sprint(time.Now()))
		userMap, helpList, helpPetList, code, err := svc.MyHome(user)
		log.Debug("█████结束MyHome时间：", fmt.Sprint(time.Now()))
		if err != nil {
			return sUserRequest{code, err.Error(), nil, nil, nil}, nil
		}
		return sUserRequest{code, "", userMap, helpList, helpPetList}, nil
	}
}

func makeRegEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
		}
		code, userId, havePassword, err := svc.Reg(user)
		if err != nil {
			return regLoginResponse{code, err.Error(), "", ""}, nil
		}
		return regLoginResponse{code, "", userId, havePassword}, nil
	}
}

func makeLoginEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
		}
		code, userId, havePassword, err := svc.Login(user)
		if err != nil {
			return regLoginResponse{code, err.Error(), "", ""}, nil
		}
		return regLoginResponse{code, "", userId, havePassword}, nil
	}
}

func makeShortcutLoginEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
		}
		code, userId, havePassword, err := svc.ShortcutLogin(user)
		if err != nil {
			return regLoginResponse{code, err.Error(), "", ""}, nil
		}
		return regLoginResponse{code, "", userId, havePassword}, nil
	}
}

func makeFindPasswordEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
		}
		code, userId, err := svc.FindPassword(user)
		if err != nil {
			return userIdResponse{code, err.Error(), "", ""}, nil
		}
		return userIdResponse{code, "", userId, ""}, nil
	}
}

func makeChangePhoneNoEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
		}
		code, err := svc.ChangePhoneNo(user)
		if err != nil {
			return userResponse{code, err.Error()}, nil
		}
		return userResponse{code, ""}, nil
	}
}

func makeOtherLoginEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
		}
		code, userId, phoneNo, havePassword, err := svc.OtherLogin(user)
		if err != nil {
			return otherLoginResponse{code, err.Error(), "", "", ""}, nil
		}
		return otherLoginResponse{code, "", userId, phoneNo, havePassword}, nil
	}
}

func makeUpdateUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(userRequest)
		user := &User{
			PhoneNo:         req.PhoneNo,
			Password:        req.Password,
			Sina_uid:        req.Sina_uid,
			Sina_name:       req.Sina_name,
			Sina_iconurl:    req.Sina_iconurl,
			Sina_gender:     req.Sina_gender,
			Wechat_uid:      req.Wechat_uid,
			Wechat_name:     req.Wechat_name,
			Wechat_iconurl:  req.Wechat_iconurl,
			Wechat_gender:   req.Wechat_gender,
			QQ_uid:          req.QQ_uid,
			QQ_name:         req.QQ_name,
			QQ_iconurl:      req.QQ_iconurl,
			QQ_gender:       req.QQ_gender,
			UserId:          req.UserId,
			TrueName:        req.TrueName,
			NickName:        req.NickName,
			BirthDay:        req.BirthDay,
			ChineseZodiac:   req.ChineseZodiac,
			Sex:             req.Sex,
			HomeAddress:     req.HomeAddress,
			ImageName:       req.ImageName,
			ChatName:        req.ChatName,
			ChatPwd:         req.ChatPwd,
			Hometown:        req.Hometown,
			Description:     req.Description,
			PlatForm:        req.PlatForm,
			OpenId:          req.OpenId,
			BackgroundImg:   req.BackgroundImg,
			OtherLogienType: req.OtherLogienType,
			Id:              req.Id,
		}
		code, err := svc.UpdateUser(user)
		if err != nil {
			return userResponse{code, err.Error()}, nil
		}
		return userResponse{code, ""}, nil
	}
}

func makeCheckPhoneBookEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PhoneBook)
		phoneBook := &PhoneBook{UserId: req.UserId, HashCode: req.HashCode}
		code, err, flag := svc.CheckPhoneBook(phoneBook)
		if err != nil {
			return HashCodeResponse{code, err.Error(), ""}, nil
		}
		return HashCodeResponse{code, "", flag}, nil
	}
}

func makePhoneBookUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PhoneBook)
		phoneBook := &PhoneBook{UserId: req.UserId, Phones: req.Phones}
		code, err, haveUser, noUser := svc.PhoneBookUser(phoneBook)
		if err != nil {
			return PhonesBookResponse{code, err.Error(), nil, nil}, nil
		}
		return PhonesBookResponse{code, "", haveUser, noUser}, nil
	}
}

func makeActiveUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(userRequest)
		code, err := svc.ActiveUser(&User{UserId: req.UserId})
		if err != nil {
			return userResponse{code, err.Error()}, nil
		}
		return userResponse{code, ""}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request userRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodePhonnesBookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request PhoneBook
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(context.TODO(), w, response)
}

type userRequest struct {
	LookUserId      string `json:"lookUserId"`
	UserId          string `json:"userId"`
	PhoneNo         string `json:"phoneNo"`
	Password        string `json:"password"`
	TrueName        string `json:"trueName"`
	NickName        string `json:"nickName"`
	BirthDay        string `json:"birthDay"`
	ChineseZodiac   string `json:"chineseZodiac"`
	Sex             string `json:"sex"`
	HomeAddress     string `json:"homeAddress"`
	ImageName       string `json:"imageName"`
	ChatName        string `json:"chatName"`
	ChatPwd         string `json:"chatPwd"`
	MtalkNo         string `json:"mtalkNo"`
	Hometown        string `json:"hometown"`
	Description     string `json:"description"`
	PlatForm        string `json:"platForm"`
	UUID            string `json:"UUID"`
	OpenId          string `json:"openId"`
	BackgroundImg   string `json:"backgroundImg"`
	Sina_uid        string `json:"Sina_uid"`
	Sina_name       string `json:"Sina_name"`
	Sina_iconurl    string `json:"Sina_iconurl"`
	Sina_gender     string `json:"Sina_gender"`
	Wechat_uid      string `json:"Wechat_uid"`
	Wechat_name     string `json:"Wechat_name"`
	Wechat_iconurl  string `json:"Wechat_iconurl"`
	Wechat_gender   string `json:"Wechat_gender"`
	QQ_uid          string `json:"QQ_uid"`
	QQ_name         string `json:"QQ_name"`
	QQ_iconurl      string `json:"QQ_iconurl"`
	QQ_gender       string `json:"QQ_gender"`
	Ad              string `json:"ad"`
	OtherLogienType string `json:"otherLogienType"` //0微博,1微信，2qq
	Id              string `json:"id"`              //数据导入临时使用，无实际业务意义
}

type userResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
}
type userIdResponse struct {
	Code    string `json:"code"`
	Err     string `json:"err"`
	UserId  string `json:"userId"`
	PhoneNo string `json:"phoneNo"`
}

type otherLoginResponse struct {
	Code         string `json:"code"`
	Err          string `json:"err"`
	UserId       string `json:"userId"`
	PhoneNo      string `json:"phoneNo"`
	HavePassword string `json:"havePassword"`
}

type regLoginResponse struct {
	Code         string `json:"code"`
	Err          string `json:"err"`
	UserId       string `json:"userId"`
	HavePassword string `json:"havePassword"`
}

type sUserRequest struct {
	Code        string            `json:"code"`
	Err         string            `json:"err"`
	User        map[string]string `json:"user"`
	HelpList    []string          `json:"helpList"`
	HelpPetList []string          `json:"helpPetList"`
}

type SloganResponse struct {
	Code   string              `json:"code"`
	Err    string              `json:"err"`
	Slogan []map[string]string `json:"slogan"`
}

type HashCodeResponse struct {
	Code string `json:"code"`
	Err  string `json:"err"`
	Flag string `json:"flag"`
}

type PhonesBookResponse struct {
	Code     string              `json:"code"`
	Err      string              `json:"err"`
	HavaUser []map[string]string `json:"havaUser"`
	NoUser   []map[string]string `json:"noUser"`
}

package main

import (
	"errors"
	U "qx_user/caller"
)

//手机号注册加锁添加用户
func SynRegAddUser(ur *U.UserRequest) (string, error) { //添加方法加锁,返回userId,code,err
	ex, err2 := redisClient.SETNX("muser:P:"+ur.PhoneNo, "1")
	defer redisClient.SetStringAndExpire("muser:P:"+ur.PhoneNo, "1", 1)
	if err2 != nil {
		return "", err2
	}
	if ex != 1 {
		return "", errors.New("请不要重复提交")
	}
	userId := ""
	req, err := userClient.SearchUsers(&U.UserRequest{PhoneNo: ur.PhoneNo}) //查询数据库数据
	if err != nil {
		return "", err
	}
	if req.Err != "" {
		return "", errors.New(req.Err)
	}
	if len(req.Data) == 0 {
		rep, err1 := userClient.AddUser(&U.UserRequest{PhoneNo: ur.PhoneNo})
		if err1 != nil {
			return "", err1
		}
		if rep.Err != "" {
			return "", errors.New(rep.Err)
		}
		userId = rep.Data.UserId
	}
	if userId == "" {
		return "", errors.New("请不要重复提交")
	}
	return userId, nil
}

//第三方注册加锁添加用户
func SynOtherAddUser(ur *U.UserRequest) (string, error) { //添加方法加锁
	userReq := &U.UserRequest{}
	if ur.Type == "1" { //微信
		ex, err2 := redisClient.SETNX("muser:W:"+ur.Wechat_uid, "1")
		defer redisClient.SetStringAndExpire("muser:W:"+ur.Wechat_uid, "1", 1)
		if err2 != nil {
			return "", err2
		}
		if ex != 1 {
			return "", errors.New("请不要重复提交")
		}
		userReq.Wechat_uid = ur.Wechat_uid
		userReq.Wechat_name = ur.Wechat_name
		userReq.Wechat_iconurl = ur.Wechat_iconurl
		userReq.Wechat_gender = ur.Wechat_gender
	} else if ur.Type == "2" { //QQ
		ex, err2 := redisClient.SETNX("muser:Q:"+ur.QQ_uid, "1")
		defer redisClient.SetStringAndExpire("muser:Q:"+ur.QQ_uid, "1", 1)
		if err2 != nil {
			return "", err2
		}
		if ex != 1 {
			return "", errors.New("请不要重复提交")
		}
		userReq.QQ_uid = ur.QQ_uid
		userReq.QQ_name = ur.QQ_name
		userReq.QQ_iconurl = ur.QQ_iconurl
		userReq.QQ_gender = ur.QQ_gender
	}

	userId := ""
	request, err := userClient.SearchUsers(userReq) //查询数据库数据
	if err != nil {
		return "", err
	}
	if request.Err != "" {
		return "", errors.New(request.Err)
	}
	if len(request.Data) == 0 {
		rep, err1 := userClient.AddUser(userReq)
		if err1 != nil {
			return "", err1
		}
		if rep.Err != "" {
			return "", errors.New(rep.Err)
		}
		userId = rep.Data.UserId
	}
	if userId == "" {
		return "", errors.New("请不要重复提交")
	}

	return userId, nil
}

func UserRequest2mapi(ur *U.UserRequest) map[string]interface{} {
	returnMap := map[string]interface{}{}
	returnMap["userId"] = ur.UserId
	returnMap["phoneNo"] = ur.PhoneNo
	returnMap["password"] = ur.Password
	returnMap["email"] = ur.Email
	returnMap["trueName"] = ur.TrueName
	returnMap["nickName"] = ur.NickName
	returnMap["birthDay"] = ur.BirthDay
	returnMap["chineseZodiac"] = ur.ChineseZodiac
	returnMap["qrImageName"] = ur.QrImageName
	returnMap["sex"] = ur.Sex
	returnMap["homeAddress"] = ur.HomeAddress
	returnMap["imageName"] = ur.ImageName
	returnMap["chatName"] = ur.ChatName
	returnMap["chatPwd"] = ur.ChatPwd
	returnMap["mtalkNo"] = ur.MtalkNo
	returnMap["hometown"] = ur.Hometown
	returnMap["description"] = ur.Description
	returnMap["platForm"] = ur.PlatForm
	returnMap["UUID"] = ur.UUID
	returnMap["openId"] = ur.OpenId
	returnMap["backgroundImg"] = ur.BackgroundImg
	returnMap["Wechat_uid"] = ur.Wechat_uid
	returnMap["Wechat_name"] = ur.Wechat_name
	returnMap["Wechat_iconurl"] = ur.Wechat_iconurl
	returnMap["Wechat_gender"] = ur.Wechat_gender
	returnMap["QQ_uid"] = ur.QQ_uid
	returnMap["QQ_name"] = ur.QQ_name
	returnMap["QQ_iconurl"] = ur.QQ_iconurl
	returnMap["QQ_gender"] = ur.QQ_gender
	returnMap["isWater"] = ur.IsWater
	returnMap["volunteer"] = ur.Volunteer
	return returnMap
}

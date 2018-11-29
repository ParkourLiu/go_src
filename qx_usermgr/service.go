package main

import (
	"errors"
	"fmt"
	U "qx_user/caller"
)

type UserService interface {
	RegAndLogin(ur *U.UserRequest) (map[string]interface{}, string, error) //一键登录,传入手机号和操作状态，返回注册状态和userId，前端再根据此状态判断是登录还是继续其他操作
	OtherLogin(ur *U.UserRequest) (map[string]interface{}, string, error)  //第三方登录（微信，QQ）
	SearchUser(ur *U.UserRequest) (map[string]interface{}, string, error)  //查询用户信息
	UpdateUser(ur *U.UserRequest) (map[string]interface{}, string, error)  //更新用户信息
	ChangeBind(ur *U.UserRequest) (map[string]interface{}, string, error)  //绑定/解绑新电话号码，微信，QQ
}

type userService struct{}

//登录方式的Key,登录时存储，用于解绑时判断 0手机登录，1qq登录，2微信登录，3微博登录
func LoginFlag(userId string) string {
	return "LoginFlag:" + userId
}

//一键登录,传入手机号和操作状态，返回注册状态和userId，前端再根据此状态判断是登录还是继续其他操作
func (service userService) RegAndLogin(ur *U.UserRequest) (map[string]interface{}, string, error) {
	log.Debug("userService_method_start", "RegAndLogin", "input", fmt.Sprint(ur))
	returnMap := map[string]interface{}{}
	userReq, err := userClient.SearchUsers(&U.UserRequest{PhoneNo: ur.PhoneNo})
	userId := ""
	if err != nil {
		return returnMap, "", err
	}
	if userReq.Err != "" {
		return returnMap, "", errors.New(userReq.Err)
	}
	if ur.Type == "0" { //判断是否注册,未注册直接返回未注册状态
		if len(userReq.Data) < 1 { //未注册
			return returnMap, code1, nil //未注册状态
		}
	}
	if len(userReq.Data) < 1 { //未注册
		userId, err = SynRegAddUser(ur)
		if err != nil {
			return returnMap, "", err
		}
	} else {
		userId = userReq.Data[0].UserId
	}
	returnMap["userId"] = userId
	//设置登录方式为手机登录
	err = redisClient.Set(LoginFlag(userId), loginFlag_Phone)
	if err != nil {
		return returnMap, "", nil
	}
	log.Debug("userService_method_end", "RegAndLogin", "status", "success")
	return returnMap, "100", nil
}

//第三方登录（微信，QQ）
func (service userService) OtherLogin(ur *U.UserRequest) (map[string]interface{}, string, error) {
	log.Debug("userService_method_start", "OtherLogin", "input", fmt.Sprint(ur))
	returnMap := map[string]interface{}{}
	userId := ""
	//微信登录
	if "1" == ur.Type {
		request, err := userClient.SearchUsers(&U.UserRequest{Wechat_uid: ur.Wechat_uid})
		if err != nil {
			return returnMap, "", err
		}
		if request.Err != "" {
			return returnMap, "", errors.New(request.Err)
		}
		if len(request.Data) < 1 { //以前没有登录过，注册微信信息
			userId, err = SynOtherAddUser(ur) //加锁添加用户
			if err != nil {
				return returnMap, "", err
			}
		}
		if len(request.Data) == 1 { //以前登录过，更新微信信息
			request1, err1 := userClient.UpdateUser(ur)
			if err1 != nil {
				return returnMap, "", err1
			}
			if request1.Err != "" {
				return returnMap, "", errors.New(request1.Err)
			}
			userId = request.Data[0].UserId
		}
		//设置登录方式为微信登录
		err = redisClient.Set(LoginFlag(userId), loginFlag_Wechat)
		if err != nil {
			return returnMap, "", err
		}
	}

	//QQ登录
	if "2" == ur.Type {
		request, err := userClient.SearchUsers(&U.UserRequest{QQ_uid: ur.QQ_uid})
		if err != nil {
			return returnMap, "", err
		}
		if request.Err != "" {
			return returnMap, "", errors.New(request.Err)
		}
		if len(request.Data) < 1 { //以前没有登录过，注册QQ信息
			userId, err = SynOtherAddUser(ur) //加锁添加用户
			if err != nil {
				return returnMap, "", err
			}
		}
		if len(request.Data) == 1 { //以前登录过，更新QQ信息
			request1, err1 := userClient.UpdateUser(ur)
			if err1 != nil {
				return returnMap, "", err1
			}
			if request1.Err != "" {
				return returnMap, "", errors.New(request1.Err)
			}
			userId = request.Data[0].UserId
		}
		//设置登录方式为QQ登录
		err = redisClient.Set(LoginFlag(userId), loginFlag_QQ)
		if err != nil {
			return returnMap, "", err
		}
	}
	log.Debug("userService_method_end", "OtherLogin", "status", "success")
	return returnMap, "100", nil
}

//查询用户信息
func (service userService) SearchUser(ur *U.UserRequest) (map[string]interface{}, string, error) {
	log.Debug("userService_method_start", "SearchUser", "input", fmt.Sprint(ur))
	returnMap := map[string]interface{}{}
	req, err := userClient.SearchUserById(ur)
	if err != nil {
		return returnMap, "", err
	}
	if req.Err != "" {
		return returnMap, "", errors.New(req.Err)
	}
	returnMap = UserRequest2mapi(&req.Data)
	log.Debug("userService_method_end", "SearchUser", "status", "success")
	return returnMap, "100", nil
}

//修改用户信息
func (service userService) UpdateUser(ur *U.UserRequest) (map[string]interface{}, string, error) {
	log.Debug("userService_method_start", "UpdateUser", "input", fmt.Sprint(ur))
	returnMap := map[string]interface{}{}
	req, err := userClient.UpdateUser(ur)
	if err != nil {
		return returnMap, "", err
	}
	if req.Err != "" {
		return returnMap, "", errors.New(req.Err)
	}
	returnMap = UserRequest2mapi(&req.Data)
	log.Debug("userService_method_end", "UpdateUser", "status", "success")
	return returnMap, "100", nil
}

//绑定，解绑用户
func (service userService) ChangeBind(ur *U.UserRequest) (map[string]interface{}, string, error) {
	log.Debug("userService_method_start", "ChangeBind", "input", fmt.Sprint(ur))
	returnMap := map[string]interface{}{}
	if ur.Wechat_uid != "" && "0" != ur.Wechat_uid { //绑定微信前判断
		request, err := userClient.SearchUsers(&U.UserRequest{Wechat_uid: ur.Wechat_uid})
		if err != nil {
			return returnMap, "", err
		}
		if request.Err != "" {
			return returnMap, "", errors.New(request.Err)
		}
		if len(request.Data) > 0 { //已绑定
			return returnMap, code2, nil //该平台下此用户已绑定过了
		}
	} else if ur.QQ_uid != "" && "0" != ur.QQ_uid { //绑定qq前判断
		request, err := userClient.SearchUsers(&U.UserRequest{QQ_uid: ur.QQ_uid})
		if err != nil {
			return returnMap, "", err
		}
		if request.Err != "" {
			return returnMap, "", errors.New(request.Err)
		}
		if len(request.Data) > 0 { //已绑定
			return returnMap, code2, nil //该平台下此用户已绑定过了
		}
	}

	//判断是否是在操作解绑当前登录账号
	loginFlag, errl := redisClient.Get(LoginFlag(ur.UserId))
	if errl != nil {
		if errl.Error() != "redigo: nil returned" {
			return returnMap, "", errl
		}
	}
	if ur.QQ_uid == "0" {
		if loginFlag == loginFlag_QQ {
			return returnMap, code3, nil //不可解绑以当前平台登录的平台
		}
	} else if ur.Wechat_uid == "0" {
		if loginFlag == loginFlag_Wechat {
			return returnMap, code3, nil //不可解绑以当前平台登录的平台
		}
	}
	request1, err1 := userClient.UpdateUser(ur)
	if err1 != nil {
		return returnMap, "", err1
	}
	if request1.Err != "" {
		return returnMap, "", errors.New(request1.Err)
	}
	log.Debug("userService_method_end", "ChangeBind", "status", "success")
	return returnMap, "100", nil
}

package main

import (
	"errors"
	"fmt"
	search "search/client"
	"strings"
	"time"
	UserClient "user/client"
)

type UserService interface {
	Reg(*User) (string, string, string, error)                           //注册
	Login(*User) (string, string, string, error)                         //登录
	ShortcutLogin(*User) (string, string, string, error)                 //快捷登录
	FindPassword(*User) (string, string, error)                          //找回密码
	ChangePhoneNo(*User) (string, error)                                 //更换手机号
	OtherLogin(*User) (string, string, string, string, error)            //第三方登录
	UpdateUser(*User) (string, error)                                    //修改用户
	MyHome(*User) (map[string]string, []string, []string, string, error) //个人主页信息
	HomePageSlogan() ([]map[string]string, string, error)                //首页广告语
	ActiveUser(user *User) (string, error)                               //统计活跃用户

	//4.1.0
	CheckPhoneBook(*PhoneBook) (string, error, string)                                  //Check通讯录是否变化,返回 code； err； 判断结果0无，1有
	PhoneBookUser(*PhoneBook) (string, error, []map[string]string, []map[string]string) //通讯录用户
}

type userService struct{}

type User struct {
	LookUserId      string
	UserId          string
	PhoneNo         string
	Password        string
	Email           string
	TrueName        string
	NickName        string
	BirthDay        string
	ChineseZodiac   string
	Sex             string
	HomeAddress     string
	ImageName       string
	ChatName        string
	ChatPwd         string
	MtalkNo         string
	Hometown        string
	Description     string
	PlatForm        string
	OpenId          string
	BackgroundImg   string
	Sina_uid        string
	Sina_name       string
	Sina_iconurl    string
	Sina_gender     string
	Wechat_uid      string
	Wechat_name     string
	Wechat_iconurl  string
	Wechat_gender   string
	QQ_uid          string
	QQ_name         string
	QQ_iconurl      string
	QQ_gender       string
	Ad              string
	OtherLogienType string //0微博,1微信，2qq
	Id              string //数据导入临时使用，无实际业务意义
}

type PhoneBook struct {
	UserId   string              `json:"userId"`   //用户id
	HashCode string              `json:"hashCode"` //用户通讯录list的hash值
	Phones   []string            `json:"phones"`   //用户通讯录电话
	Data     []map[string]string `json:"data"`     //用户通讯录电话
}

func (user *User) String() string {
	return "UserId: " + user.UserId + ", PhoneNo: " + user.PhoneNo + ", Password: " + user.Password + ", Email: " + user.Email + ", TrueName: " + user.TrueName + ", BirthDay: " + user.BirthDay + ", ChineseZodiac: " + user.ChineseZodiac + ", Sex: " + user.Sex + ", HomeAddress: " + user.HomeAddress + ", ImageName: " + user.ImageName + ", ChatName: " + user.ChatName + ", ChatPwd: " + user.ChatPwd + ", MtalkNo: " + user.MtalkNo + ", Hometown: " + user.Hometown + ", Description: " + user.Description + ", PlatForm: " + user.PlatForm + ", OpenId: " + user.OpenId + ", BackgroundImg: " + user.BackgroundImg + "Sina_uid: " + user.Sina_uid + ", Sina_name: " + user.Sina_name + ", Sina_iconurl: " + user.Sina_iconurl + ", Sina_gender: " + user.Sina_gender + "Wechat_uid: " + user.Wechat_uid + ", Wechat_name: " + user.Wechat_name + ", Wechat_iconurl: " + user.Wechat_iconurl + ", Wechat_gender: " + user.Wechat_gender + "QQ_uid: " + user.QQ_uid + ", QQ_name: " + user.QQ_name + ", QQ_iconurl: " + user.QQ_iconurl + ", QQ_gender: " + user.QQ_gender + ", Ad: " + user.Ad + ", OtherLogienType: " + user.OtherLogienType
}

//登录方式的Key,登录时存储，用于解绑时判断 0手机登录，1qq登录，2微信登录，3微博登录
func LoginFlag(userId string) string {
	return "LoginFlag:" + userId
}

//个人主页
func (service userService) MyHome(user *User) (map[string]string, []string, []string, string, error) {
	log.Debug("method_start", "MyHome", "input", user)
	u := &UserClient.UserRequest{UserId: user.UserId}

	log.Debug("█████myHome开始调用SearchUser时间：", fmt.Sprint(time.Now()))
	req, err := userClient.SearchUserMap(u)
	log.Debug("█████myHome结束调用SearchUser时间：", fmt.Sprint(time.Now()))
	if err != nil {
		log.Debug("████1", err.Error())
		return nil, nil, nil, "", err
	}
	userMap := req.User
	if userMap == nil {
		return map[string]string{}, nil, nil, "100", nil
	}
	//查看此人粉丝数
	fansCount, err4 := redisClient.Zlen("grpfans:uIdFansList:" + user.UserId)
	if err4 != nil {
		if "redigo: nil returned" != err4.Error() {
			log.Debug("████2", err4.Error())
			return nil, nil, nil, "", err4
		}
	}

	//查看我是否关注此人
	score, err5 := redisClient.Zscore("grpfans:uIdFansList:"+user.UserId, user.LookUserId)
	if err5 != nil {
		if "redigo: nil returned" != err5.Error() {
			log.Debug("████3", err5.Error())
			return nil, nil, nil, "", err5
		}

	}
	likeFlag := "0"
	if score != 0 {
		likeFlag = "1"
	}

	//查看此人参与了多少活动
	activityCount, err6 := redisClient.Zlen("mactivity:UIdAttendActivityList:" + user.UserId)
	if err6 != nil {
		if "redigo: nil returned" != err6.Error() {
			log.Debug("████activityCount", fmt.Sprint(activityCount), "████4", err6.Error())
			return nil, nil, nil, "", err6
		}

	}

	likeCount := int64(0)
	//查看此人关注了多少人
	count, err7 := redisClient.Zlen("grpfans:attentionUIdList:" + user.UserId)
	if err7 != nil {
		if "redigo: nil returned" != err7.Error() {
			log.Debug("████5", err7.Error())
			return nil, nil, nil, "", err7
		}
	}
	likeCount += count
	//查看此人关注了多少社群
	count, err7 = redisClient.Zlen("grpfans:attentionGroupList:" + user.UserId)
	if err7 != nil {
		if "redigo: nil returned" != err7.Error() {
			log.Debug("████5", err7.Error())
			return nil, nil, nil, "", err7
		}
	}
	likeCount += count
	//查看此人关注了多少活动
	count, err7 = redisClient.Zlen("grpfans:attentionActivityList:" + user.UserId)
	if err7 != nil {
		if "redigo: nil returned" != err7.Error() {
			log.Debug("████5", err7.Error())
			return nil, nil, nil, "", err7
		}
	}
	likeCount += count

	//查看此人是否发布求助
	helpList, err8 := redisClient.LrangeAll("UC2C:" + user.UserId)
	if err8 != nil {
		if "redigo: nil returned" != err8.Error() {
			log.Debug("████8", err8.Error())
			return nil, nil, nil, "", err8
		}

	}
	helpList_A := []string{}
	for _, v := range helpList {
		status, _ := redisClient.Hget("CALL:"+v, "status")
		if status == "0" {
			helpList_A = append(helpList_A, v)
		}
	}

	helpPetList, err9 := redisClient.LrangeAll("UC2Cp:" + user.UserId)
	if err9 != nil {
		if "redigo: nil returned" != err9.Error() {
			log.Debug("████9", err9.Error())
			return nil, nil, nil, "", err9
		}

	}
	helpPetList_A := []string{}
	for _, v := range helpPetList {
		status, _ := redisClient.Hget("CALL:"+v, "status")
		if status == "0" {
			helpPetList_A = append(helpPetList_A, v)
		}
	}
	userMap["haveHelp"] = "0"
	if helpList_A != nil && len(helpList_A) > 0 {
		userMap["haveHelp"] = "1"
	}
	if helpPetList_A != nil && len(helpPetList_A) > 0 {
		userMap["haveHelp"] = "1"
	}
	userMap["fansCount"] = fmt.Sprint(fansCount)
	userMap["activityCount"] = fmt.Sprint(activityCount)
	userMap["likeFlag"] = likeFlag
	userMap["likeCount"] = fmt.Sprint(likeCount)
	userMap["havePassword"] = "0"
	if userMap["password"] != "" {
		userMap["havePassword"] = "1"
	}
	delete(userMap, "password")
	log.Debug("method_end", "MyHome", "status", "success")
	return userMap, helpList_A, helpPetList_A, "100", nil
}

//手机号注册
func (service userService) Reg(user *User) (string, string, string, error) {

	log.Debug("method_start", "Reg", "input", user)
	havePassword := "0"
	if user.Password != "" {
		havePassword = "1"
	}
	userId := ""
	u := &UserClient.UserRequest{PhoneNo: user.PhoneNo}
	request, err := userClient.SearchUsers(u)
	if err != nil {
		return "", "", "", err
	}
	if request.Err != "" {
		return "", "", "", errors.New(request.Err)
	}
	if len(request.Users) > 0 {
		return flag5, "", "", nil //该手机号已注册
	} else { //开始注册
		addUserId, code, err1 := SynRegAddUser(user)
		if err1 != nil {
			return "", "", "", err1
		}
		if code == flag5 {
			return flag5, "", "", nil //该手机号已注册
		}
		userId = addUserId
	}

	//设置登录方式为手机登录
	err = redisClient.Set(LoginFlag(userId), loginFlag_Phone)
	if err != nil {
		return "", "", "", err
	}

	//信息添加到solr,添加紧急联系人
	user.UserId = userId
	go upSolr(user)
	go CTR(user)

	log.Debug("method_end", "Reg", "status", "success")
	return "100", userId, havePassword, nil
}

//手机号注册加锁添加用户
func SynRegAddUser(user *User) (string, string, error) { //添加方法加锁,返回userId,code,err
	ex, err2 := redisClient.SETNX("muser:P:"+user.PhoneNo, "1")
	defer redisClient.SetStringAndExpire("muser:P:"+user.PhoneNo, "1", 1)
	if err2 != nil {
		return "", "", err2
	}
	if ex != 1 {
		return "", "", errors.New("请不要重复提交")
	}
	userId := ""
	u := &UserClient.UserRequest{PhoneNo: user.PhoneNo}
	request, err := userClient.SearchUsers(u) //查询数据库数据
	if err != nil {
		return "", "", err
	}
	if request.Err != "" {
		return "", "", errors.New(request.Err)
	}
	if len(request.Users) > 0 {
		return "", flag5, nil
	} else {
		u.Password = user.Password
		rep, err1 := userClient.AddUser(u)
		if err1 != nil {
			return "", "", err1
		}
		if rep.Err != "" {
			return "", "", errors.New(rep.Err)
		}
		userId = rep.UserId
	}
	if userId == "" {
		return "", "", errors.New("请不要重复提交")
	}
	return userId, "", nil
}

func CTR(user *User) error {
	ctId := "CT" + idGenClient.GetUniqueId()
	ctMap := map[string]string{}
	ctMap["phoneNo"] = user.PhoneNo
	ctMap["trueName"] = "钦家用户"
	ctMap["isDef"] = "1" //设置成默认联系人
	err := redisClient.Hmset("CTR:"+ctId, ctMap)
	if err != nil {
		return err
	}
	return redisClient.LpushFromHead("U2CTR:"+user.UserId, ctId)
}

//手机号，密码登陆
func (service userService) Login(user *User) (string, string, string, error) {
	userId := ""
	havePassword := "0"
	log.Debug("method_start", "Login", "input", user)
	if "" == user.PhoneNo || "" == user.Password {
		return flag2, "", "", nil //账号或密码不可以为空
	}
	log.Debug("███████0:", "")
	u := &UserClient.UserRequest{PhoneNo: user.PhoneNo, Password: user.Password}
	log.Debug("███████1:", fmt.Sprint(u))
	request, err := userClient.SearchUsers(u)
	if err != nil {
		return "", "", "", err
	}
	if request.Err != "" {
		return "", "", "", errors.New(request.Err)
	}
	if len(request.Users) < 1 {
		return flag1, "", "", nil //账号或密码错误
	}
	if len(request.Users) > 1 {
		return flag3, "", "", nil //数据库相同手机密码有多个
	}
	if request.Users[0].PhoneNo != user.PhoneNo {
		return flag1, "", "", nil //账号或密码错误
	}
	if request.Users[0].Password != user.Password {
		return flag1, "", "", nil //账号或密码错误
	}
	userId = request.Users[0].UserId
	if request.Users[0].Password != "" {
		havePassword = "1"
	}
	log.Debug("███████userId:", fmt.Sprint(userId))
	//设置登录方式为手机登录
	err = redisClient.Set(LoginFlag(userId), loginFlag_Phone)

	log.Debug("███████err:", fmt.Sprint(err))
	if err != nil {
		return "", "", "", err
	}
	log.Debug("method_end", "Login", "status", "success")
	return "100", userId, havePassword, nil
}

//手机号快捷登录
func (service userService) ShortcutLogin(user *User) (string, string, string, error) {
	log.Debug("method_start", "ShortcutLogin", "input", user)
	userId := ""
	havePassword := "0"
	u := &UserClient.UserRequest{PhoneNo: user.PhoneNo}
	request, err := userClient.SearchUsers(u)
	log.Debug("███手机号快捷登录request:", fmt.Sprint(request), "长度", fmt.Sprint(request.Users))
	if err != nil {
		return "", "", "", err
	}
	if request.Err != "" {
		return "", "", "", errors.New(request.Err)
	}
	if len(request.Users) < 1 { //该手机号以前未注册
		code := ""
		code, userId, _, err = service.Reg(&User{PhoneNo: user.PhoneNo})
		if err != nil {
			return code, "", "", err
		}
	} else if len(request.Users) == 1 {
		userId = request.Users[0].UserId
		if request.Users[0].Password != "" {
			havePassword = "1"
		}
	} else {
		return flag3, "", "", nil
	}
	//设置登录方式为手机登录
	err = redisClient.Set(LoginFlag(userId), loginFlag_Phone)
	if err != nil {
		return "", "", "", err
	}
	log.Debug("method_end", "ShortcutLogin", "status", "success")
	return "100", userId, havePassword, nil
}

//找回密码
func (service userService) FindPassword(user *User) (string, string, error) {
	log.Debug("method_start", "FindPassword", "input", user)
	userId := ""
	u := &UserClient.UserRequest{PhoneNo: user.PhoneNo}
	request, err := userClient.SearchUsers(u)
	if err != nil {
		return "", "", err
	}
	if request.Err != "" {
		return "", "", errors.New(request.Err)
	}
	u.UserId = request.Users[0].UserId
	u.Password = user.Password
	request1, err1 := userClient.UpdateUser(u)
	if err1 != nil {
		return "", "", err1
	}
	if request1.Err != "" {
		return "", "", errors.New(request1.Err)
	}
	userId = request.Users[0].UserId
	log.Debug("method_end", "FindPassword", "status", "success")
	return "100", userId, nil
}

//更换手机号
func (service userService) ChangePhoneNo(user *User) (string, error) {
	log.Debug("method_start", "ChangePhoneNo", "input", user)

	u := &UserClient.UserRequest{UserId: user.UserId, PhoneNo: user.PhoneNo}
	request, err := userClient.UpdateUser(u)
	if err != nil {
		return "", err
	}
	if request.Err != "" {
		return "", errors.New(request.Err)
	}
	fmt.Println("==========================", request)
	log.Debug("method_end", "ChangePhoneNo", "status", "success")
	return "100", nil
}

//第三方登录
func (service userService) OtherLogin(user *User) (string, string, string, string, error) { //code userId, phoneNo, havePassword, nil
	log.Debug("method_start", "OtherLogin", "input", user)
	havePassword := "0"
	phoneNo := ""
	userId := ""
	if "0" != user.OtherLogienType && "1" != user.OtherLogienType && "2" != user.OtherLogienType {
		return flag4, "", "", "", nil
	}

	//微博登录
	if "0" == user.OtherLogienType {
		u := &UserClient.UserRequest{Sina_uid: user.Sina_uid}
		request, err := userClient.SearchUsers(u)
		if err != nil {
			return "", "", "", "", err
		}
		if request.Err != "" {
			return "", "", "", "", errors.New(request.Err)
		}
		if len(request.Users) < 1 { //以前没有登录过，注册微博信息
			userId, err = SynOtherAddUser(u, user, "0") //加锁添加用户
			if err != nil {
				return "", "", "", "", err
			}
		}
		if len(request.Users) == 1 { //以前登录过，更新微博信息
			if request.Users[0].Password != "" {
				havePassword = "1"
			}
			phoneNo = request.Users[0].PhoneNo
			u.UserId = request.Users[0].UserId
			u.Sina_name = user.Sina_name
			u.Sina_iconurl = user.Sina_iconurl
			u.Sina_gender = user.Sina_gender
			request1, err1 := userClient.UpdateUser(u)
			if err1 != nil {
				return "", "", "", "", err1
			}
			if request1.Err != "" {
				return "", "", "", "", errors.New(request1.Err)
			}
			userId = request.Users[0].UserId
		}
		if len(request.Users) > 1 { //数据库有多条ad为a的Sina_uid信息
			return flag3, "", "", "", nil
		}

		//设置登录方式为微博登录
		err = redisClient.Set(LoginFlag(userId), loginFlag_Sina)
		if err != nil {
			return "", "", "", "", err
		}
	}

	//微信登录
	if "1" == user.OtherLogienType {
		u := &UserClient.UserRequest{Wechat_uid: user.Wechat_uid}
		request, err := userClient.SearchUsers(u)
		if err != nil {
			return "", "", "", "", err
		}
		if request.Err != "" {
			return "", "", "", "", errors.New(request.Err)
		}
		if len(request.Users) < 1 { //以前没有登录过，注册微信信息
			userId, err = SynOtherAddUser(u, user, "1") //加锁添加用户
			if err != nil {
				return "", "", "", "", err
			}
		}
		if len(request.Users) == 1 { //以前登录过，更新微信信息
			if request.Users[0].Password != "" {
				havePassword = "1"
			}
			phoneNo = request.Users[0].PhoneNo
			u.UserId = request.Users[0].UserId
			u.Wechat_name = user.Wechat_name
			u.Wechat_iconurl = user.Wechat_iconurl
			u.Wechat_gender = user.Wechat_gender
			request1, err1 := userClient.UpdateUser(u)
			if err1 != nil {
				return "", "", "", "", err1
			}
			if request1.Err != "" {
				return "", "", "", "", errors.New(request1.Err)
			}
			userId = request.Users[0].UserId
		}
		if len(request.Users) > 1 { //数据库有多条ad为a的Wechat_uid信息
			return flag3, "", "", "", nil
		}
		//设置登录方式为微信登录
		err = redisClient.Set(LoginFlag(userId), loginFlag_Wechat)
		if err != nil {
			return "", "", "", "", err
		}
	}

	//QQ登录
	if "2" == user.OtherLogienType {
		u := &UserClient.UserRequest{QQ_uid: user.QQ_uid}
		request, err := userClient.SearchUsers(u)
		if err != nil {
			return "", "", "", "", err
		}
		if request.Err != "" {
			return "", "", "", "", errors.New(request.Err)
		}
		if len(request.Users) < 1 { //以前没有登录过，注册QQ信息
			userId, err = SynOtherAddUser(u, user, "2") //加锁添加用户
			if err != nil {
				return "", "", "", "", err
			}
		}
		if len(request.Users) == 1 { //以前登录过，更新QQ信息
			if request.Users[0].Password != "" {
				havePassword = "1"
			}
			phoneNo = request.Users[0].PhoneNo
			u.UserId = request.Users[0].UserId
			u.QQ_name = user.QQ_name
			u.QQ_iconurl = user.QQ_iconurl
			u.QQ_gender = user.QQ_gender
			request1, err1 := userClient.UpdateUser(u)
			if err1 != nil {
				return "", "", "", "", err1
			}
			if request1.Err != "" {
				return "", "", "", "", errors.New(request1.Err)
			}
			userId = request.Users[0].UserId
		}
		if len(request.Users) > 1 { //数据库有多条ad为a的QQ_uid信息
			return flag3, "", "", "", nil
		}
		//设置登录方式为QQ登录
		err = redisClient.Set(LoginFlag(userId), loginFlag_QQ)
		if err != nil {
			return "", "", "", "", err
		}
	}

	log.Debug("method_end", "OtherLogin", "status", "success")
	return "100", userId, phoneNo, havePassword, nil
}

func SynOtherAddUser(u *UserClient.UserRequest, user *User, from string) (string, error) { //添加方法加锁
	if from == "0" { //新浪
		ex, err2 := redisClient.SETNX("muser:S:"+user.Sina_uid, "1")
		defer redisClient.SetStringAndExpire("muser:S:"+user.Sina_uid, "1", 1)
		if err2 != nil {
			return "", err2
		}
		if ex != 1 {
			return "", errors.New("请不要重复提交")
		}
		u.Sina_name = user.Sina_name
		u.Sina_iconurl = user.Sina_iconurl
		u.Sina_gender = user.Sina_gender
	} else if from == "1" { //微信
		ex, err2 := redisClient.SETNX("muser:W:"+user.Wechat_uid, "1")
		defer redisClient.SetStringAndExpire("muser:W:"+user.Wechat_uid, "1", 1)
		if err2 != nil {
			return "", err2
		}
		if ex != 1 {
			return "", errors.New("请不要重复提交")
		}
		u.Wechat_name = user.Wechat_name
		u.Wechat_iconurl = user.Wechat_iconurl
		u.Wechat_gender = user.Wechat_gender
	} else if from == "2" { //QQ
		ex, err2 := redisClient.SETNX("muser:Q:"+user.QQ_uid, "1")
		defer redisClient.SetStringAndExpire("muser:Q:"+user.QQ_uid, "1", 1)
		if err2 != nil {
			return "", err2
		}
		if ex != 1 {
			return "", errors.New("请不要重复提交")
		}
		u.QQ_name = user.QQ_name
		u.QQ_iconurl = user.QQ_iconurl
		u.QQ_gender = user.QQ_gender
	}

	userId := ""
	request, err := userClient.SearchUsers(u) //查询数据库数据
	if err != nil {
		return "", err
	}
	if request.Err != "" {
		return "", errors.New(request.Err)
	}
	if len(request.Users) == 0 {
		rep, err1 := userClient.AddUser(u)
		if err1 != nil {
			return "", err1
		}
		if rep.Err != "" {
			return "", errors.New(rep.Err)
		}
		userId = rep.UserId
	}
	if userId == "" {
		return "", errors.New("请不要重复提交")
	}

	return userId, nil
}

func (service userService) UpdateUser(user *User) (string, error) {
	log.Debug("method_start", "UpdateUser", "input", user)
	if user.Sina_uid != "" && "0" != user.Sina_uid { //绑定微博前判断
		u := &UserClient.UserRequest{Sina_uid: user.Sina_uid}
		request, err := userClient.SearchUsers(u)
		if err != nil {
			return "", err
		}
		if request.Err != "" {
			return "", errors.New(request.Err)
		}
		if len(request.Users) > 0 { //已绑定
			return flag6, nil //该平台下此用户已绑定过了
		}
	} else if user.Wechat_uid != "" && "0" != user.Wechat_uid { //绑定微信前判断
		u := &UserClient.UserRequest{Wechat_uid: user.Wechat_uid}
		request, err := userClient.SearchUsers(u)
		if err != nil {
			return "", err
		}
		if request.Err != "" {
			return "", errors.New(request.Err)
		}
		if len(request.Users) > 0 { //已绑定
			return flag6, nil //该平台下此用户已绑定过了
		}
	} else if user.QQ_uid != "" && "0" != user.QQ_uid { //绑定qq前判断
		u := &UserClient.UserRequest{QQ_uid: user.QQ_uid}
		request, err := userClient.SearchUsers(u)
		if err != nil {
			return "", err
		}
		if request.Err != "" {
			return "", errors.New(request.Err)
		}
		if len(request.Users) > 0 { //已绑定
			return flag6, nil //该平台下此用户已绑定过了
		}
	}

	//判断是否是在操作解绑当前登录账号
	loginFlag, errl := redisClient.Get(LoginFlag(user.UserId))
	if errl != nil {
		if errl.Error() != "redigo: nil returned" {
			return "", errl
		}
	}
	if user.QQ_uid == "0" {
		if loginFlag == loginFlag_QQ {
			return flag7, nil
		}
	} else if user.Wechat_uid == "0" {
		if loginFlag == loginFlag_Wechat {
			return flag7, nil
		}
	} else if user.Sina_uid == "0" {
		if loginFlag == loginFlag_Sina {
			return flag7, nil
		}
	}

	//用户名拦截
	if user.NickName != "" {
		nameFlag := qinJiaName(user.NickName)
		if nameFlag == false {
			return flag8, nil
		}
	}
	if user.TrueName != "" {
		nameFlag := qinJiaName(user.TrueName)
		if nameFlag == false {
			return flag8, nil
		}
	}

	u := &UserClient.UserRequest{
		PhoneNo:        user.PhoneNo,
		Password:       user.Password,
		Sina_uid:       user.Sina_uid,
		Sina_name:      user.Sina_name,
		Sina_iconurl:   user.Sina_iconurl,
		Sina_gender:    user.Sina_gender,
		Wechat_uid:     user.Wechat_uid,
		Wechat_name:    user.Wechat_name,
		Wechat_iconurl: user.Wechat_iconurl,
		Wechat_gender:  user.Wechat_gender,
		QQ_uid:         user.QQ_uid,
		QQ_name:        user.QQ_name,
		QQ_iconurl:     user.QQ_iconurl,
		QQ_gender:      user.QQ_gender,
		UserId:         user.UserId,
		TrueName:       user.TrueName,
		NickName:       user.NickName,
		BirthDay:       user.BirthDay,
		ChineseZodiac:  user.ChineseZodiac,
		Sex:            user.Sex,
		HomeAddress:    user.HomeAddress,
		ImageName:      user.ImageName,
		ChatName:       user.ChatName,
		ChatPwd:        user.ChatPwd,
		Hometown:       user.Hometown,
		Description:    user.Description,
		PlatForm:       user.PlatForm,
		OpenId:         user.OpenId,
		BackgroundImg:  user.BackgroundImg,
		MtalkNo:        user.MtalkNo,
		Id:             user.Id,
	}
	request1, err1 := userClient.UpdateUser(u)
	if err1 != nil {
		return "", err1
	}
	if request1.Err != "" {
		return "", errors.New(request1.Err)
	}

	if user.NickName != "" {

		//将数据导入solr
		list := []search.SearchInfo{}
		p := map[string]interface{}{"id": user.UserId, "uname": user.NickName, "ad": "a"}
		s := search.SearchInfo{Param: p, UpdateType: "0"}
		list = append(list, s)
		i := &search.Infos{TableName: "mtalk4_user", SearchInfos: list}
		err := searchClient.Update(i)

		//更新solr新数据前会覆盖其他字段，所以先查询出其他字段出来一起更新
		user.PhoneNo, err1 = redisClient.Hget("U:"+user.UserId, "phoneNo")
		if err1 != nil {
			return "", err1
		}
		//信息导入solr
		err = upSolr(user)

		if err != nil {
			return "", err
		}
	}

	log.Debug("method_end", "UpdateUser", "status", "success")
	return "100", nil
}

//将用户数据导入solr
func upSolr(user *User) error {
	//将数据导入solr
	list := []search.SearchInfo{}
	p := map[string]interface{}{"id": user.UserId, "phone_no": user.PhoneNo, "uname": user.NickName, "ad": "a"}
	s := search.SearchInfo{Param: p, UpdateType: "0"}
	list = append(list, s)
	i := &search.Infos{TableName: "mtalk4_user", SearchInfos: list}
	return searchClient.Update(i)
}

//首页广告语
func (service userService) HomePageSlogan() ([]map[string]string, string, error) {
	log.Debug("method_start", "homePageSlogan", "input", "")
	returnList := []map[string]string{}
	helpSlogan, err := redisClient.HgetAllMap("helpSlogan") //求助
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return nil, "", err
		}
	}
	publicBenefitSlogan, err1 := redisClient.HgetAllMap("publicBenefitSlogan") //公益
	if err1 != nil {
		if err1.Error() != "redigo: nil returned" {
			return nil, "", err1
		}
	}
	returnList = append(returnList, helpSlogan)
	returnList = append(returnList, publicBenefitSlogan)
	log.Debug("method_end", "homePageSlogan", "status", "success")
	return returnList, "100", nil
}

func PhoneBookHashKey(userId string) string {
	//string类型
	return "PhoneBookHash:" + userId
}

//check用户通讯录
func (service userService) CheckPhoneBook(phoneBook *PhoneBook) (string, error, string) {
	log.Debug("method_start", "CheckPhoneBook", "input", fmt.Sprint(phoneBook))
	hashCode, err := redisClient.Get(PhoneBookHashKey(phoneBook.UserId))
	if err != nil {
		return "", err, ""
	}
	hashCodeFlag := "0"
	if hashCode == phoneBook.HashCode {
		hashCodeFlag = "1"
		return "100", nil, hashCodeFlag
	}
	err = redisClient.Set(PhoneBookHashKey(phoneBook.UserId), phoneBook.HashCode)
	if err != nil {
		return "", err, ""
	}
	log.Debug("method_end", "CheckPhoneBook", "status", "success")
	return "100", nil, hashCodeFlag
}

func PhoneBookKey(userId string) string {
	//list
	return "PhoneBook:" + userId
}

//通讯录用户
func (service userService) PhoneBookUser(phoneBook *PhoneBook) (string, error, []map[string]string, []map[string]string) {
	log.Debug("method_start", "PhoneBookUser", "input", fmt.Sprint(phoneBook))
	if len(phoneBook.Phones) < 1 { //取redis数据
		err, haveUser, noUser := RedisPhonesUser(phoneBook.UserId)
		if err != nil {
			return "", err, nil, nil
		}
		return "100", nil, haveUser, noUser
	}
	//取传进来的数据
	err := redisClient.Del(PhoneBookKey(phoneBook.UserId))
	if err != nil {
		return "", err, nil, nil
	}
	c, err1 := redisClient.GetPipeline()
	if err1 != nil {
		return "", err1, nil, nil
	}
	for _, v := range phoneBook.Phones {
		redisClient.PipelineLpushFromTail(c, PhoneBookKey(phoneBook.UserId), v)
	}
	err = redisClient.ExecutePipeline(c)
	if err != nil {
		return "", err, nil, nil
	}
	err, haveUser, noUser := RedisPhonesUser(phoneBook.UserId)
	if err != nil {
		return "", err, nil, nil
	}
	log.Debug("method_end", "PhoneBookUser", "status", "success")
	return "100", nil, haveUser, noUser
}

//统计活跃用户
func (service userService) ActiveUser(user *User) (string, error) {
	log.Debug("method_start", "ActiveUser", "input", user)
	if user.UserId != "" {
		score, err := redisClient.Zscore("A", user.UserId)
		if err != nil {
			if err.Error() != "redigo: nil returned" {
				return "", err
			}
		}
		if score == 0 {
			_, err := redisClient.Zadd("A", int64(time.Now().Unix()), user.UserId)
			if err != nil {
				return "", err
			}
		}
	}

	log.Debug("method_end", "ActiveUser", "status", "success")
	return "100", nil
}

func RedisPhonesUser(userId string) (error, []map[string]string, []map[string]string) {
	haveUser := []map[string]string{}
	noUser := []map[string]string{}
	phones, err := redisClient.LrangeAll(PhoneBookKey(userId))
	if err != nil {
		return err, nil, nil
	}
	if len(phones) > 0 {
		u := &UserClient.UserRequest{Phones: phones}
		request, err := userClient.SearchUsersByPhones(u)
		if err != nil {
			return err, nil, nil
		}
		if request.Err != "" {
			return errors.New(request.Err), nil, nil
		}
		haveUser = request.Users
	} else {
		return nil, haveUser, noUser
	}
	//挑出已未注册用户
	for _, phone := range phones {
		userFlag := 1
		for _, userMap := range haveUser {
			if phone == userMap["phoneNo"] {
				userFlag = 0
				break
			}
		}
		if userFlag == 1 {
			noUserMap := map[string]string{}
			noUserMap["phoneNo"] = phone
			noUserMap["imageName"] = ""
			noUserMap["userId"] = ""
			noUserMap["nickName"] = ""
			noUserMap["status"] = "2"
			noUserMap["localName"] = ""
			noUser = append(noUser, noUserMap)
		}
	}
	//查看已注册用户是否关注
	for i, userMap := range haveUser {
		//查看我是否关注此人
		score, err5 := redisClient.Zscore("grpfans:uIdFansList:"+userMap["userId"], userId)
		if err5 != nil {
			return err5, nil, nil
		}
		if score != 0 {
			haveUser[i]["status"] = "1"
		} else {
			haveUser[i]["status"] = "0"
		}
		haveUser[i]["localName"] = ""
	}
	return nil, haveUser, noUser
}

//用户名拦截,
func qinJiaName(nickName string) bool {
	names := strings.Split(QinJiaName, ",")
	for _, v := range names {
		if v == nickName {
			return false
		}
	}
	return true
}

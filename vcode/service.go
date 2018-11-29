package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	logger "mtcomm/log"
	"regexp"
	sms "sms/client"
	"strings"
	"time"
	UserClient "user/client"
)

const (
	regular  = "^(1[0-9])\\d{9}$"
	reqular2 = "^\\d{4}$"
)

type CodeService interface {
	CheckCode(Code) (string, error)

	GetCode(Code) (string, error)
}

type codeService struct {
}

type Code struct {
	UserId  string
	PhoneNo string
	Type    string
	Vcode   string
}

func (code *Code) String() string {
	return "UserId: " + code.UserId + "PhoneNo: " + code.PhoneNo + ", Type: " + code.Type + ",Vocode:" + code.Vcode
}

func (service codeService) GetCode(code Code) (string, error) {
	log := logger.GetDefaultLogger()
	log.Debug("method_start", "GetCode", "check", code)
	log.Debug("███0█code", fmt.Sprint(code))

	/* check */
	//	err := getCodeCheck(log, code)
	//	if err != nil {
	//		return "", err
	//	}

	if "0" != code.Type && "1" != code.Type && "2" != code.Type && "3" != code.Type && "4" != code.Type && "5" != code.Type && "6" != code.Type && "7" != code.Type {
		return flag1, nil
	}
	/* service */ //验证类型  0:注册；1：更改密码；2：修改绑定手机号；3：绑定紧急联系人；4：停用定位贴；5：快捷登录；6找回密码;,7创建学校验证（用于直接验证手机号，不需要验证用户信息）
	//先判断是否是注册,如果是注册,必须先查询该手机号码是是否注册过
	if "0" == code.Type || "6" == code.Type { //0:注册，6找回密码
		//判断是否注册过
		u := &UserClient.UserRequest{PhoneNo: code.PhoneNo}
		users, err1 := userClient.SearchUsers(u)
		log.Debug("███0█users", users)
		if err1 != nil {
			return "", err1
		}
		log.Debug("type_0_users:", users)
		if "0" == code.Type && len(users.Users) > 0 {
			return flag2, nil //该用户已注册
		} else if "6" == code.Type && len(users.Users) == 0 {
			return flag3, nil //该用户不存在
		}

	} else if "1" == code.Type || "2" == code.Type || "3" == code.Type || "4" == code.Type { //1：更改密码，2修改绑定手机号
		redisUser, err0 := redisClient.HgetAllMap("U:" + code.UserId)
		if err0 != nil {
			if err0.Error() != "redigo: nil returned" {
				return "", err0
			}
		}
		if len(redisUser) < 1 {
			return flag3, nil
		}
		if redisUser["phoneNo"] != code.PhoneNo {
			log.Debug(redisUser, "██redisUser['phoneNo']"+redisUser["phoneNo"])
			return flag4, nil
		}
	}

	//判断用户是否在57s内申请过验证码
	log.Debug("███0█Redis前", "进入Redis")
	b, _ := redisClient.Exist(getIsCanSendMsgAgainKey(code))
	log.Debug("███0█Redis后", "进入Redis")
	if b {
		return flag5, nil
	}

	//将code存入redis 并调用短信服务发送短信服务

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%04v", rnd.Int31n(10000))
	log.Debug("███1█vcode", vcode)
	param := []string{vcode}
	phoneNos := []string{code.PhoneNo}
	req := &sms.SmsRequest{
		Params: param,
		Mobile: phoneNos,
		Tpl_id: "57981",
	}
	log.Debug("███2█vcode", vcode)
	err := smsClient.Sms(req)
	if err != nil {
		return "", err
	}
	c, _ := redisClient.GetPipeline()
	log.Debug("███3█vcode", vcode)
	redisClient.PipelineSetStringAndExpire(c, getCodePrimaryKey(code), vcode, uint32(60*30))    //设置缓存30分钟
	redisClient.PipelineSetStringAndExpire(c, getIsCanSendMsgAgainKey(code), vcode, uint32(57)) //设置57s内用户不能再次申请验证码
	err1 := redisClient.ExecutePipeline(c)
	if err1 != nil {
		log.Debug("████err1", err1.Error())
		return "", err1
	}
	log.Debug("███4█vcode", vcode)

	log.Debug("method_end", "GetCode", "status", "success")
	return "100", nil
}

func (service codeService) CheckCode(code Code) (string, error) {
	log := logger.GetDefaultLogger()
	log.Debug("method_start", "CheckCode", "check", code)

	/* check */
	//	err := checkCodeCheck(log, code)
	//	if err != nil {
	//		return "", err
	//	}
	//判断验证码是否正确
	v, _ := redisClient.Get(getCodePrimaryKey(code))
	log.Debug("●●●●●v", v)
	if v != code.Vcode {
		return flag6, nil
	}
	log.Debug("method_end", "CheckCode", "status", "success")
	return "100", nil
}

func getCodePrimaryKey(code Code) string { //拼接得到验证码的主key
	var buffer bytes.Buffer
	buffer.WriteString("vcode:")
	buffer.WriteString(code.PhoneNo)
	buffer.WriteString(":")
	buffer.WriteString(code.Type)
	return buffer.String()
}

//判断是发可以再次获取验证码的key
func getIsCanSendMsgAgainKey(code Code) string {
	var buffer bytes.Buffer
	buffer.WriteString("vcode:")
	buffer.WriteString(code.PhoneNo)
	buffer.WriteString(":")
	buffer.WriteString(code.Type)
	buffer.WriteString(":")
	buffer.WriteString("limits")
	return buffer.String()
}

//验证手机号码
func valiPhoneNo(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//验证类型  0:注册；1：更改密码；2：修改绑定手机号；3：绑定紧急联系人；4：停用定位贴
func valiType(Type string) bool {
	return strings.Contains("01234", Type)
}

//验证是否是4为数
func valiVcode(vcode string) bool {
	reg := regexp.MustCompile(reqular2)
	return reg.MatchString(vcode)
}

func getCodeCheck(log *logger.Logger, code Code) error { //获取验证码的check
	log.Debug("method_start", "check", "check", code)
	if !valiPhoneNo(code.PhoneNo) || !valiType(code.Type) {
		msg := "Parameter Check Error"
		log.Debug("method_end", "check", "status", "fail", "msg", msg)
		return errors.New(msg)
	}
	log.Debug("method_end", "check", "status", "success")
	return nil
}

func checkCodeCheck(log *logger.Logger, code Code) error { //校验验证码的check
	log.Debug("method_start", "check", "check", code)
	if !valiPhoneNo(code.PhoneNo) || !valiType(code.Type) || !valiVcode(code.Vcode) {
		msg := "Parameter Check Error"
		log.Debug("method_end", "check", "status", "fail", "msg", msg)
		return errors.New(msg)
	}
	log.Debug("method_end", "check", "status", "success")
	return nil
}

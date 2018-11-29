package main

import (
	"errors"
	"fmt"
	logger "mtcomm/log"
	push "push/client"
	"strings"
	"time"
	UserClient "user/client"
)

// StringService provides operations on strings.
type TokenService interface {
	CreateToken(*Token) (string, error)
	DeleteToken(*Token) error
}

type tokenService struct{}

type Token struct {
	//IOS, WeChat, 其他
	PlatForm string
	UUID     string
	Token    string
	UserId   string
}

func (token *Token) String() string {
	return "PlatForm" + token.PlatForm + "UserId" + token.UserId + "UUID: " + token.UUID + ", Token: " + token.Token
}

func createTokenCheck(log *logger.Logger, token *Token) error {
	log.Debug("method_start", "createTokenCheck", "input", token)
	if token.PlatForm == "" || token.UUID == "" || token.UserId == "" {
		msg := "Parameter Check Error"
		log.Debug("method_end", "createTokenCheck", "status", "fail", "msg", msg)
		return errors.New(msg)
	}
	log.Debug("method_end", "createTokenCheck", "status", "success")
	return nil
}

func (service tokenService) CreateToken(token *Token) (string, error) {
	log.Debug("method_start", "createToken", "input", token)

	/* check */
	err := createTokenCheck(log, token)
	if err != nil {
		return "", err
	}
	c, err := redisClient.GetPipeline()
	if err != nil {
		return "", err
	}
	userId := token.UserId
	UUID := token.UUID

	u := &UserClient.UserRequest{UserId: userId}
	mysqlUser, err := userClient.SearchUser(u)
	if err != nil {
		return "", err
	}
	if mysqlUser.Err != "" {
		return "", errors.New(mysqlUser.Err)
	}
	//	redisMap, err := redisClient.HgetAllMap("U:" + userId) //redis上user信息
	//	if err != nil {
	//		return "", err
	//	}
	if mysqlUser.User.UserId != "" {
		flag1 := mysqlUser.User.PlatForm == ""
		flag2 := token.PlatForm == "weChat" || token.PlatForm == "WeChat" || token.PlatForm == "wechat"
		if !flag1 && !flag2 {
			//if mysqlUser.User.PlatForm != "weChat" || mysqlUser.User.PlatForm != "WeChat" || mysqlUser.User.PlatForm != "wechat" {
			if UUID != mysqlUser.User.UUID {
				//推送消息=========================================================
				key, err1 := redisClient.Get(serviceName + ":app:" + token.UserId)
				//log.Debug("========================================key:", key, "========================================err1:", err1.Error())
				if err1 != nil {
					if err1.Error() != "redigo: nil returned" {
						log.Debug("redisClient.Get(serviceName + : + token.UserId=================)", err1.Error())
						return "", err1
					}
				}
				redisClient.Del(serviceName + ":" + key)
				redisClient.Del(serviceName + ":app:" + token.UserId)
				list := []push.PushInfo{}
				pushInfo0 := push.PushInfo{
					Title: "您的账号在别的手机登录，您已被强制下线。",
					Text:  "{\"type\":\"U1\"}",
					Alias: getCid(userId), //23C0AD2F15F84E4A882156049C6FC7A732980060
				} //6685c1b99792c68fcef07029f3ded9c2 可用的

				list = append(list, pushInfo0)
				req := &push.PushRequest{
					PushType:  "1", //0代表推送通知，1代表透传消息
					PushInfos: list,
				}
				pushClient.PushErrorLog(req)
			}
		}
		if flag1 && flag2 {
			key, err1 := redisClient.Get(serviceName + ":WeChat:" + token.UserId)
			//log.Debug("========================================key:", key, "========================================err1:", err1.Error())
			if err1 != nil {
				if err1.Error() != "redigo: nil returned" {
					log.Debug("redisClient.Get(serviceName + : + token.UserId=================)", err1.Error())
					return "", err1
				}
			}
			redisClient.Del(serviceName + ":" + key)
			redisClient.Del(serviceName + ":WeChat:" + token.UserId)
		}
	}

	//从redis生成不可重复数
	incr, err := redisClient.Incr(serviceName + ":GETTOKEN")
	if err != nil {
		return "", err
	}

	//获取当前时间
	nowTime := fmt.Sprint(time.Now().Unix())
	tokenString := nowTime + fmt.Sprint(incr)
	//通过UserId,token双向保存
	if token.PlatForm == "IOS" {
		redisClient.PipelineSet(c, serviceName+":0:"+token.UUID+":"+tokenString, token.UserId)
		redisClient.PipelineSet(c, serviceName+":app:"+token.UserId, "0:"+token.UUID+":"+tokenString)
	} else if token.PlatForm == "WeChat" {
		redisClient.PipelineSet(c, serviceName+":1:"+token.UUID+":"+tokenString, token.UserId)
		redisClient.PipelineSet(c, serviceName+":WeChat:"+token.UserId, "1:"+token.UUID+":"+tokenString)
	} else {
		redisClient.PipelineSet(c, serviceName+":2:"+token.UUID+":"+tokenString, token.UserId)
		redisClient.PipelineSet(c, serviceName+":app:"+token.UserId, "2:"+token.UUID+":"+tokenString)
	}
	//保存user信息
	if token.PlatForm != "WeChat" {
		redisClient.PipelineHset(c, "U:"+token.UserId, "UUID", token.UUID)
		redisClient.PipelineHset(c, "U:"+token.UserId, "platForm", token.PlatForm)
	}
	err = redisClient.ExecutePipeline(c)
	if err != nil {
		log.Debug("redisClient.ExecutePipeline(c)++++++++++++++++++++++++++++++err", err.Error())
		return "", err
	}
	log.Debug("method_end", "createToken", "status", "success :token =", tokenString)
	return tokenString, nil
}

func deleteTokenCheck(log *logger.Logger, token *Token) error {
	log.Debug("method_start", "deleteTokenCheck", "input", token)
	if token.PlatForm == "" || token.UUID == "" || token.Token == "" {
		msg := "Parameter Check Error"
		log.Debug("method_end", "deleteTokenCheck", "status", "fail", "msg", msg)
		return errors.New(msg)
	}
	log.Debug("method_end", "deleteTokenCheck", "status", "success")
	return nil
}

func (service tokenService) DeleteToken(token *Token) error {
	log.Debug("method_start", "deleteToken", "input", token)

	/* check */
	err := deleteTokenCheck(log, token)
	if err != nil {
		return err
	}
	c, err2 := redisClient.GetPipeline()
	if err2 != nil {
		return err2
	}
	//del token
	if token.PlatForm == "IOS" {
		uId, err1 := redisClient.Get(serviceName + ":0:" + token.UUID + ":" + token.Token)
		token.UserId = uId
		if err1 != nil {
			if err1.Error() != "redigo: nil returned" {
				return err1
			}
		}
		if uId != "" {
			redisClient.PipelineDel(c, serviceName+":app:"+uId)
		}
		redisClient.PipelineDel(c, serviceName+":0:"+token.UUID+":"+token.Token)
	} else if token.PlatForm == "WeChat" {
		uId, err1 := redisClient.Get(serviceName + ":1:" + token.UUID + ":" + token.Token)
		if err1 != nil {
			return err1
		}
		if uId != "" {
			redisClient.PipelineDel(c, serviceName+":WeChat:"+uId)
		}
		redisClient.PipelineDel(c, serviceName+":1:"+token.UUID+":"+token.Token)
	} else {
		uId, err1 := redisClient.Get(serviceName + ":2:" + token.UUID + ":" + token.Token)
		if err1 != nil {
			return err1
		}
		if uId != "" {
			redisClient.PipelineDel(c, serviceName+":app:"+uId)
		}
		redisClient.PipelineDel(c, serviceName+":2:"+token.UUID+":"+token.Token)
	}

	//清除user里面的信息
	if token.PlatForm != "WeChat" {
		redisClient.PipelineHset(c, "U:"+token.UserId, "UUID", "")
		redisClient.PipelineHset(c, "U:"+token.UserId, "platForm", "")
	}
	err = redisClient.ExecutePipeline(c)
	if err != nil {
		return err
	}
	log.Debug("method_end", "deleteToken", "status", "success")
	return nil
}

func getCid(userId string) string {
	//推送的别名设置，由于推送别名不能超过40字节。
	//推送别名为：UUID（去掉横杆）后32位 + UserId的后8位。
	//目前UUID正好是32位，但是不确定以后会不会变动，所以判断超过32位的话取后32位
	UUID, _ := redisClient.Hget("U:"+userId, "UUID")
	if len(UUID) > 32 {
		UUID = UUID[len(UUID)-32 : len(UUID)]
	}
	returnString := ""
	UUIDs := strings.Split(UUID, "-")
	for _, v := range UUIDs {
		returnString += v
	}
	lenUId := len(userId)
	if lenUId > 8 {
		returnString += userId[lenUId-8 : lenUId]
	}
	return returnString
}

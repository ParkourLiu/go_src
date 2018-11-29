package main

import (
	"strconv"
	"strings"
)

// StringService provides operations on strings.
type PushService interface {
	Push(pushInfoList *pushInfoList) (string, error)
}

type pushService struct{}

type pushInfoList struct {
	PushType  string     `json:"pushType"` //0代表通知消息，1代表透传消息,all_0全量推送通知，all_1全量推送透传消息，
	PushInfos []pushInfo `json:"pushInfos"`
}

type pushInfo struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Alias    string `json:"alias"`
	JsonInfo string `json:"jsonInfo"`
}

func (pushInfoList *pushInfoList) String() string {
	str := "PushType:" + pushInfoList.PushType
	list := pushInfoList.PushInfos
	for i := 0; i < len(list); i++ {
		str += "title:" + strconv.Itoa(i) + ":" + list[i].Title + ",text" + strconv.Itoa(i) + ":" + list[i].Text + ",alias" + strconv.Itoa(i) + ":" + list[i].Alias
	}
	return str
}

func getAlias(userId string) string {
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
	returnString += userId[lenUId-8 : lenUId]
	return returnString
}

func (service pushService) Push(pushInfoList *pushInfoList) (string, error) {
	log.Debug("method_start", "Push", "input", pushInfoList)

	/* service */
	msg := ""
	if pushInfoList.PushType == "0" { //通知消息
		msg, _ = PushList(pushInfoList)

	} else if pushInfoList.PushType == "1" { //透传消息
		msg, _ = OSPFList(pushInfoList)

	} else if pushInfoList.PushType == "all_0" { //给所有人推送通知
		msg, _ = PushAll(pushInfoList)
	} /*else if pushInfoList.PushType == "all_1" { //给所有人推送透传消息
		pushList := allUser(pushInfoList.PushInfos[0])
		msg, _ = OSPFList(pushList)
	}*/ //此处打开即可给所有人推送，但是及其影响性能，慎用
	log.Debug("method_end", "Push", "status", "success")
	return msg, nil
}

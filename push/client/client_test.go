package client_test

import (
	"context"
	"mtcomm/db/redis"
	logger "mtcomm/log"
	push "push/client"
	"strings"
	"testing"
)

var (
	pushClient  push.PushCaller
	redisClient redis.RedisClient
)

func TestClient(t *testing.T) {
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:       context.TODO(),
		Logger:    logger.GetDefaultLogger(),
		RedisHost: "127.0.0.1:6379",
		//RedisHost:     "106.14.216.4:8888",
		//RedisPassword: "zaq12wsx1",
	})
	pushClient = push.NewPushCaller(redisClient)

	list := []push.PushInfo{}
	pushInfo0 := push.PushInfo{
		Title:    "全部推送",
		Text:     "iiiii", //当text参数不必要时，需要传一个空格
		JsonInfo: "{\"id\":\"1·234\",\"name\":\"5678\"}",
		Alias:    getAlias("U15064872359141091"), //U14921352095531169 //U146970017247836
	} //6685c1b99792c68fcef07029f3ded9c2 可用的

	list = append(list, pushInfo0)
	req := &push.PushRequest{
		PushType:  "all_0", //0代表推送消息，1代表透传消息,all_0代表给所有用户推送通知
		PushInfos: list,
	}
	pushClient.PushErrorLog(req)
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

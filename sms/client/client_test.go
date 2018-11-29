package client_test

import (
	"context"
	"fmt"
	"mtcomm/db/redis"
	logger "mtcomm/log"
	sms "sms/client"
	"testing"
)

var (
	smsClient   sms.SmsCaller
	redisClient redis.RedisClient
)

func TestClient(t *testing.T) {
	redisClient = redis.NewRedisClient(&redis.RedisServerInfo{
		Ctx:    context.TODO(),
		Logger: logger.GetDefaultLogger(),
		//RedisHost: "127.0.0.1:6379",
		RedisHost:     "106.15.156.236:8888",
		RedisPassword: "zaq12wsx1",
	})
	smsClient = sms.NewSmsCaller(redisClient)
	param := []string{"1111"}
	phoneNos := []string{"17671774535"}
	req := &sms.SmsRequest{
		Params: param,
		Mobile: phoneNos,
		Tpl_id: "57981",
	}
	err := smsClient.Sms(req)
	fmt.Println("err::::::::::::::::::::", err)
}

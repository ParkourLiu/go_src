package producer_test

import (
	"context"
	"fmt"
	"mtcomm/db/redis"
	logger "mtcomm/log"
	csr "mtcomm/queue/consumer"
	prd "mtcomm/queue/producer"
	"testing"
	"time"
)

var client redis.RedisClient

func init() {
	logger.SetDefaultLogLevel(logger.LevelDebug)
	info := &redis.RedisServerInfo{
		Ctx:       context.TODO(),
		Logger:    logger.GetDefaultLogger(),
		RedisHost: "127.0.0.1:6379",
	}
	client = redis.NewRedisClient(info)
}

func TestMq(t *testing.T) {
	cc1 := csr.NewConsumer(client)
	cc2 := csr.NewConsumer(client)
	pc := prd.NewProducer(client)
	pc.ProduceData("test", "1")
	pc.ProduceData("test", "2")
	pc.ProduceData("test", "3")
	pc.ProduceData("test", "4")
	pc.ProduceData("test", "5")
	pc.ProduceData("test", "6")

	callback := func(value string) error {
		fmt.Println(value)
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	go cc1.ConsumeData("test", callback)
	go cc2.ConsumeData("test", callback)
	time.Sleep(5 * time.Second)
	return
}

package client

import (
	"encoding/json"
	"mtcomm/db/redis"
	prd "mtcomm/queue/producer"
)

type PushCaller interface {
	PushErrorLog(req *PushRequest) error
}

type pushCaller struct {
	pp prd.Producer
}

func NewPushCaller(redisClient redis.RedisClient) PushCaller {
	pp := prd.NewProducer(redisClient)
	return &pushCaller{pp: pp}
}

type PushRequest struct {
	PushType  string     `json:"pushType"` //0代表通知消息，1代表透传消息
	PushInfos []PushInfo `json:"pushInfos"`
}

type PushInfo struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Alias    string `json:"alias"`
	JsonInfo string `json:"jsonInfo"`
}

func (m *pushCaller) PushErrorLog(req *PushRequest) error {
	data, err := json.Marshal(*req)
	if err != nil {
		return err
	}
	return m.pp.ProduceData("push", string(data))
}

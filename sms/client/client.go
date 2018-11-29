package client

import (
	"encoding/json"
	"mtcomm/db/redis"
	prd "mtcomm/queue/producer"
)

type SmsCaller interface {
	Sms(req *SmsRequest) error
}

type smsCaller struct {
	pp prd.Producer
}

func NewSmsCaller(redisClient redis.RedisClient) SmsCaller {
	pp := prd.NewProducer(redisClient)
	return &smsCaller{pp: pp}
}

type SmsRequest struct {
	Params []string `json:"params"`
	Mobile []string `json:"mobile"`
	Tpl_id string   `json:"tpl_id"`
}

func (m *smsCaller) Sms(req *SmsRequest) error {
	data, err := json.Marshal(*req)
	if err != nil {
		return err
	}
	return m.pp.ProduceData("sms", string(data))
}

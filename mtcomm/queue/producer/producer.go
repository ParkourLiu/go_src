package producer

import "mtcomm/db/redis"

type Producer interface {
	ProduceData(queueName string, struction string) error
}

type producer struct {
	Client redis.RedisClient
}

func NewProducer(c redis.RedisClient) Producer {
	return &producer{Client: c}
}

func (p *producer) ProduceData(queueName string, data string) error {
	return p.Client.Rpush("queue:"+queueName, data)
}

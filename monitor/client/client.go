package client

import (
	"encoding/json"
	"mtcomm/db/redis"
	prd "mtcomm/queue/producer"
)

type MonitorCaller interface {
	PushErrorLog(req *MonitorRequest) error
}

type monitorCaller struct {
	pp prd.Producer
}

func NewMonitorCaller(redisClient redis.RedisClient) MonitorCaller {
	pp := prd.NewProducer(redisClient)
	return &monitorCaller{pp: pp}
}

type MonitorRequest struct {
	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`
	Param       string `json:"param"`
	ErrorMsg    string `json:"errorMsg"`
	ErrorTime   string `json:"errorTime"`
}

func (m *monitorCaller) PushErrorLog(req *MonitorRequest) error {
	data, err := json.Marshal(*req)
	if err != nil {
		return err
	}
	return m.pp.ProduceData("monitor", string(data))
}

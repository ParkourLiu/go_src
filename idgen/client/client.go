package client

import (
	caller "idgen/caller"
	k8s "mtcomm/k8s"
	logger "mtcomm/log"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type IdGenClient interface {
	GetUniqueId() string
}

type idGenClient struct {
	Caller caller.IdGeneraterCaller
	log    *logger.Logger
}

func NewIdGenClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace string) IdGenClient {
	if !k8sClient.IsClusterEnv() {
		return nil
	}
	log := logger.GetDefaultLogger()
	c := caller.NewIdGeneraterCaller(k8sClient, tracer, namespace, "idgenClient", "idgen", "8888")
	client := &idGenClient{
		Caller: c,
		log:    log,
	}
	return client
}

func (client *idGenClient) GetUniqueId() string {
	ids := getIdsFormRemote(client.log, client.Caller, uint32(1))
	return ids[0]
}

func getIdsFormRemote(log *logger.Logger, c caller.IdGeneraterCaller, count uint32) []string {
	req := &caller.GenerateUniqueIdV1Request{
		Count: count,
	}
	resp, err := c.GenerateUniqueIdV1(req)
	if err != nil {
		log.Warn("method", "getIdsFormRemote", "msg", err, "retry", "will retry after 10s.")
		time.Sleep(10 * time.Second)
		return getIdsFormRemote(log, c, count)
	}
	return resp.Ids
}

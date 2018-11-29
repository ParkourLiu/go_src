package main

import (
	"errors"
	"fmt"
	logger "mtcomm/log"
	"runtime"
)

type GatewayService interface {
	Forward(jsonData *reqData) (string, error)
}

type gatewayService struct{}

func (gatewayService) Forward(jsonData *reqData) (result string, err error) {
	defer func() {
		if x := recover(); x != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			msg := fmt.Sprintf("Panic Msg: %v\n%s", x, buf)

			log := logger.GetDefaultLogger()
			log.Error("panic", "true", "msg", msg)
			result, err = "", errors.New("Gateway Forward Panic! "+msg)
		}
	}()

	caller := NewProxyCaller(k8sClient, tracer, namespace, "gateway", jsonData.CalledServiceName, "8888")
	result, err = caller.Forward(jsonData)
	return
}

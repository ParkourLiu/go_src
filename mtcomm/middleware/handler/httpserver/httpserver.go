package httpserver

import (
	"context"
	"fmt"
	monitor "monitor/client"
	logger "mtcomm/log"
	"net/http"
	"runtime"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
)

type HttpMtalkServer struct {
	*httptransport.Server
	ServiceName   string
	MethodName    string
	MonitorClient monitor.MonitorCaller
	PanicResponse interface{} //err 为 panic, 其他字段为空即可。
}

func NewHttpMtalkServer(httpServer *httptransport.Server, serviceName string, methodName string, monitorClient monitor.MonitorCaller, panicResponse interface{}) *HttpMtalkServer {
	s := &HttpMtalkServer{
		ServiceName:   serviceName,
		MethodName:    methodName,
		MonitorClient: monitorClient,
		PanicResponse: panicResponse,
	}
	s.Server = httpServer
	return s
}

func (s HttpMtalkServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			msg := fmt.Sprintf("Panic Msg: %v\n%s", x, buf)

			log := logger.GetDefaultLogger()
			log.Error("panic", "true", "msg", msg)
			//monitor
			go func() {
				req := &monitor.MonitorRequest{
					ServiceName: s.ServiceName,
					MethodName:  s.MethodName,
					Param:       "unknown",
					ErrorMsg:    msg,
					ErrorTime:   time.Now().String(),
				}
				s.MonitorClient.PushErrorLog(req)
			}()
			// 设置response
			httptransport.EncodeJSONResponse(context.TODO(), w, s.PanicResponse)
		}
	}()
	s.Server.ServeHTTP(w, r)
}

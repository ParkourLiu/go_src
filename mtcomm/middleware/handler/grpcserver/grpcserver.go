package grpcserver

import (
	monitor "monitor/client"
	logger "mtcomm/log"
	"time"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
)

type GrpcMtalkServer struct {
	*grpctransport.Server
	ServiceName   string
	MethodName    string
	MonitorClient monitor.MonitorCaller
}

func NewGrpcMtalkServer(grpcServer *grpctransport.Server, serviceName string, methodName string, monitorClient monitor.MonitorCaller) *GrpcMtalkServer {
	s := &GrpcMtalkServer{
		ServiceName:   serviceName,
		MethodName:    methodName,
		MonitorClient: monitorClient,
	}
	s.Server = grpcServer
	return s
}

func (s GrpcMtalkServer) ServeGRPC(ctx oldcontext.Context, req interface{}) (retctx oldcontext.Context, resp interface{}, err error) {
	defer func() {
		if x := recover(); x != nil {
			log := logger.GetDefaultLogger()
			log.Error("panic", "true")
			//monitor
			go func() {
				req := &monitor.MonitorRequest{
					ServiceName: s.ServiceName,
					MethodName:  s.MethodName,
					Param:       "unknown",
					ErrorMsg:    "panic occurr!",
					ErrorTime:   time.Now().String(),
				}
				s.MonitorClient.PushErrorLog(req)
			}()
		}
	}()
	return s.Server.ServeGRPC(ctx, req)
}

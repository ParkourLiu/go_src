package main

import (
	"strings"
	"time"
	monitor "monitor/client"
)

type loggingMiddleware struct {
	next GatewayService
}

func (mw loggingMiddleware) Forward(jsonData *reqData) (result string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "Forward",
		"input", jsonData,
	)

	defer func(begin time.Time) {
		noErrStr := "\"err\":\"\""
		ok := strings.Contains(result, noErrStr)
		if ok && err==nil {
			log.Info(
				"write_log", "yes",
				"method_end", "Forward",
				"input", jsonData,
//				"return", result,
				"status", "success",
				"took", time.Since(begin),
			)
		}else if err != nil {
			log.Error(
				"write_log", "yes",
				"method_end", "Forward",
				"input", jsonData,
				"msg", err,
				"status", "fail",
				"took", time.Since(begin),
			)
			//monitor
			go func() {
				req := &monitor.MonitorRequest{
					ServiceName: serviceName,
					MethodName:  "Forward",
					Param:       jsonData.String(),
					ErrorMsg:    err.Error(),
					ErrorTime:   time.Now().String(),
				}
				monitorClient.PushErrorLog(req)
			}()
		}else{
			log.Error(
				"write_log", "yes",
				"method_end", "Forward",
				"input", jsonData,
				"return", result,
				"status", "fail",
				"took", time.Since(begin),
			)
			//monitor
			go func() {
				req := &monitor.MonitorRequest{
					ServiceName: serviceName,
					MethodName:  "Forward",
					Param:       jsonData.String(),
					ErrorMsg:    result,
					ErrorTime:   time.Now().String(),
				}
				monitorClient.PushErrorLog(req)
			}()
		}
	}(time.Now())
	
	return mw.next.Forward(jsonData)
}

package main

import (
	"time"
)

type loggingMiddleware struct {
	next PushService
}

func (mw loggingMiddleware) Push(pushInfoList *pushInfoList) (msg string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "PushService",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "Push",
				"input", pushInfoList,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "Push",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	msg, err = mw.next.Push(pushInfoList)
	return
}

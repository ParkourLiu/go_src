package main

import (
	"time"
)

type loggingMiddleware struct {
	next SmsService
}

func (mw loggingMiddleware) Sms(info *SmsInfo) (msg string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SmsService",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "Sms",
				"input", info,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "Sms",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	msg, err = mw.next.Sms(info)
	return
}

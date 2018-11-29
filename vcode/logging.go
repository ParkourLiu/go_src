package main

import (
	"time"

	logger "mtcomm/log"
)

type loggingMiddleware struct {
	next CodeService
}

func (mw loggingMiddleware) CheckCode(code Code) (str string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "CheckCode",
		"code", code,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "CheckCode",
				"code", code,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "CheckCode",
				"code", code,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	str, err = mw.next.CheckCode(code)
	return
}
func (mw loggingMiddleware) GetCode(code Code) (str string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "GetCode",
		"code", code,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "GetCode",
				"code", code,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "GetCode",
				"code", code,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	str, err = mw.next.GetCode(code)
	return
}

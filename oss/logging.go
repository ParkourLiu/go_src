package main

import (
	"time"
)

type loggingMiddleware struct {
	next Oss
}

func (mw loggingMiddleware) GetOssTokenForWeb() (ossToken map[string]interface{}, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "GetOssTokenForWeb",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "GetOssTokenForWeb",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "GetOssTokenForWeb",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	ossToken, err = mw.next.GetOssTokenForWeb()
	return
}

package main

import (
	"time"
)

type loggingMiddleware struct {
	next IdGeneraterService
}

func (mw loggingMiddleware) GenerateUniqueIdV1(count uint32) (ids []string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "GenerateUniqueIdV1",
		"input", count,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "GenerateUniqueIdV1",
				"input", count,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "GenerateUniqueIdV1",
				"input", count,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	ids, err = mw.next.GenerateUniqueIdV1(count)
	return
}

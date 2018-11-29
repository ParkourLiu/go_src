package main

import (
	"time"
)

type loggingMiddleware struct {
	next VersionService
}

func (mw loggingMiddleware) VersionInfo(version *Version) (requestList map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "VersionInfo",
		"input", version,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "VersionInfo",
				"input", version,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "VersionInfo",
				"input", version,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	requestList, code, err = mw.next.VersionInfo(version)
	return
}

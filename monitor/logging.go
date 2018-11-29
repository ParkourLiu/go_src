package main

import (
	"time"
)

type loggingMiddleware struct {
	next MonitorService
}

func (mw loggingMiddleware) monitor(m *monitorRequest) {
	log.Info(
		"check_log", "yes",
		"method_start", "monitor",
		"input", m,
	)
	defer func(begin time.Time) {
		log.Info(
			"check_log", "yes",
			"method_end", "monitor",
			"input", m,
			"status", "success",
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.next.monitor(m)
	return
}

package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           VersionService
}

func (mw instrumentingMiddleware) VersionInfo(version *Version) (request map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "VersionInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	request, code, err = mw.next.VersionInfo(version)
	return
}

package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Oss
}

func (mw instrumentingMiddleware) GetOssTokenForWeb() (ossToken map[string]interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetOssTokenForWeb", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ossToken, err = mw.next.GetOssTokenForWeb()
	return
}

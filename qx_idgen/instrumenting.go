package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           IdGeneraterService
}

func (mw instrumentingMiddleware) GenerateUniqueIdV1(count uint32) (ids []string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GenerateUniqueIdV1", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ids, err = mw.next.GenerateUniqueIdV1(count)
	return
}

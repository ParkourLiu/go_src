package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           CodeService
}

func (mw instrumentingMiddleware) CheckCode(code Code) (str string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CheckCode", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	str, err = mw.next.CheckCode(code)
	return
}

func (mw instrumentingMiddleware) GetCode(code Code) (str string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetCode", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	str, err = mw.next.GetCode(code)
	return
}

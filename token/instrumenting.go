package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           TokenService
}

func (mw instrumentingMiddleware) CreateToken(token *Token) (tk string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "IdCreate", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	tk, err = mw.next.CreateToken(token)
	return
}

func (mw instrumentingMiddleware) DeleteToken(token *Token) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "deleteToken", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteToken(token)
	return
}

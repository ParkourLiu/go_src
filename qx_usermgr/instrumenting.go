package main

import (
	"fmt"
	"time"

	U "qx_user/caller"
	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           UserService
}

func (mw instrumentingMiddleware) RegAndLogin(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "RegAndLogin", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.RegAndLogin(ur)
	return
}

func (mw instrumentingMiddleware) OtherLogin(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "OtherLogin", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.OtherLogin(ur)
	return
}

func (mw instrumentingMiddleware) SearchUser(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.SearchUser(ur)
	return
}

func (mw instrumentingMiddleware) UpdateUser(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.UpdateUser(ur)
	return
}

func (mw instrumentingMiddleware) ChangeBind(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ChangeBind", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.ChangeBind(ur)
	return
}

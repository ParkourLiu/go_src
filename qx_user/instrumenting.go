package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           UserService
}

func (mw instrumentingMiddleware) SearchUserById(user *User) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchUserById", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.SearchUserById(user)
	return
}
func (mw instrumentingMiddleware) SearchUsers(user *User) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchUsers", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.SearchUsers(user)
	return
}
func (mw instrumentingMiddleware) AddUser(user *User) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchUserById", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.AddUser(user)
	return
}
func (mw instrumentingMiddleware) UpdateUser(user *User) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchUserById", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.UpdateUser(user)
	return
}

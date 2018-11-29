package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           AutismService
}

func (mw instrumentingMiddleware) StarDetails(autism *Autism) (returnMap map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "StarDetails", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	returnMap, code, err = mw.next.StarDetails(autism)
	return
}

func (mw instrumentingMiddleware) SaveComment(autism *Autism) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveComment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SaveComment(autism)
	return
}

//==============================================================================================
func (mw instrumentingMiddleware) StarList(autism *Autism) (brightCount string, likeCount string, rolls []map[string]string, starList []map[string]string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "StarList", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	brightCount, likeCount, rolls, starList, err = mw.next.StarList(autism)
	return
}
func (mw instrumentingMiddleware) Likes(autism *Autism) (err error) {
	log.Debug("method_start", "Likes")
	err = mw.next.Likes(autism)
	log.Debug("method_end", "Likes", "status", "success")
	return
}

func (mw instrumentingMiddleware) GetUnionid(autism *Autism) (data Un_ionid, err error) {
	log.Debug("method_start", "GetUnionid")
	data, err = mw.next.GetUnionid(autism)
	log.Debug("method_end", "GetUnionid", "status", "success")
	return
}

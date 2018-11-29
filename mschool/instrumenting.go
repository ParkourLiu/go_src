package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           MschoolService
}

func (mw instrumentingMiddleware) CreateSchoolYear(school *School) (request map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateSchoolYear", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	request, code, err = mw.next.CreateSchoolYear(school)
	return
}

func (mw instrumentingMiddleware) SearchSchool(school *School) (request map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchSchool", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	request, code, err = mw.next.SearchSchool(school)
	return
}

func (mw instrumentingMiddleware) MySchool(school *School) (request [][]map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "MySchool", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	request, code, err = mw.next.MySchool(school)
	return
}

func (mw instrumentingMiddleware) SetWorkDay(school *School) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SetWorkDay", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SetWorkDay(school)
	return
}
func (mw instrumentingMiddleware) LookWorkDay(school *School) (request []map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "LookWorkDay", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	request, code, err = mw.next.LookWorkDay(school)
	return
}

func (mw instrumentingMiddleware) WorkDay(school *School) (request []map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "LookWorkDay", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	request, code, err = mw.next.WorkDay(school)
	return
}

func (mw instrumentingMiddleware) UpSchoolYear(school *School) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpSchoolYear", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.UpSchoolYear(school)
	return
}

func (mw instrumentingMiddleware) FaceDataGather() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FaceDataGather", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.FaceDataGather()
	return
}

func (mw instrumentingMiddleware) LabelGuidDataGather() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "LabelGuidDataGather", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.LabelGuidDataGather()
	return
}

func (mw instrumentingMiddleware) GetDataFileUrl(school *School) (returnMap map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetDataFileUrl", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	returnMap, code, err = mw.next.GetDataFileUrl(school)
	return
}
func (mw instrumentingMiddleware) DelDataFile(school *School) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "LabelGuidDataGather", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.DelDataFile(school)
	return
}

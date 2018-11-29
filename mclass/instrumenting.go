package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           McreateClassService
}

func (mw instrumentingMiddleware) SelectByScInviteCode(sy *SchoolYear) (schoolYear map[string]interface{},err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SelectByScInviteCode", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	schoolYear, err, code = mw.next.SelectByScInviteCode(sy)
	return
}
func (mw instrumentingMiddleware) CreateClass(sy *SchoolYear) (schoolYear map[string]string, err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateClass", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	schoolYear, err, code = mw.next.CreateClass(sy)
	return
}

func (mw instrumentingMiddleware) TeacherJoinClass(sy *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "TeacherJoinClass", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	 err, code = mw.next.TeacherJoinClass(sy)
	return
}
func (mw instrumentingMiddleware) FamilyJoinClass(sy *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FamilyJoinClass", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err, code = mw.next.FamilyJoinClass(sy)
	return
}
func (mw instrumentingMiddleware) NewMember(gp *SchoolYear) (data []map[string]string, err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewMember", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, err, code = mw.next.NewMember(gp)
	return
}
func (mw instrumentingMiddleware) ApproveMembers(gp *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ApproveMembers", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err, code = mw.next.ApproveMembers(gp)
	return
}
func (mw instrumentingMiddleware) ManagerMember(gp *SchoolYear) (data []map[string]string, err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ManagerMember", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, err, code = mw.next.ManagerMember(gp)
	return
}
func (mw instrumentingMiddleware) OperateMember(gp *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "OperateMember", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err, code = mw.next.OperateMember(gp)
	return
}
func (mw instrumentingMiddleware) ClassQrCode(gp *SchoolYear) (data map[string]string,err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ClassQrCode", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data,err, code = mw.next.ClassQrCode(gp)
	return
}
func (mw instrumentingMiddleware) UpdateStudent(gp *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateStudent", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err, code = mw.next.UpdateStudent(gp)
	return
}
func (mw instrumentingMiddleware) FindAllMember(gp *SchoolYear) (data []map[string]string ,err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindAllMember", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data,err, code = mw.next.FindAllMember(gp)
	return
}
func (mw instrumentingMiddleware) FindTeacherMember(gp *SchoolYear) (data []map[string]string ,err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindTeacherMember", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data,err, code = mw.next.FindTeacherMember(gp)
	return
}
func (mw instrumentingMiddleware) UpdateGroupChatInfo(gp *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateGroupChatInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err, code = mw.next.UpdateGroupChatInfo(gp)
	return
}
func (mw instrumentingMiddleware) UpdateTeachInfo(gp *SchoolYear) (err error, code string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateTeachInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err, code = mw.next.UpdateTeachInfo(gp)
	return
}
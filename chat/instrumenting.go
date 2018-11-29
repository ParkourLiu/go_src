package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           RYToKenService
}

func (mw instrumentingMiddleware) GetRYToken(info *ToKenInfo) (data interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetRYToken", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.GetRYToken(info)
	return data, code, err
}
func (mw instrumentingMiddleware) GetUserInfo(info *ToKenInfo) (data interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.GetUserInfo(info)
	return data, code, err
}

func (mw instrumentingMiddleware) CreateGroupChat(info *GroupChatInfo) (data string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateGroupChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.CreateGroupChat(info)
	return data, code, err
}
func (mw instrumentingMiddleware) JoinGroupChat(info *GroupChatInfo) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "JoinGroupChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.JoinGroupChat(info)
	return code, err
}
func (mw instrumentingMiddleware) QuitGroupChat(info *GroupChatInfo) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "QuitGroupChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.QuitGroupChat(info)
	return code, err
}

func (mw instrumentingMiddleware) Dismiss(info *GroupChatInfo) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Dismiss", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.Dismiss(info)
	return code, err
}
func (mw instrumentingMiddleware) QueryGroupChatMemberList(info *GroupChatInfo) (data map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "QueryGroupChatMemberList", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.QueryGroupChatMemberList(info)
	return data, code, err
}
func (mw instrumentingMiddleware) GetArrayGroupInfo(info *GroupChatInfo) (data interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetArrayGroupInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.GetArrayGroupInfo(info)
	return data, code, err
}
func (mw instrumentingMiddleware) GetMyGroupChat(info *GroupChatInfo) (data interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetMyGroupChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.GetMyGroupChat(info)
	return data, code, err
}
func (mw instrumentingMiddleware) UpdateGroupChat(info *GroupChatInfo) (data interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateGroupChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.UpdateGroupChat(info)
	return data, code, err
}

//parkour======================================================start
func (mw instrumentingMiddleware) AddOfficialMSG(chat *Chat) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddOfficialMSG", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.AddOfficialMSG(chat)
	return
}
func (mw instrumentingMiddleware) LookOfficialMSG() (data []map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "LookOfficialMSG", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	data, code, err = mw.next.LookOfficialMSG()
	return
}
func (mw instrumentingMiddleware) InformChat(chat *Chat) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "InformChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.InformChat(chat)
	return
}

//parkour======================================================end

func (mw instrumentingMiddleware) SearchChatInfo(info *GroupChatInfo) (data map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchChatInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.SearchChatInfo(info)
	return data, code, err
}
func (mw instrumentingMiddleware) GetClassId(info *GroupChatInfo) (data string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SearchChatInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.GetClassId(info)
	return data, code, err
}

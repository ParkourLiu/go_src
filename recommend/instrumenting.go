package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Recommend
}

func (mw instrumentingMiddleware) PopuserUsers(pageInfo *PageInfo) (retuenList []map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PopuserUsers", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	retuenList, code, err = mw.next.PopuserUsers(pageInfo)
	return
}
func (mw instrumentingMiddleware) RmduserUsers() (retuenList []map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "RmduserUsers", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	retuenList, code, err = mw.next.RmduserUsers()
	return
}

func (mw instrumentingMiddleware) HotDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HotDynamic", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.HotDynamic(pageInfo)
	return
}

func (mw instrumentingMiddleware) NewDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewDynamic", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.NewDynamic(pageInfo)
	return
}

func (mw instrumentingMiddleware) ReverseHotDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HotDynamic", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.ReverseHotDynamic(pageInfo)
	return
}

func (mw instrumentingMiddleware) ReverseNewDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "NewDynamic", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.ReverseNewDynamic(pageInfo)
	return
}

func (mw instrumentingMiddleware) SavePopuserUsers() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SavePopuserUsers", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SavePopuserUsers()
	return
}

func (mw instrumentingMiddleware) SaveRmduserUsers() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveRmduserUsers", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SaveRmduserUsers()
	return
}

func (mw instrumentingMiddleware) SaveHotDynamic() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveHotDynamic", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SaveHotDynamic()
	return
}

func (mw instrumentingMiddleware) FriendRecommend(pageInfo *PageInfo) (data []map[string]interface{}, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FriendRecommend", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.FriendRecommend(pageInfo)
	return
}

func (mw instrumentingMiddleware) SaveFriendRecommend() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveFriendRecommend", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SaveFriendRecommend()
	return
}

func (mw instrumentingMiddleware) SearchRecommend() (data []string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveFriendRecommend", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, code, err = mw.next.SearchRecommend()
	return
}

func (mw instrumentingMiddleware) PushFavour() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PushFavour", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.PushFavour()
	return
}

func (mw instrumentingMiddleware) PushStartActivity() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PushStartActivity", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.PushStartActivity()
	return
}

func (mw instrumentingMiddleware) SaveHomePageCache() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SaveHomePageCache", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.SaveHomePageCache()
	return
}

func (mw instrumentingMiddleware) AddDynamicFans() (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddDynamicFans", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.AddDynamicFans()
	return
}

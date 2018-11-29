package main

import (
	"fmt"
	"time"
)

type loggingMiddleware struct {
	next Recommend
}

func (mw loggingMiddleware) PopuserUsers(pageInfo *PageInfo) (retuenList []map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "PopuserUsers",
		"input", pageInfo,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "PopuserUsers",
				"input", pageInfo,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "PopuserUsers",
				"input", pageInfo,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	retuenList, code, err = mw.next.PopuserUsers(pageInfo)
	return
}

func (mw loggingMiddleware) RmduserUsers() (retuenList []map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "RmduserUsers",
		"input", "nil",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "RmduserUsers",
				"input", "nil",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "RmduserUsers",
				"input", "nil",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	retuenList, code, err = mw.next.RmduserUsers()
	return
}

func (mw loggingMiddleware) HotDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "HotDynamic",
		"input", pageInfo,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "HotDynamic",
				"input", pageInfo,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "HotDynamic",
				"input", pageInfo,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.HotDynamic(pageInfo)
	return
}

func (mw loggingMiddleware) NewDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "NewDynamic",
		"input", pageInfo,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "NewDynamic",
				"input", pageInfo,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "NewDynamic",
				"input", pageInfo,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.NewDynamic(pageInfo)
	return
}

func (mw loggingMiddleware) ReverseHotDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "ReverseHotDynamic",
		"input", pageInfo,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ReverseHotDynamic",
				"input", pageInfo,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ReverseHotDynamic",
				"input", pageInfo,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.ReverseHotDynamic(pageInfo)
	return
}

func (mw loggingMiddleware) ReverseNewDynamic(pageInfo *PageInfo) (retuenList []map[string]interface{}, pageFlag string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "ReverseNewDynamic",
		"input", pageInfo,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ReverseNewDynamic",
				"input", pageInfo,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ReverseNewDynamic",
				"input", pageInfo,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	retuenList, pageFlag, code, err = mw.next.ReverseNewDynamic(pageInfo)
	return
}

func (mw loggingMiddleware) SavePopuserUsers() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SavePopuserUsers",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SavePopuserUsers",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SavePopuserUsers",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SavePopuserUsers()
	return
}

func (mw loggingMiddleware) SaveRmduserUsers() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SaveRmduserUsers",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SaveRmduserUsers",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SaveRmduserUsers",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SaveRmduserUsers()
	return
}

func (mw loggingMiddleware) SaveHotDynamic() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SaveHotDynamic",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SaveHotDynamic",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SaveHotDynamic",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SaveHotDynamic()
	return
}

func (mw loggingMiddleware) FriendRecommend(pageInfo *PageInfo) (data []map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "FriendRecommend",
		"input", fmt.Sprint(pageInfo),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "FriendRecommend",
				"input", fmt.Sprint(pageInfo),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SaveHotDynamic",
				"input", fmt.Sprint(pageInfo),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.FriendRecommend(pageInfo)
	return
}

func (mw loggingMiddleware) SaveFriendRecommend() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SaveFriendRecommend",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SaveFriendRecommend",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SaveFriendRecommend",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SaveFriendRecommend()
	return
}

func (mw loggingMiddleware) SearchRecommend() (data []string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SearchRecommend",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SearchRecommend",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SearchRecommend",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.SearchRecommend()
	return
}

func (mw loggingMiddleware) PushFavour() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "PushFavour",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "PushFavour",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "PushFavour",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.PushFavour()
	return
}

func (mw loggingMiddleware) PushStartActivity() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "PushStartActivity",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "PushStartActivity",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "PushStartActivity",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.PushStartActivity()
	return
}

func (mw loggingMiddleware) SaveHomePageCache() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SaveHomePageCache",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SaveHomePageCache",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SaveHomePageCache",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SaveHomePageCache()
	return
}

func (mw loggingMiddleware) AddDynamicFans() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "AddDynamicFans",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "AddDynamicFans",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "AddDynamicFans",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.AddDynamicFans()
	return
}

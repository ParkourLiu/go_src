package main

import (
	"time"

	logger "mtcomm/log"
)

type loggingMiddleware struct {
	next RYToKenService
}

func (mw loggingMiddleware) GetRYToken(info *ToKenInfo) (data interface{}, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "GetRYToken",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "GetRYToken",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "GetRYToken",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.GetRYToken(info)
	return data, code, err
}
func (mw loggingMiddleware) GetUserInfo(info *ToKenInfo) (data interface{}, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "GetUserInfo",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "GetUserInfo",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "GetUserInfo",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.GetUserInfo(info)
	return data, code, err
}

func (mw loggingMiddleware) CreateGroupChat(info *GroupChatInfo) (data string, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "CreateGroupChat",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "CreateGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "CreateGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.CreateGroupChat(info)
	return data, code, err
}

func (mw loggingMiddleware) JoinGroupChat(info *GroupChatInfo) (code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "JoinGroupChat",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "JoinGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "JoinGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.JoinGroupChat(info)
	return code, err
}

func (mw loggingMiddleware) QuitGroupChat(info *GroupChatInfo) (code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "QuitGroupChat",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "QuitGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "QuitGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.QuitGroupChat(info)
	return code, err
}

func (mw loggingMiddleware) Dismiss(info *GroupChatInfo) (code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "Dismiss",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "Dismiss",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "Dismiss",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.Dismiss(info)
	return code, err
}

func (mw loggingMiddleware) QueryGroupChatMemberList(info *GroupChatInfo) (data map[string]interface{}, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "QueryGroupChatMemberList",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "QueryGroupChatMemberList",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "QueryGroupChatMemberList",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.QueryGroupChatMemberList(info)
	return data, code, err
}
func (mw loggingMiddleware) GetArrayGroupInfo(info *GroupChatInfo) (data interface{}, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "GetArrayGroupInfo",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "GetArrayGroupInfo",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "GetArrayGroupInfo",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.GetArrayGroupInfo(info)
	return data, code, err
}
func (mw loggingMiddleware) GetMyGroupChat(info *GroupChatInfo) (data interface{}, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "GetMyGroupChat",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "GetMyGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "GetMyGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.GetMyGroupChat(info)
	return data, code, err
}
func (mw loggingMiddleware) UpdateGroupChat(info *GroupChatInfo) (data interface{}, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "UpdateGroupChat",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "UpdateGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "UpdateGroupChat",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.UpdateGroupChat(info)
	return data, code, err
}

//parkour=========================================================start
func (mw loggingMiddleware) AddOfficialMSG(chat *Chat) (code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "AddOfficialMSG",
		"input", chat,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "AddOfficialMSG",
				"input", chat,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "AddOfficialMSG",
				"input", chat,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.AddOfficialMSG(chat)
	return
}
func (mw loggingMiddleware) LookOfficialMSG() (data []map[string]string, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "LookOfficialMSG",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "LookOfficialMSG",
				"input", "",
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "LookOfficialMSG",
				"input", "",
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.LookOfficialMSG()
	return
}
func (mw loggingMiddleware) InformChat(chat *Chat) (code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "InformChat",
		"input", chat,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "InformChat",
				"input", chat,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "InformChat",
				"input", chat,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.InformChat(chat)
	return
}

//parkour======================================================end
func (mw loggingMiddleware) SearchChatInfo(info *GroupChatInfo) (data map[string]string, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "SearchChatInfo",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "SearchChatInfo",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "SearchChatInfo",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.SearchChatInfo(info)
	return data, code, err
}
func (mw loggingMiddleware) GetClassId(info *GroupChatInfo) (data string, code string, err error) {
	log := logger.GetDefaultLogger()
	log.Info(
		"method", "GetClassId",
		"input", info,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "GetClassId",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "GetClassId",
				"input", info,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.GetClassId(info)
	return data, code, err
}

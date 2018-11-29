package main

import (
	"time"
)

type loggingMiddleware struct {
	next McreateClassService
}

func (mw loggingMiddleware) SelectByScInviteCode(sy *SchoolYear) (schoolYear map[string]interface{}, err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "SelectByScInviteCode",
		"input", sy,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SelectByScInviteCode",
				"input", sy,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SelectByScInviteCode",
				"input", sy,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	schoolYear,err, code = mw.next.SelectByScInviteCode(sy)
	return
}
func (mw loggingMiddleware) CreateClass(sy *SchoolYear) (schoolYear map[string]string, err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "CreateClass",
		"input", sy,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "CreateClass",
				"input", sy,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "CreateClass",
				"input", sy,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	schoolYear, err, code = mw.next.CreateClass(sy)
	return
}
func (mw loggingMiddleware) TeacherJoinClass(sy *SchoolYear) (err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "TeacherJoinClass",
		"input", sy,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "TeacherJoinClass",
				"input", sy,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "TeacherJoinClass",
				"input", sy,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	 err, code = mw.next.TeacherJoinClass(sy)
	return
}

func (mw loggingMiddleware) FamilyJoinClass(sy *SchoolYear) (err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "FamilyJoinClass",
		"input", sy,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "FamilyJoinClass",
				"input", sy,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "FamilyJoinClass",
				"input", sy,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err, code = mw.next.FamilyJoinClass(sy)
	return
}
func (mw loggingMiddleware) NewMember(gp *SchoolYear) (data []map[string]string, err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "NewMember",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "NewMember",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "NewMember",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, err, code = mw.next.NewMember(gp)
	return
}
func (mw loggingMiddleware) ApproveMembers(gp *SchoolYear) (err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "ApproveMembers",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ApproveMembers",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ApproveMembers",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err, code = mw.next.ApproveMembers(gp)
	return
}
func (mw loggingMiddleware) ManagerMember(gp *SchoolYear) (data []map[string]string, err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "ManagerMember",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ManagerMember",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ManagerMember",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, err, code = mw.next.ManagerMember(gp)
	return
}
func (mw loggingMiddleware) OperateMember(gp *SchoolYear) (err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "OperateMember",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "OperateMember",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "OperateMember",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)

		}
	}(time.Now())

	err, code = mw.next.OperateMember(gp)
	return
}
func (mw loggingMiddleware) ClassQrCode(gp *SchoolYear) (data map[string]string,err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "ClassQrCode",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ClassQrCode",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ClassQrCode",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)

		}
	}(time.Now())

	data,err, code = mw.next.ClassQrCode(gp)
	return
}
func (mw loggingMiddleware) UpdateStudent(gp *SchoolYear) (err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "UpdateStudent",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "UpdateStudent",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "UpdateStudent",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err, code = mw.next.UpdateStudent(gp)
	return
}
func (mw loggingMiddleware) FindAllMember(gp *SchoolYear) (data []map[string]string,err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "FindAllMember",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "FindAllMember",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "FindAllMember",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data,err, code = mw.next.FindAllMember(gp)
	return
}
func (mw loggingMiddleware) FindTeacherMember(gp *SchoolYear) (data []map[string]string,err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "FindTeacherMember",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "FindTeacherMember",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "FindTeacherMember",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data,err, code = mw.next.FindTeacherMember(gp)
	return
}
func (mw loggingMiddleware) UpdateGroupChatInfo(gp *SchoolYear) (err error, code string) {
	log.Info(
		"check_log", "yes",
		"method_start", "UpdateGroupChatInfo",
		"input", gp,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "UpdateGroupChatInfo",
				"input", gp,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "UpdateGroupChatInfo",
				"input", gp,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err, code = mw.next.UpdateGroupChatInfo(gp)
	return
}
func (mw loggingMiddleware) UpdateTeachInfo(gp *SchoolYear) (err error, code string) {
		log.Info(
			"check_log", "yes",
			"method_start", "UpdateTeachInfo",
			"input", gp,
		)
		defer func(begin time.Time) {
			if err != nil {
				log.Error(
					"check_log", "yes",
					"method_end", "UpdateTeachInfo",
					"input", gp,
					"status", "fail",
					"msg", err,
					"took", time.Since(begin),
				)
			} else {
				log.Info(
					"check_log", "yes",
					"method_end", "UpdateTeachInfo",
					"input", gp,
					"status", "success",
					"took", time.Since(begin),
				)
			}
		}(time.Now())

		err, code = mw.next.UpdateTeachInfo(gp)
		return
	}
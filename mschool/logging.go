package main

import (
	"time"
)

type loggingMiddleware struct {
	next MschoolService
}

func (mw loggingMiddleware) CreateSchoolYear(school *School) (request map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "CreateSchoolYear",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "CreateSchoolYear",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "CreateSchoolYear",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	request, code, err = mw.next.CreateSchoolYear(school)
	return
}

func (mw loggingMiddleware) SearchSchool(school *School) (request map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SearchSchool",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SearchSchool",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SearchSchool",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	request, code, err = mw.next.SearchSchool(school)
	return
}

func (mw loggingMiddleware) MySchool(school *School) (request [][]map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "MySchool",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "MySchool",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "MySchool",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	request, code, err = mw.next.MySchool(school)
	return
}

func (mw loggingMiddleware) SetWorkDay(school *School) (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SetWorkDay",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SetWorkDay",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SetWorkDay",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SetWorkDay(school)
	return
}

func (mw loggingMiddleware) LookWorkDay(school *School) (request []map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "LookWorkDay",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "LookWorkDay",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "LookWorkDay",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	request, code, err = mw.next.LookWorkDay(school)
	return
}

func (mw loggingMiddleware) WorkDay(school *School) (request []map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "WorkDay",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "WorkDay",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "WorkDay",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	request, code, err = mw.next.WorkDay(school)
	return
}

func (mw loggingMiddleware) UpSchoolYear(school *School) (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "UpSchoolYear",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "UpSchoolYear",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "UpSchoolYear",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.UpSchoolYear(school)
	return
}

func (mw loggingMiddleware) FaceDataGather() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "FaceDataGather",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "FaceDataGather",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "FaceDataGather",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.FaceDataGather()
	return
}

func (mw loggingMiddleware) LabelGuidDataGather() (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "LabelGuidDataGather",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "LabelGuidDataGather",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "LabelGuidDataGather",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.LabelGuidDataGather()
	return
}

func (mw loggingMiddleware) GetDataFileUrl(school *School) (returnMap map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "LabelGuidDataGather",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "LabelGuidDataGather",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "LabelGuidDataGather",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	returnMap, code, err = mw.next.GetDataFileUrl(school)
	return
}
func (mw loggingMiddleware) DelDataFile(school *School) (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "LabelGuidDataGather",
		"input", school,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "LabelGuidDataGather",
				"input", school,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "LabelGuidDataGather",
				"input", school,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.DelDataFile(school)
	return
}

package main

import (
	"time"
)

type loggingMiddleware struct {
	next AutismService
}

func (mw loggingMiddleware) StarDetails(autism *Autism) (returnMap map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "StarDetails",
		"input", autism,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "StarDetails",
				"input", autism,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "StarDetails",
				"input", autism,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	returnMap, code, err = mw.next.StarDetails(autism)
	return
}

func (mw loggingMiddleware) SaveComment(autism *Autism) (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SaveComment",
		"input", autism,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SaveComment",
				"input", autism,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SaveComment",
				"input", autism,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.SaveComment(autism)
	return
}

//=======================================================================================================
func (mw loggingMiddleware) StarList(autism *Autism) (brightCount string, likeCount string, rolls []map[string]string, starList []map[string]string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "StarList",
		"input", autism,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "StarList",
				"input", autism,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "StarList",
				"input", autism,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())
	brightCount, likeCount, rolls, starList, err = mw.next.StarList(autism)
	return
}

func (mw loggingMiddleware) Likes(autism *Autism) (err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "Likes",
		"input", autism,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "Likes",
				"input", autism,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "Likes",
				"input", autism,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())
	err = mw.next.Likes(autism)
	return
}
func (mw loggingMiddleware) GetUnionid(autism *Autism) (data Un_ionid, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "GetUnionid",
		"input", autism,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "GetUnionid",
				"input", autism,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "GetUnionid",
				"input", autism,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())
	data, err = mw.next.GetUnionid(autism)
	return
}

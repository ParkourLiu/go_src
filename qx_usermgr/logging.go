package main

import (
	"fmt"
	"time"
	U "qx_user/caller"
)

type loggingMiddleware struct {
	next UserService
}

func (mw loggingMiddleware) RegAndLogin(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "RegAndLogin",
		"input", fmt.Sprint(ur),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "RegAndLogin",
				"input", fmt.Sprint(ur),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "RegAndLogin",
				"input", fmt.Sprint(ur),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.RegAndLogin(ur)
	return
}
func (mw loggingMiddleware) OtherLogin(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "OtherLogin",
		"input", fmt.Sprint(ur),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "OtherLogin",
				"input", fmt.Sprint(ur),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "OtherLogin",
				"input", fmt.Sprint(ur),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.OtherLogin(ur)
	return
}
func (mw loggingMiddleware) SearchUser(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SearchUser",
		"input", fmt.Sprint(ur),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SearchUser",
				"input", fmt.Sprint(ur),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SearchUser",
				"input", fmt.Sprint(ur),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.SearchUser(ur)
	return
}
func (mw loggingMiddleware) UpdateUser(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "UpdateUser",
		"input", fmt.Sprint(ur),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "UpdateUser",
				"input", fmt.Sprint(ur),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "UpdateUser",
				"input", fmt.Sprint(ur),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.UpdateUser(ur)
	return
}
func (mw loggingMiddleware) ChangeBind(ur *U.UserRequest) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "ChangeBind",
		"input", fmt.Sprint(ur),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ChangeBind",
				"input", fmt.Sprint(ur),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ChangeBind",
				"input", fmt.Sprint(ur),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.ChangeBind(ur)
	return
}

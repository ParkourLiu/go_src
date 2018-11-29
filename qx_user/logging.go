package main

import (
	"time"
)

type loggingMiddleware struct {
	next UserService
}

func (mw loggingMiddleware) SearchUserById(user *User) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SearchUserById",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SearchUserById",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SearchUserById",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.SearchUserById(user)
	return
}

func (mw loggingMiddleware) SearchUsers(user *User) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "SearchUsers",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "SearchUsers",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "SearchUsers",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.SearchUsers(user)
	return
}

func (mw loggingMiddleware) AddUser(user *User) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "AddUser",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "AddUser",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "AddUser",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.AddUser(user)
	return
}

func (mw loggingMiddleware) UpdateUser(user *User) (data map[string]interface{}, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "UpdateUser",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "UpdateUser",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "UpdateUser",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	data, code, err = mw.next.UpdateUser(user)
	return
}

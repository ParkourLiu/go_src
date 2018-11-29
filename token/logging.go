package main

import (
	"time"
)

type loggingMiddleware struct {
	next TokenService
}

func (mw loggingMiddleware) CreateToken(token *Token) (tk string, err error) {
	log.Info(
		"method", "createToken",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "createToken",
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "createToken",
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())
	tk, err = mw.next.CreateToken(token)
	return
}

func (mw loggingMiddleware) DeleteToken(token *Token) (err error) {
	log.Info(
		"method", "deleteToken",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method", "deleteToken",
				"err", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method", "deleteToken",
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	err = mw.next.DeleteToken(token)
	return
}

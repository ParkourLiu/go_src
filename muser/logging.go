package main

import (
	"fmt"
	"time"
)

type loggingMiddleware struct {
	next UserService
}

func (mw loggingMiddleware) MyHome(user *User) (userMap map[string]string, helpList []string, helpPetList []string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "MyHome",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "MyHome",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "MyHome",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	userMap, helpList, helpPetList, code, err = mw.next.MyHome(user)
	return
}

func (mw loggingMiddleware) HomePageSlogan() (userMap []map[string]string, code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "HomePageSlogan",
		"input", "",
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "HomePageSlogan",
				"input", "",
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "HomePageSlogan",
				"input", "",
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	userMap, code, err = mw.next.HomePageSlogan()
	return
}

func (mw loggingMiddleware) Reg(user *User) (code string, userId string, havePassword string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "Reg",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "Reg",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "Reg",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, userId, havePassword, err = mw.next.Reg(user)
	return
}

func (mw loggingMiddleware) Login(user *User) (code string, userId string, havePassword string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "Login",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "Login",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "Login",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, userId, havePassword, err = mw.next.Login(user)
	return
}

func (mw loggingMiddleware) ShortcutLogin(user *User) (code string, userId string, havePassword string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "ShortcutLogin",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ShortcutLogin",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ShortcutLogin",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, userId, havePassword, err = mw.next.ShortcutLogin(user)
	return
}

func (mw loggingMiddleware) FindPassword(user *User) (code string, userId string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "FindPassword",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "FindPassword",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "FindPassword",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, userId, err = mw.next.FindPassword(user)
	return
}

func (mw loggingMiddleware) ChangePhoneNo(user *User) (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "ChangePhoneNo",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ChangePhoneNo",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ChangePhoneNo",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.ChangePhoneNo(user)
	return
}

func (mw loggingMiddleware) OtherLogin(user *User) (code string, userId string, phoneNo string, havePassword string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "OtherLogin",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "OtherLogin",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "OtherLogin",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, userId, phoneNo, havePassword, err = mw.next.OtherLogin(user)
	return
}

func (mw loggingMiddleware) UpdateUser(user *User) (code string, err error) {
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

	code, err = mw.next.UpdateUser(user)
	return
}

func (mw loggingMiddleware) CheckPhoneBook(phoneBook *PhoneBook) (code string, err error, hashFlag string) {
	log.Info(
		"check_log", "yes",
		"method_start", "CheckPhoneBook",
		"input", fmt.Sprint(phoneBook),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "CheckPhoneBook",
				"input", fmt.Sprint(phoneBook),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "CheckPhoneBook",
				"input", fmt.Sprint(phoneBook),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err, hashFlag = mw.next.CheckPhoneBook(phoneBook)
	return code, err, hashFlag
}

func (mw loggingMiddleware) PhoneBookUser(phoneBook *PhoneBook) (code string, err error, haveuser []map[string]string, nouser []map[string]string) {
	log.Info(
		"check_log", "yes",
		"method_start", "PhoneBookUser",
		"input", fmt.Sprint(phoneBook),
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "PhoneBookUser",
				"input", fmt.Sprint(phoneBook),
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "PhoneBookUser",
				"input", fmt.Sprint(phoneBook),
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err, haveuser, nouser = mw.next.PhoneBookUser(phoneBook)
	return code, err, haveuser, nouser
}

func (mw loggingMiddleware) ActiveUser(user *User) (code string, err error) {
	log.Info(
		"check_log", "yes",
		"method_start", "ActiveUser",
		"input", user,
	)
	defer func(begin time.Time) {
		if err != nil {
			log.Error(
				"check_log", "yes",
				"method_end", "ActiveUser",
				"input", user,
				"status", "fail",
				"msg", err,
				"took", time.Since(begin),
			)
		} else {
			log.Info(
				"check_log", "yes",
				"method_end", "ActiveUser",
				"input", user,
				"status", "success",
				"took", time.Since(begin),
			)
		}
	}(time.Now())

	code, err = mw.next.ActiveUser(user)
	return
}

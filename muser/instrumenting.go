package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           UserService
}

func (mw instrumentingMiddleware) MyHome(user *User) (userMap map[string]string, helpList []string, helpPetList []string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "MyHome", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	userMap, helpList, helpPetList, code, err = mw.next.MyHome(user)
	return
}

func (mw instrumentingMiddleware) HomePageSlogan() (userMap []map[string]string, code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HomePageSlogan", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	userMap, code, err = mw.next.HomePageSlogan()
	return
}

func (mw instrumentingMiddleware) Reg(user *User) (code string, userId string, havePassword string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Reg", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, userId, havePassword, err = mw.next.Reg(user)
	return
}

func (mw instrumentingMiddleware) Login(user *User) (code string, userId string, havePassword string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, userId, havePassword, err = mw.next.Login(user)
	return
}

func (mw instrumentingMiddleware) ShortcutLogin(user *User) (code string, userId string, havePassword string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ReShortcutLoging", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, userId, havePassword, err = mw.next.ShortcutLogin(user)
	return
}

func (mw instrumentingMiddleware) FindPassword(user *User) (code string, userId string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindPassword", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, userId, err = mw.next.FindPassword(user)
	return
}

func (mw instrumentingMiddleware) ChangePhoneNo(user *User) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ChangePhoneNo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.ChangePhoneNo(user)
	return
}

func (mw instrumentingMiddleware) OtherLogin(user *User) (code string, userId string, phoneNo string, havePassword string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "OtherLogin", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, userId, phoneNo, havePassword, err = mw.next.OtherLogin(user)
	return
}

func (mw instrumentingMiddleware) UpdateUser(user *User) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.UpdateUser(user)
	return
}

func (mw instrumentingMiddleware) CheckPhoneBook(phoneBook *PhoneBook) (code string, err error, hashFlag string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CheckPhoneBook", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err, hashFlag = mw.next.CheckPhoneBook(phoneBook)
	return code, err, hashFlag
}

func (mw instrumentingMiddleware) PhoneBookUser(phoneBook *PhoneBook) (code string, err error, haveuser []map[string]string, nouser []map[string]string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PhoneBookUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err, haveuser, nouser = mw.next.PhoneBookUser(phoneBook)
	return code, err, haveuser, nouser
}

func (mw instrumentingMiddleware) ActiveUser(user *User) (code string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ActiveUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	code, err = mw.next.ActiveUser(user)
	return
}

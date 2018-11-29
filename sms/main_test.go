package main

import (
	"testing"
)

func TestSaveUser_1(t *testing.T) {
	var svc SmsService
	svc = smsService{}
	svc = loggingMiddleware{svc}
	a := []string{"aaa"}
	b := []string{"17671774535"}
	SmsInfo1 := &SmsInfo{
		Params: a,
		Mobile: b,
		Tpl_id: "57999",
	}
	_, err := svc.Sms(SmsInfo1)
	if err != nil {
		t.Error(err)
	}

	a = []string{"aaa", "bbb"}
	b = []string{"17671774535"}
	SmsInfo2 := &SmsInfo{
		Params: a,
		Mobile: b,
		Tpl_id: "57999",
	}
	_, err = svc.Sms(SmsInfo2)
	if err == nil {
		t.Error(err)
	}

	a = []string{"aaa"}
	b = []string{"17671774535"}
	SmsInfo3 := &SmsInfo{
		Params: a,
		Mobile: b,
		Tpl_id: "",
	}
	_, err = svc.Sms(SmsInfo3)
	if err == nil {
		t.Error(err)
	}

}

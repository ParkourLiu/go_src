package main

import (
	"testing"
)

func TestPushUser_1(t *testing.T) {
	var svc PushService
	svc = pushService{}
	svc = loggingMiddleware{svc}
	list := []pushInfo{}
	pushInfo0 := pushInfo{
		Title: "测试标题",
		Text:  "测试内容",
		Alias:   "6685c1b99792c68fcef07029f3ded9c2",
	}

	list = append(list, pushInfo0)
	pushInfoList0 := &pushInfoList{
		PushType:  "0",
		PushInfos: list,
	}
	_, err := svc.Push(pushInfoList0)
	if err != nil {
		t.Error(err)
	}

	pushInfoList1 := &pushInfoList{
		PushType:  "1",
		PushInfos: list,
	}
	_, err = svc.Push(pushInfoList1)
	if err != nil {
		t.Error(err)
	}

	pushInfoList2 := &pushInfoList{
		PushType:  "3",
		PushInfos: list,
	}
	_, err = svc.Push(pushInfoList2)
	if err == nil {
		t.Error(err)
	}

	pushInfoList3 := &pushInfoList{
		PushType: "1",
	}
	_, err = svc.Push(pushInfoList3)
	if err == nil {
		t.Error(err)
	}

	list1 := []pushInfo{}
	pushInfo1 := pushInfo{
		Title: "rerwer",
		Text:  "rewtwert",
		Alias:   "",
	}

	list1 = append(list1, pushInfo1)
	pushInfoList4 := &pushInfoList{
		PushType:  "1",
		PushInfos: list1,
	}
	_, err = svc.Push(pushInfoList4)
	if err == nil {
		t.Error(err)
	}

	list2 := []pushInfo{}
	pushInfo2 := pushInfo{
		Title: "rerwer",
		Text:  "rewtwert",
		Alias:   "eqewq",
	}

	list2 = append(list2, pushInfo2)
	pushInfoList5 := &pushInfoList{
		PushType:  "1",
		PushInfos: list2,
	}
	_, err = svc.Push(pushInfoList5)
	if err == nil {
		t.Error(err)
	}

	pushInfoList6 := &pushInfoList{
		PushType:  "0",
		PushInfos: list2,
	}
	_, err = svc.Push(pushInfoList6)
	if err == nil {
		t.Error(err)
	}
	//	srcMail := &User{
	//		UserId:   "title1",
	//		UserName: "text1",
	//	}
	//	m, err := json.Marshal(srcMail)
	//	if err != nil {
	//		t.Error(err)
	//	}

	//	//create producer
	//	s := prd.NewSender(log, "amqp://lio:zaq12wsx1@localhost:5672/", "ex_monitor")
	//	defer s.Close()
	//	err = s.Publish("mail", string(m))
	//	if err != nil {
	//		t.Error(err)
	//	}

	//	srcMail = &User{
	//		UserId:   "title2",
	//		UserName: "text2",
	//	}
	//	m, err = json.Marshal(srcMail)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	err = s.Publish("mail", string(m))
	//	if err != nil {
	//		t.Error(err)
	//	}

	//	srcMail = &User{
	//		UserId:   "title3",
	//		UserName: "text3",
	//	}
	//	m, err = json.Marshal(srcMail)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	err = s.Publish("mail", string(m))
	//	if err != nil {
	//		t.Error(err)
	//	}

	//	srcMail = &User{
	//		UserId:   "title4",
	//		UserName: "text4",
	//	}
	//	m, err = json.Marshal(srcMail)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	err = s.Publish("mail", string(m))
	//	if err != nil {
	//		t.Error(err)
	//	}

	//	srcMail = &User{
	//		UserId:   "title5",
	//		UserName: "text5",
	//	}
	//	m, err = json.Marshal(srcMail)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	err = s.Publish("mail", string(m))
	//	if err != nil {
	//		t.Error(err)
	//	}
}

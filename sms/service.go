package main

import (
	"fmt"
	monitor "monitor/client"
	"time"
)

// StringService provides operations on strings.
type SmsService interface {
	Sms(info *SmsInfo) (string, error)
}

type smsService struct{}

type SmsInfo struct {
	Params []string `json:"params"`
	Mobile []string `json:"mobile"`
	Tpl_id string   `json:"tpl_id"`
}

func (info *SmsInfo) String() string {
	params := info.Params
	mobile := info.Mobile
	tpl_id := info.Tpl_id
	return "params:" + fmt.Sprint(params) + ";mobile:" + fmt.Sprint(mobile) + ";tpl_id:" + tpl_id
}

func (service smsService) Sms(info *SmsInfo) (string, error) {
	log.Debug("method_start", "Sms", "input", info)
	smsReq, url := smsRequestInit(info)
	sms := NewMore_Template_qcloudSMS(url)
	par, err := sms.More_Template_qcloudSMS(smsReq)
	log.Debug("███par:", par, "err:", err)
	if err != nil {
		return "", err
	}
	if "OK" != par.Par["errmsg"] {
		monitorClient.PushErrorLog(&monitor.MonitorRequest{serviceName, "Sms", fmt.Sprint(par.Par["errmsg"]), "短信发送异常", fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))})
		return "调用异常", nil
	}
	redisClient.Incr(serviceName) //记录成功发送短信的次数
	log.Debug("method_end", "Sms", "status", "success")
	return "调用正常", nil
}

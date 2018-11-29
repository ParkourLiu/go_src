package main

import (
	"mtcomm/common"
	"mtcomm/db/mysql"
	"sync"
	"time"
	mail "mail/client"
	sms "sms/client"
	"github.com/bluele/gcache"
)

// StringService provides operations on strings.
type MonitorService interface {
	monitor(*monitorRequest)
}

var pMutex sync.Mutex
var eMutex sync.Mutex

type monitorService struct{}

func (service monitorService) monitor(m *monitorRequest) {
	log.Debug("method_start", "monitor", "input", m)

	if m.ServiceName != "sms" {
		go func() {
			ps, err := getPhoneNoes()
			if err == nil && len(ps) > 0 {
				// 10 分钟只能发送一次
				_, err :=redisClient.Get("monitor:sms_flag")
				if err != nil {
					//缓存不存在或者出错，发送通知
					p := &sms.SmsRequest{
						Params: []string{"mtalk"},
						Mobile: ps,
						Tpl_id: "104683",	//系统报错短信模板
					}
					err1 := smsClient.Sms(p)
					if err1 != nil {
						log.Error("send sms", "fail", "sms_err_msg", err1.Error())
					}else{
						log.Debug("send sms", "sucess")
					}
					//如果发送成功
					redisClient.SetStringAndExpire("monitor:sms_flag", "1", uint32(600))
				}
			}
		}()
	}
	
	if m.ServiceName != "mail" {
		go func() {
			es, err := getEmailAddress()
			if err == nil && len(es) > 0 {
				// 10 分钟只能发送一次
				_, err :=redisClient.Get("monitor:email_flag")
				if err != nil {
					//缓存不存在或者出错，发送通知
					p := &mail.MailRequest{
						To:              common.Slice2StringBySemi(es),
						SubjectTemplate: "1",
						Text:            err.Error(),
						BodyTemplate:    "1",
					}
					err1 := mailClient.PushEmail(p)
					if err1 != nil {
						log.Error("send mail", "fail", "mail_err_msg", err1.Error())
					}else{
						log.Debug("send mail", "sucess")
					}
					//如果发送成功
					redisClient.SetStringAndExpire("monitor:email_flag", "1", uint32(600))
				}
			}
		}()
	}
	
	// 错误信息保存到数据库
	go func() {
		saveMonitorData2DB(m)
	}()

	log.Debug("method_end", "monitor", "status", "success")
}

// 错误信息保存到数据库
func saveMonitorData2DB(m *monitorRequest) error {
	return mysqlClient.Execute(&mysql.Stmt{Sql: "insert into errlog(id, serviceName, methodName, param, errorMsg, errorTime) values(uuid(), ?, ?, ?, ?, ?)", Args: []interface{}{m.ServiceName, m.MethodName, m.Param, m.ErrorMsg, m.ErrorTime}})
}

//取得电话
func getPhoneNoes() ([]string, error) {
	ps, err := cacheClient.Get("monitor_phoneNoes")
	if err != nil && err != gcache.KeyNotFoundError {
		// 缓存出错
		return getPhoneNoesFromDB()
	} else if err == gcache.KeyNotFoundError {
		// not in cache
		psdb, err := getPhoneNoesFromDB()
		if err != nil {
			return nil, err
		}

		psstr := common.Slice2String(psdb)

		pMutex.Lock()
		_, err1 := cacheClient.Get("monitor_phoneNoes")
		if err1 == gcache.KeyNotFoundError {
			cacheClient.SetWithExpire("monitor_phoneNoes", psstr, time.Duration(10*time.Minute))
		}
		pMutex.Unlock()

		return psdb, nil
	} else {
		s := ps.(string)
		return common.String2Slice(s), nil
	}
}

//取得邮件地址
func getEmailAddress() ([]string, error) {
	ps, err := cacheClient.Get("monitor_email")
	if err != nil && err != gcache.KeyNotFoundError {
		// 缓存出错
		return getEmailAddressFromDB()
	} else if err == gcache.KeyNotFoundError {
		// not in cache
		psdb, err := getEmailAddressFromDB()
		if err != nil {
			return nil, err
		}

		psstr := common.Slice2String(psdb)

		eMutex.Lock()
		_, err1 := cacheClient.Get("monitor_email")
		if err1 == gcache.KeyNotFoundError {
			cacheClient.SetWithExpire("monitor_email", psstr, time.Duration(10*time.Minute))
		}
		eMutex.Unlock()

		return psdb, nil
	} else {
		s := ps.(string)
		return common.String2Slice(s), nil
	}
}

//从数据库取得运维人员的电话号码
func getPhoneNoesFromDB() ([]string, error) {
	log.Debug("msg", "get phone data from db")
	ms, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: "select distinct phoneNo from notifymembers where ad='A'", Args: []interface{}{}})
	if err != nil {
		return []string{}, nil
	}
	result := []string{}
	for _, m := range ms {
		result = append(result, m["phoneNo"])
	}
	return result, nil
}

//从数据库取得运维人员的邮件地址
func getEmailAddressFromDB() ([]string, error) {
	log.Debug("msg", "get email data from db")
	ms, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: "select distinct email from notifymembers where ad='A'", Args: []interface{}{}})
	if err != nil {
		return []string{}, nil
	}
	result := []string{}
	for _, m := range ms {
		result = append(result, m["email"])
	}
	return result, nil
}

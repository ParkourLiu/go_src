package main

import (
	"time"
	"testing"
	"mtcomm/db/mysql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	mail "mail/client"
	sms "sms/client"
)

type idMock struct {
	mock.Mock
}

func (m *idMock) Sms(req *sms.SmsRequest) error {
	args := m.Called()
	return args.Error(0)
}

func (m *idMock) PushEmail(info *mail.MailRequest) error {
	args := m.Called()
	return args.Error(0)

}

func TestMonitor(t *testing.T) {
	idm := new(idMock)
	idm.On("Sms").Return(nil)
	idm.On("PushEmail").Return(nil)
	mailClient = idm
	smsClient = idm
	
	//测试插入数据
	mysqlClient.Execute(&mysql.Stmt{Sql: "delete from errlog", Args: []interface{}{}})
	mysqlClient.Execute(&mysql.Stmt{Sql: "delete from notifymembers", Args: []interface{}{}})
	redisClient.Del("monitor:sms_flag", "monitor:email_flag")
	//after
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from notifymembers", Args: []interface{}{}})
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from errlog", Args: []interface{}{}})
	defer redisClient.Del("monitor:sms_flag", "monitor:email_flag")
	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into notifymembers values(1, 'Lio', '15900452449', 'lio.zhu@mm-world.com', 'A')", Args: []interface{}{}})
	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into notifymembers values(2, 'Parkour', '17671774535', 'parkour.liu@mm-world.com', 'A')", Args: []interface{}{}})

	/* create service */
	svc = monitorService{}
	svc = loggingMiddleware{svc}

	req := &monitorRequest{
		ServiceName: "lio-monitor",
		MethodName: "test",
		Param:	"1",
		ErrorMsg: "err",
		ErrorTime: "20180411",
	}
	svc.monitor(req)
	time.Sleep( 3 * time.Second)
	// comfirm
	u, err1 := mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from errlog where serviceName in ('lio-monitor')", Args: []interface{}{}})
	assert.NoError(t, err1)
	assert.Equal(t, "lio-monitor", u["serviceName"], "error")
	assert.Equal(t, "test", u["methodName"], "error")
	assert.Equal(t, "1", u["param"], "error")
	assert.Equal(t, "err", u["errorMsg"], "error")
	assert.Equal(t, "20180411", u["errorTime"], "error")
	
	sf, err2 := redisClient.Get("monitor:sms_flag")
	assert.NoError(t, err2)
	assert.Equal(t, "1", sf, "error")
	
	ef, err3 := redisClient.Get("monitor:email_flag")
	assert.NoError(t, err3)
	assert.Equal(t, "1", ef, "error")
}

func TestMonitor2(t *testing.T) {
	idm := new(idMock)
	idm.On("Sms").Return(nil)
	idm.On("PushEmail").Return(nil)
	mailClient = idm
	smsClient = idm
	
	//测试插入数据
	mysqlClient.Execute(&mysql.Stmt{Sql: "delete from errlog", Args: []interface{}{}})
	mysqlClient.Execute(&mysql.Stmt{Sql: "delete from notifymembers", Args: []interface{}{}})
	redisClient.Del("monitor:sms_flag", "monitor:email_flag")
	//after
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from notifymembers", Args: []interface{}{}})
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from errlog", Args: []interface{}{}})
	defer redisClient.Del("monitor:sms_flag", "monitor:email_flag")
	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into notifymembers values(1, 'Lio', '15900452449', 'lio.zhu@mm-world.com', 'A')", Args: []interface{}{}})
	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into notifymembers values(2, 'Parkour', '17671774535', 'parkour.liu@mm-world.com', 'A')", Args: []interface{}{}})

	/* create service */
	svc = monitorService{}
	svc = loggingMiddleware{svc}

	req := &monitorRequest{
		ServiceName: "lio-monitor",
		MethodName: "test",
		Param:	"1",
		ErrorMsg: "err",
		ErrorTime: "20180411",
	}
	svc.monitor(req)
	time.Sleep( 3 * time.Second)
	// comfirm
	u, err1 := mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from errlog where serviceName in ('lio-monitor')", Args: []interface{}{}})
	assert.NoError(t, err1)
	assert.Equal(t, "lio-monitor", u["serviceName"], "error")
	assert.Equal(t, "test", u["methodName"], "error")
	assert.Equal(t, "1", u["param"], "error")
	assert.Equal(t, "err", u["errorMsg"], "error")
	assert.Equal(t, "20180411", u["errorTime"], "error")
	
	sf, err2 := redisClient.Get("monitor:sms_flag")
	assert.NoError(t, err2)
	assert.Equal(t, "1", sf, "error")
	
	ef, err3 := redisClient.Get("monitor:email_flag")
	assert.NoError(t, err3)
	assert.Equal(t, "1", ef, "error")
	
	// monitor again
	req = &monitorRequest{
		ServiceName: "lio-monitor2",
		MethodName: "test",
		Param:	"1",
		ErrorMsg: "err",
		ErrorTime: "20180411",
	}
	svc.monitor(req)
	time.Sleep( 3 * time.Second)
	// comfirm
	u, err1 = mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from errlog where serviceName in ('lio-monitor2')", Args: []interface{}{}})
	assert.NoError(t, err1)
	assert.Equal(t, "lio-monitor2", u["serviceName"], "error")
	assert.Equal(t, "test", u["methodName"], "error")
	assert.Equal(t, "1", u["param"], "error")
	assert.Equal(t, "err", u["errorMsg"], "error")
	assert.Equal(t, "20180411", u["errorTime"], "error")
	
	sf, err2 = redisClient.Get("monitor:sms_flag")
	assert.NoError(t, err2)
	assert.Equal(t, "1", sf, "error")
	
	ef, err3 = redisClient.Get("monitor:email_flag")
	assert.NoError(t, err3)
	assert.Equal(t, "1", ef, "error")
}


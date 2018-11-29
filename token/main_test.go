package main

import (
	"testing"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func TestToken(t *testing.T) {
	//init
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: serviceName,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: serviceName,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	//测试创建token正常值，并验证结果
	/* create service */
	var svc TokenService
	svc = tokenService{}
	svc = loggingMiddleware{svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	token := &Token{PlatForm: "WeChat", UserId: "jewel_test", UUID: "123456789qwerty"}
	tk, err := svc.CreateToken(token)
	if err != nil {
		t.Error(err.Error())
		return
	}

	flag, err1 := redisClient.Exist("token:WeChat:123456789qwerty:" + tk)
	if err1 != nil {
		t.Error(err1.Error())
		return
	}
	if !flag {
		t.Error("测试创建token正常值 验证结果时 error")
	}
	flag, err1 = redisClient.Exist("token:jewel_test")
	if err1 != nil {
		t.Error(err1.Error())
		return
	}
	if !flag {
		t.Error("测试创建token正常值 验证结果时 error")
	}

	//测试创建token值异常
	token3 := &Token{UserId: "jewel_test", UUID: ""}
	_, err3 := svc.CreateToken(token3)
	if err3 == nil || err3.Error() != "Parameter Check Error" {
		t.Error("测试异常值 error")
		return
	}
	//测试删除token
	token4 := &Token{PlatForm: "WeChat", UserId: "jewel_test", Token: tk, UUID: "123456789qwerty"}
	err4 := svc.DeleteToken(token4)
	if err4 != nil {
		t.Error(err4.Error())
		return
	}

	flag2, err5 := redisClient.Exist("token:WeChat:123456789qwerty:" + tk)
	if err5 != nil {
		t.Error(err5.Error())
		return
	}
	if flag2 {
		t.Error("测试删除token正常值 验证结果时 error")
	}
	flag2, err5 = redisClient.Exist("token:jewel_test")
	if err5 != nil {
		t.Error(err5.Error())
		return
	}
	if flag2 {
		t.Error("测试删除token正常值 验证结果时 error")
	}

	//测试删除token值异常
	token6 := &Token{Token: "", UUID: "shqx"}
	err6 := svc.DeleteToken(token6)
	if err6 == nil || err6.Error() != "Parameter Check Error" {
		t.Error("测试删除token值异常 error")
		return
	}
}

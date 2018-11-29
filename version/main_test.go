package main

import (
	"fmt"
	"mtcomm/db/mysql"
	"testing"

	//	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	//	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

type idMock struct {
	mock.Mock
}

func (m *idMock) GetUniqueId() string {
	args := m.Called()
	return args.String(0)

}

func TestUser_1(t *testing.T) {
	idm := new(idMock)
	idm.On("GetUniqueId").Return("ttt")
	idGenClient = idm

	//测试插入数据
	mysqlClient.Execute(&mysql.Stmt{Sql: "delete from user where userId in ('Uttt')", Args: []interface{}{}})
	//after
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from user where userId in ('Uttt')", Args: []interface{}{}})

	//测试正常值，并验证结果
	/* create service */
	var svc UserService
	svc = userService{}

	user := &User{PhoneNo: "17671774535", Password: "17671774535"}
	userId, err := svc.AddUser(user)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("userId:" + userId)
	u, err1 := mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select userId, phoneNo from user where userId in ('Uttt')", Args: []interface{}{}})
	if err1 != nil {
		t.Error(err1.Error())
		return
	}
	if u["userId"] != "Uttt" || u["phoneNo"] != "17671774535" {
		fmt.Println(u["userId"])
		fmt.Println(u["phoneNo"])
		t.Error("测试正常值，并验证结果 error!")
	}

	//测试插入错误数据
	user = &User{PhoneNo: "17671774535"}
	_, err = svc.AddUser(user)
	if err == nil {
		t.Error(err.Error())
		return
	}

	//测试查询数据
	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into `user`(userId,mtalkNo,createTime,updateTime) values('qqq','111111111',now(),now())", Args: []interface{}{}})
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from user where userId in ('qqq')", Args: []interface{}{}})
	user = &User{UserId: "qqq"}
	a := map[string]string{}
	a, err = svc.SearchUser(user)
	fmt.Println(a)
	if err != nil {
		t.Error(err.Error())
		return
	}

	//测试错误数据
	user = &User{}
	a, err = svc.SearchUser(user)
	fmt.Println("█●█●█●测试：", a)
	if err == nil {
		t.Error(err.Error())
		return
	}

	//测试查询批量
	user = &User{}
	b, err1 := svc.SearchUsers(user)
	fmt.Println("█●█●█●测试：", b)
	if err1 != nil {
		t.Error(err1.Error())
		return
	}

	//测试修改数据
	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into `user`(userId,mtalkNo,createTime,updateTime) values('aaa','111111111',now(),now())", Args: []interface{}{}})
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from user where userId in ('aaa')", Args: []interface{}{}})
	user = &User{UserId: "aaa", Email: "123456789@qq.com"}
	err = svc.UpdateUser(user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	user = &User{PhoneNo: "123456789"}
	err = svc.UpdateUser(user)
	if err == nil {
		t.Error(err.Error())
		return
	}

	mysqlClient.Execute(&mysql.Stmt{Sql: "insert into `user`(userId,mtalkNo,createTime,updateTime) values('deleteTest','111111111',now(),now())", Args: []interface{}{}})
	defer mysqlClient.Execute(&mysql.Stmt{Sql: "delete from user where userId in ('deleteTest')", Args: []interface{}{}})
	user = &User{UserId: "deleteTest", PhoneNo: "123456789"}
	err = svc.DeleteUser(user)
	if err != nil {
		t.Error(err.Error())
		return
	}
	user = &User{PhoneNo: "123456789"}
	err = svc.DeleteUser(user)
	if err == nil {
		t.Error(err.Error())
		return
	}

}

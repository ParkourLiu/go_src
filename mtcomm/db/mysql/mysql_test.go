package mysql_test

import (
	"fmt"
	"mtcomm/db/mysql"
	"testing"

	logger "mtcomm/log"
)

func newMysqlInfo() *mysql.MysqlInfo {
	logger.SetDefaultLogLevel(logger.LevelDebug)
	return &mysql.MysqlInfo{
		UserName:     "root",
		Password:     "zaq12wsx1",
		IP:           "127.0.0.1",
		Port:         "3306",
		DatabaseName: "mm",
		Logger:       logger.GetDefaultLogger(),
	}
}

func TestUpdate(t *testing.T) {
	db := mysql.NewMysqlClient(newMysqlInfo())
	defer db.Close()

	//delete
	err := db.Execute(&mysql.Stmt{Sql: "delete from users where userId = '888'", Args: []interface{}{}})
	if err != nil {
		t.Error(err.Error())
	}

	err = db.Execute(&mysql.Stmt{Sql: "insert into users(userId, userName, v1, v2, v3) values(?, ?, ?, ?, now())", Args: []interface{}{"888", "444", "abc", 100}})
	if err != nil {
		t.Error(err.Error())
	}

	m := make(map[string]interface{})
	m["userName"] = "yoyo"
	m["v2"] = 200
	err = db.Update("users", "where userId='888'", m, false)
	if err != nil {
		t.Error(err.Error())
	}

	result, err4 := db.SearchOneRow(&mysql.Stmt{Sql: "select * from users where userId='888'", Args: []interface{}{}})
	if err4 != nil {
		t.Error(err4.Error())
	}
	if result["userId"] != "888" || result["userName"] != "yoyo" || result["v2"] != "200" || result["v1"] != "abc" {
		t.Error("Value error")
	}

	m["v2"] = 300
	err = db.Update("users", "where userId='888'", m, true)
	if err != nil {
		t.Error(err.Error())
	}

	//delete
	err3 := db.Execute(&mysql.Stmt{Sql: "delete from users where userId = '888'", Args: []interface{}{}})
	if err3 != nil {
		t.Error(err3.Error())
	}
}

func TestExecute(t *testing.T) {
	db := mysql.NewMysqlClient(newMysqlInfo())
	defer db.Close()
	//insert
	err := db.Execute(&mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"111", "222"}})
	if err != nil {
		t.Error(err.Error())
	}
	//count 1
	c1, err4 := db.Count(&mysql.Stmt{Sql: "select count(1) from users where userId='111'", Args: []interface{}{}})
	if err4 != nil {
		t.Error(err4.Error())
	}
	if c1 != 1 {
		t.Error("count error")
	}
	//count 2
	c1, err5 := db.Count(&mysql.Stmt{Sql: "select count(1) from users where userId=?", Args: []interface{}{"aaa"}})
	if err5 != nil {
		t.Error(err5.Error())
	}
	if c1 != 0 {
		t.Error("count error")
	}
	//update
	err2 := db.Execute(&mysql.Stmt{Sql: "update users set userName = ? where userId = '111'", Args: []interface{}{"333"}})
	if err2 != nil {
		t.Error(err2.Error())
	}
	//delete
	err3 := db.Execute(&mysql.Stmt{Sql: "delete from users where userId = '111'", Args: []interface{}{}})
	if err3 != nil {
		t.Error(err3.Error())
	}

}

func TestExecTransaction_1(t *testing.T) {
	db := mysql.NewMysqlClient(newMysqlInfo())
	defer db.Close()

	stmts := []mysql.Stmt{
		mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"111", "222"}},
		mysql.Stmt{Sql: "update users set userName = ? where userId = '111'", Args: []interface{}{"333"}},
		mysql.Stmt{Sql: "delete from users where userId = '111'", Args: []interface{}{}},
	}

	err := db.ExecTransaction(stmts)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestExecTransaction_2(t *testing.T) {
	db := mysql.NewMysqlClient(newMysqlInfo())
	defer db.Close()

	stmts := []mysql.Stmt{
		mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"aaa", "222"}},
		mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"bbb", "222"}},
		mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"bbb", "333"}},
	}

	err := db.ExecTransaction(stmts)
	if err == nil {
		t.Error()
	}
	//count 1
	c1, err4 := db.Count(&mysql.Stmt{Sql: "select count(1) from users where userId='aaa'", Args: []interface{}{}})
	if err4 != nil {
		t.Error(err4.Error())
	}
	if c1 != 0 {
		t.Error("error")
	}
}

func TestOneRow(t *testing.T) {
	db := mysql.NewMysqlClient(newMysqlInfo())
	defer db.Close()
	//insert
	err := db.Execute(&mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"111", "222"}})
	if err != nil {
		t.Error(err.Error())
	}
	err = db.Execute(&mysql.Stmt{Sql: "insert into users(userId, userName) values(?, ?)", Args: []interface{}{"666", "222"}})
	if err != nil {
		t.Error(err.Error())
	}
	err = db.Execute(&mysql.Stmt{Sql: "insert into users(userId, userName, v1, v2, v3) values(?, ?, ?, ?, now())", Args: []interface{}{"888", "444", "abc", 100}})
	if err != nil {
		t.Error(err.Error())
	}

	//search 1
	result, err4 := db.SearchOneRow(&mysql.Stmt{Sql: "select * from users where userId='111'", Args: []interface{}{}})
	if err4 != nil {
		t.Error(err4.Error())
	}
	if result["userId"] != "111" || result["userName"] != "222" || result["v1"] != "" || result["v2"] != "" || result["v3"] != "" {
		t.Error("Value error")
	}

	//search 2
	result, err4 = db.SearchOneRow(&mysql.Stmt{Sql: "select * from users where userId='aaa'", Args: []interface{}{}})
	if err4 != nil {
		t.Error(err4.Error())
	}
	if result != nil {
		t.Error("Value error")
	}

	//search 3
	result, err4 = db.SearchOneRow(&mysql.Stmt{Sql: "select * from users where userName='222'", Args: []interface{}{}})
	if err4.Error() != "Not only one row" {
		t.Error(err4.Error())
	}

	//search 4
	result, err4 = db.SearchOneRow(&mysql.Stmt{Sql: "select * from users where userId='888'", Args: []interface{}{}})
	if err4 != nil {
		t.Error(err4.Error())
	}
	if result["userId"] != "888" || result["userName"] != "444" || result["v2"] != "100" || result["v1"] != "abc" {
		t.Error("Value error")
	}
	fmt.Println("=============", result["v3"])

	///////////////SearchMutiRows//////////////
	//search 1
	result5, err5 := db.SearchMutiRows(&mysql.Stmt{Sql: "select * from users where userId='111'", Args: []interface{}{}})
	if err5 != nil {
		t.Error(err5.Error())
	}
	tmpMap := result5[0]
	if tmpMap["userId"] != "111" || tmpMap["userName"] != "222" || tmpMap["v1"] != "" || tmpMap["v2"] != "" || tmpMap["v3"] != "" {
		t.Error("Value error")
	}
	if len(result5) != 1 {
		t.Error("result count error")
	}

	//search 2
	result5, err5 = db.SearchMutiRows(&mysql.Stmt{Sql: "select * from users where userId='aaa'", Args: []interface{}{}})
	if err5 != nil {
		t.Error(err5.Error())
	}
	if len(result5) != 0 {
		t.Error("Value error")
	}
	if len(result5) != 0 {
		t.Error("result count error")
	}

	//search 3
	result5, err5 = db.SearchMutiRows(&mysql.Stmt{Sql: "select * from users where userName='222'", Args: []interface{}{}})
	if err5 != nil {
		t.Error(err5.Error())
	}
	tmpMap2 := result5[0]
	tmpMap3 := result5[1]
	if tmpMap2["userId"] != "111" || tmpMap2["userName"] != "222" || tmpMap2["v1"] != "" || tmpMap2["v2"] != "" || tmpMap2["v3"] != "" || tmpMap3["userId"] != "666" || tmpMap3["userName"] != "222" || tmpMap3["v1"] != "" || tmpMap3["v2"] != "" || tmpMap3["v3"] != "" {
		t.Error("Value error")
	}
	if len(result5) != 2 {
		t.Error("result count error")
	}

	//search 4
	result5, err5 = db.SearchMutiRows(&mysql.Stmt{Sql: "select * from users where userId='888'", Args: []interface{}{}})
	if err5 != nil {
		t.Error(err5.Error())
	}
	tmpMap4 := result5[0]
	if tmpMap4["userId"] != "888" || tmpMap4["userName"] != "444" || tmpMap4["v2"] != "100" || tmpMap4["v1"] != "abc" {
		t.Error("Value error")
	}
	if len(result5) != 1 {
		t.Error("result count error")
	}
	fmt.Println("=============", tmpMap4["v3"])

	//delete
	err3 := db.Execute(&mysql.Stmt{Sql: "delete from users where userId in ( '111', '666', '888')", Args: []interface{}{}})
	if err3 != nil {
		t.Error(err3.Error())
	}
}

package mysql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"mtcomm/common"
	logger "mtcomm/log"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient interface {
	Close()
	Execute(*Stmt) error
	ExecuteTx(*sql.Tx, *Stmt) error
	GetTransaction() (*sql.Tx, error)
	ExecTransaction(stmts []Stmt) error
	Count(*Stmt) (int, error)
	SearchOneRow(stmt *Stmt) (map[string]string, error)
	SearchMutiRows(stmt *Stmt) ([]map[string]string, error)
	//updateTime 表里是否包含updateTime字段
	Update(tableName string, whereSql string, data map[string]interface{}, updateTime bool) error
}

type Stmt struct {
	Sql  string
	Args []interface{}
}

func (s *Stmt) String() string {
	var b bytes.Buffer
	b.WriteString("sql: ")
	b.WriteString(s.Sql)
	b.WriteString(" Args: ")
	for _, arg := range s.Args {
		b.WriteString(fmt.Sprint(arg))
		b.WriteString(" ")
	}
	return b.String()
}

type mysqlClient struct {
	logger *logger.Logger
	db     *sql.DB
}

type MysqlInfo struct {
	UserName     string
	Password     string
	IP           string
	Port         string
	DatabaseName string
	Logger       *logger.Logger
	MaxIdleConns int //option
}

func getRows(rows *sql.Rows) ([]map[string]string, error) {
	results := make([]map[string]string, 0) //result
	if rows == nil {
		return nil, errors.New("rows is nil")
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var rawResult [][]byte
	var result map[string]string
	var dest []interface{}
	for rows.Next() {
		rawResult = make([][]byte, len(cols))
		result = make(map[string]string, len(cols))
		dest = make([]interface{}, len(cols))
		for i, _ := range rawResult {
			dest[i] = &rawResult[i]
		}

		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		for i, raw := range rawResult {
			if raw == nil {
				result[cols[i]] = ""
			} else {
				result[cols[i]] = string(raw)
			}
		}
		results = append(results, result)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

//返回值：nil代表没有数据。数组元素为空字符串代表null
func (c *mysqlClient) SearchOneRow(stmt *Stmt) (map[string]string, error) {
	c.logger.Debug("method_start", "SearchOneRow", "input", stmt)
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		c.logger.Debug("method_end", "SearchOneRow", "status", "fail", "msg", err.Error())
		return nil, err
	}
	defer stmtIns.Close()

	rows, err3 := stmtIns.Query(stmt.Args...)
	if err3 != nil {
		c.logger.Debug("method_end", "SearchOneRow", "status", "fail", "msg", err3.Error())
		return nil, err3
	}
	defer rows.Close()

	results, err2 := getRows(rows)
	if err2 != nil {
		c.logger.Debug("method_end", "SearchOneRow", "status", "fail", "msg", err2.Error())
		return nil, err2
	}
	if len(results) > 1 {
		msg := "Not only one row"
		c.logger.Debug("method_end", "SearchOneRow", "status", "fail", "msg", msg)
		return nil, errors.New(msg)
	} else if len(results) == 1 {
		c.logger.Debug("method_end", "SearchOneRow", "status", "success", "return", results[0])
		return results[0], nil
	} else {
		c.logger.Debug("method_end", "SearchOneRow", "status", "success", "return", "nil")
		return nil, nil
	}
}

//返回值：长度为0代表没有数据。数组元素为空字符串代表null
func (c *mysqlClient) SearchMutiRows(stmt *Stmt) ([]map[string]string, error) {
	c.logger.Debug("method_start", "SearchMutiRows", "input", stmt)
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		c.logger.Debug("method_end", "SearchMutiRows", "status", "fail", "msg", err.Error())
		return nil, err
	}
	defer stmtIns.Close()

	rows, err3 := stmtIns.Query(stmt.Args...)
	if err3 != nil {
		c.logger.Debug("method_end", "SearchMutiRows", "status", "fail", "msg", err3.Error())
		return nil, err3
	}
	defer rows.Close()

	results, err2 := getRows(rows)
	if err2 != nil {
		c.logger.Debug("method_end", "SearchMutiRows", "status", "fail", "msg", err2.Error())
		return nil, err2
	}
	c.logger.Debug("method_end", "SearchMutiRows", "status", "success", "return", results)
	return results, nil
}

//func NewMysqlClient(userName string, password string, ip string, port string, database string) *mysqlClient {
func NewMysqlClient(info *MysqlInfo) MysqlClient {
	//uri: "root:zaq12wsx1@tcp(localhost:3306)/mm?charset=utf8"
	if info.MaxIdleConns == 0 {
		info.MaxIdleConns = 30 //default
	}
	uri := info.UserName + ":" + info.Password + "@tcp(" + info.IP + ":" + info.Port + ")/" + info.DatabaseName + "?charset=utf8mb4"
	db, err := sql.Open("mysql", uri)
	if err != nil {
		info.Logger.Error("method", "NewMysqlClient", "status", "fail", "msg", err.Error())
		panic(err.Error())
	}
	db.Exec("SET NAMES 'utf8mb4'; SET CHARACTER SET utf8mb4;")
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		info.Logger.Error("method", "NewMysqlClient", "status", "fail", "msg", err.Error())
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	db.SetMaxIdleConns(30)

	return &mysqlClient{
		logger: info.Logger,
		db:     db,
	}
}

func (c *mysqlClient) Close() {
	c.db.Close()
}

func (c *mysqlClient) Count(stmt *Stmt) (int, error) {
	c.logger.Debug("method_start", "Count", "input", stmt)
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		c.logger.Debug("method_end", "Count", "status", "fail", "msg", err.Error())
		return 0, err
	}
	defer stmtIns.Close()

	var count int
	err = stmtIns.QueryRow(stmt.Args...).Scan(&count)
	if err != nil {
		c.logger.Debug("method_end", "Count", "status", "fail", "msg", err.Error())
		return 0, err
	}
	c.logger.Debug("method_end", "Count", "status", "success", "return", count)
	return count, nil
}

func (c *mysqlClient) Execute(stmt *Stmt) error {
	//c.logger.Debug("method_start", "Execute", "input", stmt)
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		c.logger.Debug("method_end", "Execute", "status", "fail", "msg", err.Error())
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(stmt.Args...)
	if err != nil {
		c.logger.Debug("method_end", "Execute", "status", "fail", "msg", err.Error())
		return err
	}
	//c.logger.Debug("method_end", "Execute", "status", "success")
	return err
}

//开启一个事务
func (c *mysqlClient) GetTransaction() (*sql.Tx, error) {
	return c.db.Begin()
}

//添加事物sql
func (c *mysqlClient) ExecuteTx(tx *sql.Tx, stmt *Stmt) error {
	c.logger.Debug("method_start", "ExecuteTx", "input", fmt.Sprint(stmt))
	_, err := tx.Exec(stmt.Sql, stmt.Args...)
	c.logger.Debug("method_end", "ExecuteTx", "status", "success")
	return err
}

func (c *mysqlClient) ExecTransaction(stmts []Stmt) error {
	c.logger.Debug("method_start", "ExecTransaction", "input", common.Slice2String(stmts))
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Debug("method_end", "ExecTransaction", "status", "fail", "msg", err.Error())
		return err
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		_, err = tx.Exec(stmt.Sql, stmt.Args...)
		if err != nil {
			c.logger.Debug("method_end", "ExecTransaction", "status", "fail", "msg", err.Error())
			return err
		}
	}
	//commit
	tx.Commit()
	c.logger.Debug("method_end", "ExecTransaction", "status", "success")
	return nil
}

func createSql(tableName string, whereSql string, data map[string]interface{}, updateTime bool) (string, []interface{}) {
	var buf bytes.Buffer
	buf.WriteString("update ")
	buf.WriteString(tableName)
	buf.WriteString(" set ")
	args := []interface{}{}
	for k, v := range data {
		buf.WriteString(k)
		buf.WriteString("= ?, ")
		args = append(args, v)
	}
	if updateTime {
		buf.WriteString("updateTime = now(), ")
	}
	suffix := buf.String()
	suffix = suffix[:len(suffix)-2]

	return suffix + " " + whereSql, args
}

func (c *mysqlClient) Update(tableName string, whereSql string, data map[string]interface{}, updateTime bool) error {
	c.logger.Debug("method_start", "Update", "input", tableName, "input", whereSql, "input", data)
	if len(data) <= 0 {
		return nil
		c.logger.Debug("method_end", "Update", "status", "success")
	}

	sql, args := createSql(tableName, whereSql, data, updateTime)
	stmt := &Stmt{
		Sql:  sql,
		Args: args,
	}
	c.logger.Debug("input", stmt)

	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		c.logger.Debug("method_end", "Update", "status", "fail", "msg", err.Error())
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(stmt.Args...)
	if err != nil {
		c.logger.Debug("method_end", "Update", "status", "fail", "msg", err.Error())
		return err
	}
	c.logger.Debug("method_end", "Update", "status", "success")
	return err
}

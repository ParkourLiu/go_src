package main

import (
	"mtcomm/db/mysql"
)

//查询数据库版本
func SVersion(client mysql.MysqlClient, version *Version) (map[string]string, error) {
	sql := "select * from version where 1=1 and flag = '" + version.Flag + "' "
	if version.IsForce != "" {
		sql += " and isForce='" + version.IsForce + "'"
	}
	sql += " order by createTime desc limit 0,1;"
	versionMap, err := client.SearchOneRow(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	return versionMap, err
}

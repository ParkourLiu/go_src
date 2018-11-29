package main

import (
	"mtcomm/db/mysql"
	"bytes"
)

//根据学校id查询出此校最新学年id
func s_syIdByScIdOneData(scId string) (map[string]string, error) {
	sql := "SELECT syId FROM `schoolyear` WHERE `scId`=? ORDER BY `createTime` DESC LIMIT 1"
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql: sql, Args: []interface{}{scId}})
}

//根据学年Id查询班级Id
func s_clIdBySyId(syId string) ([]map[string]string, error) {
	sql := "SELECT clId FROM `class` WHERE `syId`=?"
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{syId}})
}

//根据班级Id批量查询学生id
func s_stIdByclIds(clIdList []map[string]string) ([]map[string]string, error) {
	sbf := bytes.Buffer{}
	sbf.WriteString("select stId from student_class where ad='a' and `clId` in (")
	for i, v := range clIdList {
		sbf.WriteString("'")
		sbf.WriteString(v["clId"])
		sbf.WriteString("'")
		if i < len(clIdList)-1 {
			sbf.WriteString(",")
		}
	}
	sbf.WriteString(");")
	sql := sbf.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//查询一批学生中人脸被更新过的数据
func s_studentByFileAndstIds(clIdList []map[string]string) ([]map[string]string, error) {
	sbf := bytes.Buffer{}
	sbf.WriteString("select * from `student` where `facePath`!=`schoolFacePath` and `stId`in(")
	for i, v := range clIdList {
		sbf.WriteString("'")
		sbf.WriteString(v["stId"])
		sbf.WriteString("'")
		if i < len(clIdList)-1 {
			sbf.WriteString(",")
		}
	}
	sbf.WriteString(");")
	sql := sbf.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//修改一批学生中人脸被更新过的数据
func u_studentByFileAndstIds(studentList []map[string]string) error {
	sbf := bytes.Buffer{}
	sbf.WriteString("UPDATE student SET  schoolFacePath = facePath WHERE `facePath`!=`schoolFacePath` AND `stId`IN(")
	for i, v := range studentList {
		sbf.WriteString("'")
		sbf.WriteString(v["stId"])
		sbf.WriteString("'")
		if i < len(studentList)-1 {
			sbf.WriteString(",")
		}
	}
	sbf.WriteString(");")
	sql := sbf.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//查询一批学生中的Guid没同步过的数据 按stid排序
func s_guIdByStIds(studentList []map[string]string) ([]map[string]string, error) {
	sbf := bytes.Buffer{}
	sbf.WriteString("SELECT * FROM `student_guid` WHERE `synStatus`='0' AND `stId` IN(")
	for i, v := range studentList {
		sbf.WriteString("'")
		sbf.WriteString(v["stId"])
		sbf.WriteString("'")
		if i < len(studentList)-1 {
			sbf.WriteString(",")
		}
	}
	sbf.WriteString(") ORDER BY stId;")
	sql := sbf.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//修改一批学生中的Guid没同步过的数据为已同步
func u_guIdByStIds(studentList []map[string]string) error {
	sbf := bytes.Buffer{}
	sbf.WriteString("update `student_guid` set `synStatus`='1'  WHERE `synStatus`='0' AND `stId` IN(")
	for i, v := range studentList {
		sbf.WriteString("'")
		sbf.WriteString(v["stId"])
		sbf.WriteString("'")
		if i < len(studentList)-1 {
			sbf.WriteString(",")
		}
	}
	sbf.WriteString(");")
	sql := sbf.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

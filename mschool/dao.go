package main

import (
	"mtcomm/db/mysql"
	"bytes"
	"errors"
)

//动态查询学校
func s_school(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `school` where 1=1")
	if school.PrincipalUserId != "" {
		sqlBuffer.WriteString(" and principalUserId='")
		sqlBuffer.WriteString(school.PrincipalUserId)
		sqlBuffer.WriteString("'")
	}
	if school.ScId != "" {
		sqlBuffer.WriteString(" and scId='")
		sqlBuffer.WriteString(school.ScId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询学年
func s_schoolyear(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `schoolyear` where 1=1")
	if school.SyId != "" {
		sqlBuffer.WriteString(" and syId='")
		sqlBuffer.WriteString(school.SyId)
		sqlBuffer.WriteString("'")
	}
	if school.ScId != "" {
		sqlBuffer.WriteString(" and scId='")
		sqlBuffer.WriteString(school.ScId)
		sqlBuffer.WriteString("'")
	}
	if school.StudyYear != "" {
		sqlBuffer.WriteString(" and studyYear='")
		sqlBuffer.WriteString(school.StudyYear)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询我的校园
func s_myschool(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `myschool` where 1=1")
	if school.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(school.UserId)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString("ORDER BY `createTime` DESC")
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询workTime
func s_worktime(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `worktime` where 1=1")
	if school.SyId != "" {
		sqlBuffer.WriteString(" and syId='")
		sqlBuffer.WriteString(school.SyId)
		sqlBuffer.WriteString("'")
	}
	if school.GrId != "" {
		sqlBuffer.WriteString(" and grId='")
		sqlBuffer.WriteString(school.GrId)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString("ORDER BY `createTime` DESC")
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询workDay
func s_workday(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `workday` where 1=1")
	if school.SyId != "" {
		sqlBuffer.WriteString(" and syId='")
		sqlBuffer.WriteString(school.SyId)
		sqlBuffer.WriteString("'")
	}
	if school.GrId != "" {
		sqlBuffer.WriteString(" and grId='")
		sqlBuffer.WriteString(school.GrId)
		sqlBuffer.WriteString("'")
	}
	if school.ExceptionDay != "" {
		sqlBuffer.WriteString(" and exceptionDay='")
		sqlBuffer.WriteString(school.ExceptionDay)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString("ORDER BY `exceptionDay`")
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询holiday
func s_holiday(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `holiday` where 1=1")
	if school.Year != "" {
		sqlBuffer.WriteString(" and year='")
		sqlBuffer.WriteString(school.Year)
		sqlBuffer.WriteString("'")
	}
	if school.ExceptionDay != "" {
		sqlBuffer.WriteString(" and hoDay='")
		sqlBuffer.WriteString(school.ExceptionDay)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" ORDER BY `createTime` DESC")
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询grade
func s_grade(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `grade` where 1=1")
	if school.SyId != "" {
		sqlBuffer.WriteString(" and syId='")
		sqlBuffer.WriteString(school.SyId)
		sqlBuffer.WriteString("'")
	}
	if school.GrId != "" {
		sqlBuffer.WriteString(" and grId='")
		sqlBuffer.WriteString(school.GrId)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString("ORDER BY `grId` ASC") //顺序 id也是递增的，所以根据id排序
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//动态查询class
func s_class(school *School) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `class` where 1=1")
	if school.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(school.ClId)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString("ORDER BY `createTime` DESC")
	sql := sqlBuffer.String()
	result, err := mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return result, err
	}
	return result, nil
}

//插入school
func i_school(school *School) error {
	sql := "insert into `school`(`scId`,`scName`,`nature`,`address`,`type`,`location`,`scInviteCode`,`scInviteQRCode`,`principal`,`principalPhone`,`principalUserId`,`createTime`,`updateTime`) values(?,?,?,?,?,?,?,?,?,?,?,now(),now());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{school.ScId, school.ScName, school.Nature, school.Address, school.Type, school.Location, school.ScInviteCode, school.ScInviteQRCode, school.Principal, school.PrincipalPhone, school.PrincipalUserId}})
}

//插入schoolyear
func i_schoolyear(school *School) error {
	sql := "INSERT INTO `schoolyear`(`syId`,`scId`,`beginTime`,`endTime`,studyYear,`createTime`,`updateTime`) VALUES(?,?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{school.SyId, school.ScId, school.BeginTime, school.EndTime, school.StudyYear}})
}

//插入myschool
func i_myschool(school *School) error {
	sql := "INSERT INTO `myschool`(`msId`,`userId`,`type`,`title`,`contentId`,`studyYear`,`createTime`,`updateTime`) VALUES(?,?,?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{school.MsId, school.PrincipalUserId, school.MsType, school.Title, school.ContentId, school.StudyYear}})
}

//插入grade
func i_grade(school *School) error {
	sql := "INSERT INTO `grade`(`grId`,`grade`,`type`,`syId`,`createTime`,`updateTime`) VALUES(?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{school.GrId, school.Grade, school.Type, school.SyId}})
}

//插入workTime
func i_worktime(school *School) error {
	sql := "INSERT INTO `worktime`(`wtId`,`syId`,`grId`,`intoTime`,`outTime`,`beginDay`,`endDay`,`effectWeek`,`createTime`,`updateTime`)VALUES(?,?,?,?,?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{school.WtId, school.SyId, school.GrId, school.IntoTime, school.OutTime, school.BeginDay, school.EndDay, school.EffectWeek}})
}

//批量插入workday
func i_workday(school *School) error {
	if len(school.Gradess) < 1 {
		return errors.New("年级错误")
	}
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("INSERT INTO `workday`(`wdId`,`syId`,`grId`,`exceptionDay`,`work`,`createTime`,`updateTime`) VALUES")
	for i, Map := range school.Gradess {
		if i == 0 {
			sqlBuffer.WriteString("(")
		} else {
			sqlBuffer.WriteString(",(")
		}
		sqlBuffer.WriteString("'")
		sqlBuffer.WriteString(Map["wdId"])
		sqlBuffer.WriteString("',")

		sqlBuffer.WriteString("'")
		sqlBuffer.WriteString(school.SyId)
		sqlBuffer.WriteString("',")

		sqlBuffer.WriteString("'")
		sqlBuffer.WriteString(Map["grId"])
		sqlBuffer.WriteString("',")

		sqlBuffer.WriteString("'")
		sqlBuffer.WriteString(school.ExceptionDay)
		sqlBuffer.WriteString("',")

		sqlBuffer.WriteString("'")
		sqlBuffer.WriteString(Map["work"])
		sqlBuffer.WriteString("',")

		sqlBuffer.WriteString("NOW(),")
		sqlBuffer.WriteString("NOW())")
	}
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态修改workDay
func u_workday(school *School) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("UPDATE workday SET `updateTime`=NOW()")
	if school.Work != "" {
		sqlBuffer.WriteString(" ,work='")
		sqlBuffer.WriteString(school.Work)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" WHERE `syId`='")
	sqlBuffer.WriteString(school.SyId)
	sqlBuffer.WriteString("' AND grId='")
	sqlBuffer.WriteString(school.GrId)
	sqlBuffer.WriteString("' AND exceptionDay='")
	sqlBuffer.WriteString(school.ExceptionDay)
	sqlBuffer.WriteString("'")
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态修改school
func u_school(school *School) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("UPDATE school SET `updateTime`=NOW()")
	if school.ScInviteCode != "" {
		sqlBuffer.WriteString(" ,scInviteCode='")
		sqlBuffer.WriteString(school.ScInviteCode)
		sqlBuffer.WriteString("'")
	}
	if school.ScInviteQRCode != "" {
		sqlBuffer.WriteString(" ,scInviteQRCode='")
		sqlBuffer.WriteString(school.ScInviteQRCode)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" WHERE `scId`='")
	sqlBuffer.WriteString(school.ScId)
	sqlBuffer.WriteString("'")
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态修改schoolyear
func u_schoolyear(school *School) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("UPDATE schoolyear SET `updateTime`=NOW()")
	if school.BeginTime != "" {
		sqlBuffer.WriteString(" ,beginTime='")
		sqlBuffer.WriteString(school.BeginTime)
		sqlBuffer.WriteString("'")
	}
	if school.EndTime != "" {
		sqlBuffer.WriteString(" ,endTime='")
		sqlBuffer.WriteString(school.EndTime)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" WHERE `syId`='")
	sqlBuffer.WriteString(school.SyId)
	sqlBuffer.WriteString("'")
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//根据     学年，年级，异常日期     删除workDay
func d_workday(school *School) error {
	sql := "DELETE FROM workday WHERE `syId`=?  AND `exceptionDay`=?;"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{school.SyId, school.ExceptionDay}})
}

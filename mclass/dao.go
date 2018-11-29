package main

import (
	"bytes"
	"fmt"
	"mtcomm/db/mysql"
)

//SelectByScInviteCode接口所需sql
func selectByScInviteCode(sy *SchoolYear) (map[string]string, error) {
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select s.*,(select syId from schoolyear where scId=s.scId)syId from school s  where s.scInviteCode=? ", Args: []interface{}{sy.ScInviteCode}})
}
func findAll(sy *SchoolYear) ([]map[string]string, error) {
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: "select * from grade where syId=(select syId from school s join schoolyear sy on sy.scId=s.scId where scInviteCode=?) ", Args: []interface{}{sy.ScInviteCode}})
}
func selectClassInviteCode(sy *SchoolYear) (map[string]string, error) {
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from class where clInviteCode=?", Args: []interface{}{sy.ClInviteCode}})
}
func countClid(sy *SchoolYear) (int, error) {
	//return mysqlClient.Count(&mysql.Stmt{Sql: "select count(*)from student_class c1 where clId=(select c2.clId from class c2 where c2.clInviteCode=? )", Args: []interface{}{sy.ClInviteCode}})
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select count(*)from student_class c1 where c1.ad='a'")
	if sy.ClInviteCode != "" {
		sqlBuffer.WriteString(" and c1.clId=(select c2.clId from class c2 where c2.clInviteCode='")
		sqlBuffer.WriteString(sy.ClInviteCode)
		sqlBuffer.WriteString("')")
	}
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and c1.clId=")
		sqlBuffer.WriteString("'")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.Count(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

func selectByClid(sy *SchoolYear) (map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `class`  where 1=1")
	if sy.ClInviteCode != "" {
		sqlBuffer.WriteString(" and clInviteCode='")
		sqlBuffer.WriteString(sy.ClInviteCode)
		sqlBuffer.WriteString("'")
	}
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}
func selectBySyId(sy *SchoolYear) (map[string]string, error) {
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select sc.* from school sc where sc.scId=(select scId from schoolyear where syId=? )", Args: []interface{}{sy.SyId}})
}
func selectByGrid(sy *SchoolYear) (map[string]string, error) {
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from `grade`  where grId=? ", Args: []interface{}{sy.GrId}})
}

//创建班级sql
func addClass(client mysql.MysqlClient, sy *SchoolYear) error {
	return client.Execute(&mysql.Stmt{Sql: "insert into class values(?,?,?,?,?,?,?,NOW(),NOW())", Args: []interface{}{sy.ClId, sy.SyId, sy.GrId, sy.ClName, sy.DirectorUserId, sy.ClInviteCode, sy.ClInviteQRCode}})
}
func addTeacherClassMember(client mysql.MysqlClient, sy *SchoolYear) error {
	return client.Execute(&mysql.Stmt{Sql: "insert into  teacher_classmember(tcId,clId,userId,relation,type,flag,ad,createTime,updateTime,onClass)values(?,?,?,?,'0','1','a',NOW(),NOW(),'0')", Args: []interface{}{sy.TcId, sy.ClId, sy.DirectorUserId, sy.Relation}})
}
func addMySchool(client mysql.MysqlClient, sy *SchoolYear) error {
	return client.Execute(&mysql.Stmt{Sql: "insert into myschool (msId,userId,type,title,contentId,studyYear,createTime,updateTime)values (?,?,'1',?,?,?,NOW(),NOW())", Args: []interface{}{sy.MsId, sy.DirectorUserId, sy.Title, sy.ClId, sy.StudyYear}})
}
func selectByTcId(client mysql.MysqlClient, sy *SchoolYear) (map[string]string, error) {
	return client.SearchOneRow(&mysql.Stmt{Sql: "select * from teacher_classmember where tcId=?", Args: []interface{}{sy.TcId}})
}

//老师加入班级sql
func checkTeach(client mysql.MysqlClient, sy *SchoolYear) (map[string]string, error) {
	return client.SearchOneRow(&mysql.Stmt{Sql: "select * from teacher_classmember where clId=? and userId=? and ad='a'", Args: []interface{}{sy.ClId, sy.DirectorUserId}})
}
func teacherJoinClass(client mysql.MysqlClient, sy *SchoolYear) error {
	return client.Execute(&mysql.Stmt{Sql: "insert into teacher_classmember(tcId,clId,userId,relation,type,flag,ad,updateTime,createTime,onClass) values(?,?,?,?,'1','1','a',NOW(),NOW(),'0')", Args: []interface{}{sy.TcId, sy.ClId, sy.DirectorUserId, sy.Relation}})
}
func addTeacher(client mysql.MysqlClient, sy *SchoolYear) error {
	return client.Execute(&mysql.Stmt{Sql: "insert into myschool(msId,userId,type,title,contentId,studyYear,createTime,updateTime) values (?,?,'3',?,?,?,NOW(),NOW())", Args: []interface{}{sy.MsId, sy.DirectorUserId, sy.Title, sy.ClId, sy.StudyYear}})
}
func updateClass(sy *SchoolYear) error {
	return mysqlClient.Execute(&mysql.Stmt{Sql: "update `class` set clInviteCode=? , clInviteQRCode=? ,updateTime=NOW() where clId=?", Args: []interface{}{sy.ClInviteCode, sy.ClInviteQRCode, sy.ClId}})
}

func selectById(client mysql.MysqlClient, sy *SchoolYear) (map[string]string, error) {
	return client.SearchOneRow(&mysql.Stmt{Sql: "select c1.* from class c1 where clId=?", Args: []interface{}{sy.ClId}})
}

func selectStudyYear(client mysql.MysqlClient, sy *SchoolYear) (map[string]string, error) {
	return client.SearchOneRow(&mysql.Stmt{Sql: "select sc.*,s.* from schoolYear sc join school s on s.scId=sc.scId where sc.syId=?", Args: []interface{}{sy.SyId}})
}

//成员管理
func selectAllStu(client mysql.MysqlClient, sy *SchoolYear) ([]map[string]string, error) {
	return client.SearchMutiRows(&mysql.Stmt{Sql: "select sc.* ,s.* from  student_class sc join student s on s.stId=sc.stId where sc.clId=? and sc.ad='a' ORDER BY sc.`studentNum` ASC", Args: []interface{}{sy.ClId}})
}
func selectAllFamlily(client mysql.MysqlClient, sy *SchoolYear) ([]map[string]string, error) {
	return client.SearchMutiRows(&mysql.Stmt{Sql: "select c.* from family_classmember c where c.clId=?  and c.ad='a'", Args: []interface{}{sy.ClId}})
}
func selectAllTeacher(client mysql.MysqlClient, sy *SchoolYear) ([]map[string]string, error) {
	return client.SearchMutiRows(&mysql.Stmt{Sql: "select c.* from teacher_classmember c where c.clId=? and c.flag='1'and c.ad='a' order by type asc", Args: []interface{}{sy.ClId}})
}
func patriarch_student(client mysql.MysqlClient, sy *SchoolYear) ([]map[string]string, error) {
	return client.SearchMutiRows(&mysql.Stmt{Sql: "select s.*,ps.* ,sc.* from student s join patriarch_student ps on ps.stId=s.stId join student_class sc on sc.stId=s.stId where ps.userId=? and sc.clId=? ", Args: []interface{}{sy.UserId, sy.ClId}})
}

//======

func operateStu(client mysql.MysqlClient, sy *SchoolYear) error {
	if sy.Type == "0" {
		return client.Execute(&mysql.Stmt{Sql: "update student_class sc set sc.ad='d',sc.updateTime=NOW() where sc.stId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	if sy.Type == "1" {
		return client.Execute(&mysql.Stmt{Sql: "update student_class sc set sc.status='1',sc.updateTime=NOW() where sc.stId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	if sy.Type == "4" {
		return client.Execute(&mysql.Stmt{Sql: "update student_class sc set sc.status='0',sc.updateTime=NOW() where sc.stId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	return nil
}

func operateTeacher(client mysql.MysqlClient, sy *SchoolYear) error {
	//"2"表示不再任教
	if sy.Type == "2" {
		return client.Execute(&mysql.Stmt{Sql: "update teacher_classmember sc set sc.onClass='1',sc.updateTime=NOW()  where sc.userId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	//"3"表示任教班主任
	if sy.Type == "3" {
		return client.Execute(&mysql.Stmt{Sql: "update teacher_classmember sc set sc.type='0',sc.onClass='0',sc.updateTime=NOW() where sc.userId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	//"0"表示删除老师
	if sy.Type == "0" {
		return client.Execute(&mysql.Stmt{Sql: "update teacher_classmember sc set sc.ad='d',sc.updateTime=NOW() where sc.userId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	return nil
}
func operateClassOwner(client mysql.MysqlClient, sy *SchoolYear) error {
	return client.Execute(&mysql.Stmt{Sql: "update teacher_classmember sc set sc.type='1' ,sc.onclass='0',sc.updateTime=NOW() where sc.userId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.DisposeUserId, sy.ClId}})
}
func operateFamily(client mysql.MysqlClient, sy *SchoolYear) error {
	if sy.Type == "0" {
		return client.Execute(&mysql.Stmt{Sql: "update family_classmember sc set sc.ad='d',sc.updateTime=NOW() where sc.userId=? and sc.clId=? and sc.ad='a'", Args: []interface{}{sy.UserId, sy.ClId}})
	}
	return nil
}
func checkTeacher(client mysql.MysqlClient, sy *SchoolYear) (map[string]string, error) {
	return client.SearchOneRow(&mysql.Stmt{Sql: "select `type` from teacher_classmember where userId=? and clId=? and ad='a'", Args: []interface{}{sy.DisposeUserId, sy.ClId}})
}
func checkClId(client mysql.MysqlClient, sy *SchoolYear) (int, error) {
	return client.Count(&mysql.Stmt{Sql: "select count(*) from teacher_classmember where clId=? and ad='a'", Args: []interface{}{sy.ClId}})
}

func updateMschool(sy *SchoolYear) (error) {
	return mysqlClient.Execute(&mysql.Stmt{Sql: "update myschool ms set ms.type=? ,ms.updateTime=NOW() where ms.contentId=? and ms.userId=? and ms.type='3'", Args: []interface{}{sy.Type, sy.ContentId, sy.UserId}})
}
func ReupdateMschool(sy *SchoolYear) (error) {
	return mysqlClient.Execute(&mysql.Stmt{Sql: "update myschool ms set ms.type=? ,ms.updateTime=NOW() where ms.contentId=? and ms.userId=? and ms.type='1'", Args: []interface{}{sy.Type, sy.ContentId, sy.UserId}})
}

//===================================
//动态查询班级
func s_class(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `class` where 1=1")
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	if sy.SyId != "" {
		sqlBuffer.WriteString(" and syId='")
		sqlBuffer.WriteString(sy.SyId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})

}

//动态patriarch_student
func s_patriarch_student(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `patriarch_student` where 1=1")
	if sy.PsId != "" {
		sqlBuffer.WriteString(" and psId='")
		sqlBuffer.WriteString(sy.PsId)
		sqlBuffer.WriteString("'")
	}
	if sy.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(sy.UserId)
		sqlBuffer.WriteString("'")
	}
	if sy.StId != "" {
		sqlBuffer.WriteString(" and stId='")
		sqlBuffer.WriteString(sy.StId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态查询学生
func s_student(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `student` where 1=1")
	if sy.StId != "" {
		sqlBuffer.WriteString(" and stId='")
		sqlBuffer.WriteString(sy.StId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态查询family_classmember
func s_family_classmember(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `family_classmember` where ad='a'")
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	if sy.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(sy.UserId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态查询teacher_classmember
func s_teacher_classmember(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `teacher_classmember` where 1=1")
	if sy.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(sy.UserId)
		sqlBuffer.WriteString("'")
	}
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	if sy.Type != "" {
		sqlBuffer.WriteString(" and `type`='")
		sqlBuffer.WriteString(sy.Type)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" and `ad`='a'")
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态查询student_class
func s_student_class(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `student_class` where ad='a'")
	if sy.StId != "" {
		sqlBuffer.WriteString(" and stId='")
		sqlBuffer.WriteString(sy.StId)
		sqlBuffer.WriteString("'")
	}
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	if sy.StudentNum != "" {
		sqlBuffer.WriteString(" and studentNum='")
		sqlBuffer.WriteString(sy.StudentNum)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态查询check_join
func s_check_join(sy *SchoolYear) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("select * from `check_join` where ad='a'")
	if sy.CjId != "" {
		sqlBuffer.WriteString(" and cjId='")
		sqlBuffer.WriteString(sy.CjId)
		sqlBuffer.WriteString("'")
	}
	if sy.ClId != "" {
		sqlBuffer.WriteString(" and clId='")
		sqlBuffer.WriteString(sy.ClId)
		sqlBuffer.WriteString("'")
	}
	if sy.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(sy.UserId)
		sqlBuffer.WriteString("'")
	}
	if sy.Name != "" {
		sqlBuffer.WriteString(" and name='")
		sqlBuffer.WriteString(sy.Name)
		sqlBuffer.WriteString("'")
	}
	if sy.StudentNum != "" {
		sqlBuffer.WriteString(" and studentNum='")
		sqlBuffer.WriteString(sy.StudentNum)
		sqlBuffer.WriteString("'")
	}
	if sy.Birthday != "" {
		sqlBuffer.WriteString(" and birthday='")
		sqlBuffer.WriteString(sy.Birthday)
		sqlBuffer.WriteString("'")
	}
	if sy.Sex != "" {
		sqlBuffer.WriteString(" and sex='")
		sqlBuffer.WriteString(sy.Sex)
		sqlBuffer.WriteString("'")
	}
	if sy.Relation != "" {
		sqlBuffer.WriteString(" and relation='")
		sqlBuffer.WriteString(sy.Relation)
		sqlBuffer.WriteString("'")
	}
	if sy.Status != "" {
		sqlBuffer.WriteString(" and status='")
		sqlBuffer.WriteString(sy.Status)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//根据班级查询所有操作记录
func s_check_join2(sy *SchoolYear) ([]map[string]string, error) {
	sql := "SELECT * FROM `check_join`WHERE `clId`=? ORDER BY `status`,ad,createTime DESC"
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.ClId}})
}

//插入family_classmember表
func i_family_classmember(sy *SchoolYear) error {
	sql := "INSERT INTO `family_classmember`(`fcId`,`clId`,`userId`,`createTime`,`updateTime`) VALUES(?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.FcId, sy.ClId, sy.UserId}})
}

//插入myschool表
func i_myschool(sy *SchoolYear) error {
	sql := "INSERT INTO `myschool`(`msId`,`userId`,`type`,`title`,`contentId`,`stId`,`studyYear`,`createTime`,`updateTime`)VALUES(?,?,?,?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.MsId, sy.UserId, sy.Type, sy.Title, sy.ContentId, sy.StId, sy.StudyYear}})
}

//插入patriarch_student表
func i_patriarch_student(sy *SchoolYear) error {
	sql := "INSERT INTO `patriarch_student`(`psId`,`userId`,`stId`,`relation`) VALUES(?,?,?,?);"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.PsId, sy.UserId, sy.StId, sy.Relation}})
}

//插入student表
func i_student(sy *SchoolYear) error {
	sql := "INSERT INTO `student`(`stId`,`name`,`facePath`,`birthday`,`sex`,`createTime`,`updateTime`)VALUES(?,?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.StId, sy.Name, sy.FacePath, sy.Birthday, sy.Sex}})
}

//插入student_class表
func i_student_class(sy *SchoolYear) error {
	sql := "INSERT INTO `student_class`(`sclId`,`stId`,`clId`,`studentNum`,`createTime`,`updateTime`)VALUES(?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.SclId, sy.StId, sy.ClId, sy.StudentNum}})
}

//插入check_join表
func i_check_join(sy *SchoolYear) error {
	sql := "INSERT INTO `check_join`(`cjId`,`clId`,`userId`,`name`,`studentNum`,`birthday`,`sex`,`relation`,`createTime`,`updateTime`)VALUES(?,?,?,?,?,?,?,?,NOW(),NOW());"
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{sy.CjId, sy.ClId, sy.UserId, sy.Name, sy.StudentNum, sy.Birthday, sy.Sex, sy.Relation}})
}

//动态修改check_join
func u_check_join(sy *SchoolYear) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("UPDATE check_join SET `updateTime`=NOW()")
	if sy.Status != "" {
		sqlBuffer.WriteString(" ,status='")
		sqlBuffer.WriteString(sy.Status)
		sqlBuffer.WriteString("'")
	}
	if sy.Ad != "" {
		sqlBuffer.WriteString(" ,ad='")
		sqlBuffer.WriteString(sy.Ad)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" WHERE `cjId`='")
	sqlBuffer.WriteString(sy.CjId)
	sqlBuffer.WriteString("'")
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//动态修改student
func u_student(sy *SchoolYear) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("UPDATE student SET `updateTime`=NOW()")
	if sy.FacePath != "" {
		sqlBuffer.WriteString(" ,facePath='")
		sqlBuffer.WriteString(sy.FacePath)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" WHERE `stId`='")
	sqlBuffer.WriteString(sy.StId)
	sqlBuffer.WriteString("'")
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//修改student
func u_patriarch_student(sy *SchoolYear) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("UPDATE patriarch_student SET psId=psId ")
	if sy.StId != "" {
		sqlBuffer.WriteString(" ,stId='")
		sqlBuffer.WriteString(sy.StId)
		sqlBuffer.WriteString("'")
	}
	if sy.Relation != "" {
		sqlBuffer.WriteString(" ,relation='")
		sqlBuffer.WriteString(sy.Relation)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" WHERE `psId`='")
	sqlBuffer.WriteString(sy.PsId)
	sqlBuffer.WriteString("'")
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//删除视图
func myschool(sy *SchoolYear) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("delete from myschool where 1=1")
	if sy.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(sy.UserId)
		sqlBuffer.WriteString("'")
	}
	if sy.Type != "" {
		sqlBuffer.WriteString(" and type='")
		sqlBuffer.WriteString(sy.Type)
		sqlBuffer.WriteString("'")
	}
	if sy.ContentId != "" {
		sqlBuffer.WriteString(" and contentId='")
		sqlBuffer.WriteString(sy.ContentId)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}
func update_teacherMember(sy *SchoolYear) error {
	return mysqlClient.Execute(&mysql.Stmt{Sql: "update teacher_classmember set relation=? where userId=? and clId=? and ad='a'", Args: []interface{}{sy.Relation, sy.UserId, sy.ClId}})
}

//更新视图
func update_myschool(sy *SchoolYear) error {
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("update  myschool SET `updateTime`=NOW()")
	if sy.Title != "" {
		sqlBuffer.WriteString(" , title='")
		sqlBuffer.WriteString(sy.Title)
		sqlBuffer.WriteString("'")
	}
	sqlBuffer.WriteString(" where 1=1")
	if sy.MsId != "" {
		sqlBuffer.WriteString(" and msId='")
		sqlBuffer.WriteString(sy.MsId)
		sqlBuffer.WriteString("'")
	}
	if sy.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(sy.UserId)
		sqlBuffer.WriteString("'")
	}
	if sy.ContentId != "" {
		sqlBuffer.WriteString("and contentId='")
		sqlBuffer.WriteString(sy.ContentId)
		sqlBuffer.WriteString("'")
	}
	if sy.Type != "" {
		sqlBuffer.WriteString(" and type='")
		sqlBuffer.WriteString(sy.Type)
		sqlBuffer.WriteString("'")
	}
	sql := sqlBuffer.String()
	log.Debug("sql", fmt.Sprint(sql))
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

//查询学生所在学校
func SearchScIdByStId(client mysql.MysqlClient, stId string) (map[string]string, error) {
	return client.SearchOneRow(&mysql.Stmt{"SELECT s.scId, MAX(sc.`createTime`) FROM `student_class` sc,`class` c,`schoolyear` s WHERE sc.`clId`=c.`clId` AND s.`syId`=c.`syId` AND sc.`stId`=?", []interface{}{stId}})
}

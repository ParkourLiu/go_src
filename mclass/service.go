package main

import (
	"bytes"
	chat "chat/client"
	"errors"
	"fmt"
	"image/jpeg"
	"strconv"
	"time"

	"encoding/base64"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/golang/go/src/pkg/io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	push "push/client"
	"strings"
)

type SchoolYear struct {
	ScInviteCode   string              `json:"scInviteCode"`   //学校班主任邀请码
	ClId           string              `json:"clId"`           //班级主键
	SyId           string              `json:"syId"`           //学校主键
	ScId           string              `json:"scId"`           //学年表主键
	GrId           string              `json:"grId"`           //年级
	ClName         string              `json:"clName"`         //班级名字
	DirectorUserId string              `json:"directorUserId"` //班主任用户id
	ClInviteCode   string              `json:"clInviteCode"`   //班级邀请码
	ClInviteQRCode string              `json:"clInviteQRCode"` //班级邀请二维码
	Relation       string              `json:"relation"`       //称呼
	TcId           string              `json:"tcId"`           //班级成员表主键
	StudyYear      string              `json:"studyYear"`      //学年
	MsId           string              `json:"msId"`           //myschool表主键
	Title          string              `json:"title"`          //视图上显示的标题
	Flag           string              `json:"flag"`           //标志位:1表示班主任邀请码  2表示老师邀请码  3表示家长邀请码
	StId           string              `json:"stId"`           //学生主键
	Name           string              `json:"name"`           //学生姓名
	StudentNum     string              `json:"studentNum"`     //学号
	Birthday       string              `json:"birthday`        //生日
	Sex            string              `json:"sex"`            //性别
	PsId           string              `json:"psId"`           //学生家长对应关系
	SclId          string              `json:"sclId"`          //学生-班级主键
	UserId         string              `json:"userId"`
	Status         string              `json:"status"`
	DisposeUserId  string              `json:"disposeUserId"`
	Type           string              `json:"type"`
	Sign           string              `json:"sign"`
	FcId           string              `json:"fcId"`      //家长id
	ContentId      string              `json:"contentId"` //内容Id
	FacePath       string              `json:"facePath"`  //学生人脸图路径
	CjId           string              `json:"cjId"`      //审核id
	Ad             string              `json:"ad"`
	Tmember        []map[string]string `json:"tmember"`
	GroupChatId    string              `json:"groupChatId"`  //群聊id
	GroupChatUrl   string              `json:"groupChatUrl"` //群聊头像
}

type createClassService struct {
}
type McreateClassService interface {
	SelectByScInviteCode(*SchoolYear) (map[string]interface{}, error, string) //邀请码认证
	CreateClass(*SchoolYear) (map[string]string, error, string)               //创建班级
	TeacherJoinClass(*SchoolYear) (error, string)                             //老师加入班级
	FamilyJoinClass(*SchoolYear) (error, string)                              //家长加入班级
	NewMember(*SchoolYear) ([]map[string]string, error, string)               //新成员信息
	ApproveMembers(*SchoolYear) (error, string)                               //同意新成员加入
	UpdateStudent(*SchoolYear) (error, string)                                //修改学生信息
	ManagerMember(group *SchoolYear) ([]map[string]string, error, string)     //成员管理
	OperateMember(group *SchoolYear) (error, string)                          //对班级成员的操作
	ClassQrCode(group *SchoolYear) (map[string]string, error, string)         //班级信息
	FindAllMember(*SchoolYear) ([]map[string]string, error, string)
	FindTeacherMember(*SchoolYear) ([]map[string]string, error, string) //获取所有老师信息
	UpdateGroupChatInfo(*SchoolYear) (error, string)                    //修改群聊头像
	UpdateTeachInfo(*SchoolYear) (error, string)                        //修改群聊老师信息
}
type Push struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	UserId    string `json:"userId"`
	ClId      string `json:"clId"`
	HaveRobot string `json:"haveRobot"`
	StId      string `json:"stId"`
	Type      string `json:"type"`
	Name      string `json:"name"`
}

//班主任邀请码验证
func (service createClassService) SelectByScInviteCode(sy *SchoolYear) (map[string]interface{}, error, string) {
	log.Debug("method_start", "SelectByScInviteCode", "input", "sy=="+fmt.Sprint(sy))
	fg, _ := strconv.Atoi(sy.Flag)
	message := map[string]interface{}{}
	if fg == 0 {
		data, err := selectByScInviteCode(sy)
		if data == nil || err != nil {
			return nil, err, flag1
		}
		message["schoolDetail"] = data
		grade, err := findAll(sy)
		if err != nil {
			return nil, err, ""
		}
		message["gradeDetail"] = grade
		log.Debug("method_end", "SelectByScInviteCode", "status", "success", data)
		return message, err, "100"
	}
	if fg == 1 {
		var m map[string]string
		m = make(map[string]string)
		data, err := selectClassInviteCode(sy)
		if data == nil || err != nil {
			return nil, err, flag1
		}
		count, err1 := countClid(sy)
		if err1 != nil {
			return nil, err, ""
		}
		m["count"] = strconv.Itoa(count)
		classData, err2 := selectByClid(sy)
		if classData == nil || err2 != nil {
			return nil, err2, ""
		}
		//m["grId"] = classData["grId"]
		m["clName"] = classData["clName"]
		m["clId"] = classData["clId"]
		s := &SchoolYear{SyId: classData["syId"]}
		schoolData, err3 := selectBySyId(s)
		if schoolData == nil || err3 != nil {
			return nil, err3, ""
		}
		m["scName"] = schoolData["scName"]
		g := &SchoolYear{GrId: classData["grId"]}
		grade, err4 := selectByGrid(g)
		if err4 != nil {
			return nil, err, ""
		}
		m["grade"] = grade["grade"]
		message["classDetail"] = m
		log.Debug("method_end", "SelectByScInviteCode", "status", "success", m)
		return message, nil, "100"
	}
	return nil, nil, ""
}

//老师加入班级
func (service createClassService) TeacherJoinClass(sy *SchoolYear) (error, string) {
	log.Debug("method_start", "TeacherJoinClass", "input", "sy=="+fmt.Sprint(sy))
	count, errs := checkClssId(sy)
	if errs != nil {
		return errs, ""
	}
	if count == 0 {
		return errs, flag4
	}
	id := idGenClient.GetUniqueId()
	//id:="123w1sx"
	syear := getSYear()

	sy = &SchoolYear{MsId: "ms" + id, ClId: sy.ClId, DirectorUserId: sy.DirectorUserId, Relation: sy.Relation, TcId: "tc" + id, StudyYear: syear}

	//核对普通老师是否二次加入班级
	data, err1 := checkTeach(mysqlClient, sy)
	if data != nil {
		return err1, flag2
	}
	//获取学年id
	syId, err := selectByClid(sy)
	if err != nil {
		return err, ""
	}
	sy = &SchoolYear{MsId: "ms" + id, ClId: sy.ClId, SyId: syId["syId"], DirectorUserId: sy.DirectorUserId, Relation: sy.Relation, TcId: "tc" + id, StudyYear: syear, DisposeUserId: sy.DisposeUserId}
	err = teacherJoinClass(mysqlClient, sy)
	if err != nil {
		return err, ""
	}
	rel, err := selectByTcId(mysqlClient, sy)
	if err != nil {
		return err, ""
	}
	title := rel["relation"]
	sy = &SchoolYear{MsId: "ms" + id, ClId: sy.ClId, SyId: syId["syId"], DirectorUserId: sy.DirectorUserId, Relation: sy.Relation, TcId: "tc" + id, Title: title, StudyYear: syear, DisposeUserId: sy.DisposeUserId}
	err = addTeacher(mysqlClient, sy)
	if err != nil {
		return err, ""
	}
	person := []string{sy.DirectorUserId}
	groupchat := &chat.GroupChatInfo{AcId: sy.ClId, UserIdArray: person}
	response, err := chatClient.JoinGroupChat(groupchat)
	if err != nil {
		return err, ""
	}
	if response.Err != "" {
		return errors.New(response.Err), ""
	}
	log.Debug("method_end", "TeacherJoinClass", "status", "success", err)
	return nil, "100"
}

//创建班级
func (service createClassService) CreateClass(sy *SchoolYear) (map[string]string, error, string) {
	log.Debug("method_start", "CreateClass", "input", fmt.Sprint("sy", sy))
	gchat := &SchoolYear{SyId: sy.SyId, GrId: sy.GrId, ClName: sy.ClName, DirectorUserId: sy.DirectorUserId}
	var m map[string]string
	id := idGenClient.GetUniqueId()
	//id :="123qaz"
	//获取二维码
	clInviteCode, clInviteQRCode, err1 := InviteCode()
	if err1 != nil {
		return nil, err1, ""
	}
	sy = &SchoolYear{ClId: "cl" + id, ClInviteCode: clInviteCode, ClInviteQRCode: clInviteQRCode, TcId: "tc" + id, SyId: sy.SyId, Relation: sy.Relation, ClName: sy.ClName, GrId: sy.GrId, DirectorUserId: sy.DirectorUserId, MsId: "ms" + id}
	//添加数据到班级表
	err := addClass(mysqlClient, sy)
	if err != nil {
		return nil, err, ""
	}
	//添加数据到班级成员表
	err = addTeacherClassMember(mysqlClient, sy)
	if err != nil {
		return nil, err, ""
	}
	rel, err := selectByTcId(mysqlClient, sy)
	if rel == nil || err != nil {
		return nil, err, ""
	}
	title := rel["relation"]
	syear := getSYear()
	sy = &SchoolYear{MsId: "ms" + id, ClId: "cl" + id, DirectorUserId: sy.DirectorUserId, Title: title, StudyYear: syear, SyId: sy.SyId}
	//保存数据到myschool表
	err = addMySchool(mysqlClient, sy)
	if err != nil {
		return nil, err, ""
	}
	sy = &SchoolYear{SyId: sy.SyId, ClId: "cl" + id}
	classData, err := selectById(mysqlClient, sy)
	if classData == nil || err != nil {
		return nil, err, ""
	}
	clInviteCode = classData["clInviteCode"]
	clInviteQRCode = classData["clInviteQRCode"]
	schoolData, err := selectBySyId(sy)
	if schoolData == nil || err != nil {
		return nil, err, ""
	}
	m = make(map[string]string)
	m["clInviteCode"] = clInviteCode
	m["clInviteQRCode"] = clInviteQRCode
	m["syear"] = syear
	//创建班级群聊
	grade, err := selectByGrid(gchat)
	if err != nil {
		return nil, err, ""
	}
	//年级名字
	gradeName := grade["grade"]
	school, err := selectStudyYear(mysqlClient, gchat)
	if err != nil {
		return nil, err, ""
	}
	//学校名字
	scName := school["scName"]

	headPic, err := redisClient.Hget("U:"+gchat.DirectorUserId, "imageName")
	if err != nil {
		return nil, err, ""
	}
	groupchat := &chat.GroupChatInfo{AcId: "cl" + id, GroupChatName: scName + gradeName + gchat.ClName, GroupChatUrl: headPic, GroupChatNotice: "毛毛公告", CreateUserId: gchat.DirectorUserId}
	log.Debug("groupchat❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤❤", fmt.Sprint(groupchat))
	response, err := chatClient.CreateGroupChat(groupchat)
	if err != nil {
		return nil, err, ""
	}
	if response.Err != "" {
		return nil, errors.New(response.Err), ""
	}
	log.Debug("method_end", "CreateClass", "status", "success", m)
	return m, nil, "100"
}

//家长申请加入班级 birthday clId directorUserId家长 name relation称呼 sex studentNum学号
func (service createClassService) FamilyJoinClass(sy *SchoolYear) (error, string) {
	log.Debug("method_start", "FamilyJoinClass", "input", "sy=="+fmt.Sprint(sy))
	clList, err := s_class(&SchoolYear{ClId: sy.ClId}) //查询此班级
	if err != nil {
		return err, ""
	}
	if len(clList) != 1 {
		return nil, flag4 //班级id不存在
	}
	//根据传进来的数据查询check_join
	cjList, err := s_check_join(&SchoolYear{ClId: sy.ClId, UserId: sy.DirectorUserId, Name: sy.Name, StudentNum: sy.StudentNum, Birthday: sy.Birthday, Sex: sy.Sex, Relation: sy.Relation})
	if err != nil {
		return err, ""
	}
	if len(cjList) > 1 {
		return err, flag10 //数据异常
	}
	if len(cjList) == 1 { //已申请，状态可能待审核，可能已拒绝
		if cjList[0]["status"] == "0" { //待审核
			return nil, flag11 //返回正在审核
		} else if cjList[0]["status"] == "2" { //已拒绝
			err := u_check_join(&SchoolYear{CjId: cjList[0]["cjId"], Status: "0"}) //修改为待审核
			if err != nil {
				return err, ""
			}
		}
	} else if len(cjList) == 0 { //没有申请过
		//插入check_join
		id := idGenClient.GetUniqueId() //"111111111111" //
		sy.CjId = "cj" + id
		sy.UserId = sy.DirectorUserId
		err := i_check_join(sy) //sy.CjId, sy.ClId, sy.UserId, sy.Name, sy.StudentNum, sy.Birthday, sy.Sex, sy.Relation
		if err != nil {
			return err, ""
		}
	}
	log.Debug("method_end", "FamilyJoinClass", "status", "success", nil)
	return nil, "100"
}

//新成员  clId  sign 0：待审核班级成员；1：所有班级成员 (包括待审核，已拒绝);
func (service createClassService) NewMember(sy *SchoolYear) ([]map[string]string, error, string) {
	log.Debug("method_start", "NewMember", "input", "sy=="+fmt.Sprint(sy))
	if sy.Sign == "0" { //待审核
		cjList, err := s_check_join(&SchoolYear{ClId: sy.ClId, Status: "0"})
		if err != nil {
			return cjList, err, ""
		}
		for i, cjMap := range cjList {
			imageName, err := redisClient.Hget("U:"+cjMap["userId"], "imageName")
			if err != nil {
				return cjList, err, ""
			}
			cjList[i]["imageName"] = imageName
		}
		return cjList, nil, "100"
	} else if sy.Sign == "1" { //所有处理历史
		cjList, err := s_check_join2(&SchoolYear{ClId: sy.ClId}) //查询所有操作过的历史
		if err != nil {
			return cjList, err, ""
		}
		for i, cjMap := range cjList {
			imageName, err := redisClient.Hget("U:"+cjMap["userId"], "imageName")
			if err != nil {
				return cjList, err, ""
			}
			if cjList[i]["ad"] == "d" { //代表已同意加入
				cjList[i]["status"] = "1" //返回给前端时修改状态为已加入
				cjList[i]["ad"] = "a"
			}
			cjList[i]["imageName"] = imageName
		}
		return cjList, nil, "100"
	} else {
		return nil, nil, flag10
	}
}

//班主任审核家长学生加入班级 cjId   disposeUserId   status 1同意，2拒绝
func (service createClassService) ApproveMembers(sy *SchoolYear) (error, string) {
	log.Debug("method_start", "ApproveMembers", "input", "sy=="+fmt.Sprint(sy))
	cjList, err := s_check_join(&SchoolYear{CjId: sy.CjId})
	if err != nil {
		return err, ""
	}
	if len(cjList) != 1 {
		return nil, flag12 //id不存在
	}
	tcList, err := s_teacher_classmember(&SchoolYear{ClId: cjList[0]["clId"], UserId: sy.DisposeUserId, Status: "0"}) //查询此用户是否是班主任
	if err != nil {
		return err, ""
	}
	if len(tcList) != 1 {
		return err, flag3 //无权操作
	}

	//判断是否同意加入
	if sy.Status == "2" { //拒绝加入
		err := u_check_join(&SchoolYear{CjId: sy.CjId, Status: "2"}) //审核状态改成已拒绝
		if err != nil {
			return err, ""
		}
	} else if sy.Status == "1" { //同意加入
		chsStId, chs, err := classHaveStudent(cjList[0]["clId"], cjList[0]["studentNum"], cjList[0]["name"], cjList[0]["sex"], cjList[0]["birthday"]) //判断此班是否有此学生，（班级id,学号）
		if err != nil {
			return err, ""
		}
		chf, err := classHaveFamily(cjList[0]["clId"], cjList[0]["userId"]) //判断此班是否有此家长，（班级id,家长用户id）
		if err != nil {
			return err, ""
		}
		fhsStId, fhs, err := familyHaveStudent(cjList[0]["userId"], cjList[0]["name"], cjList[0]["sex"], cjList[0]["birthday"], cjList[0]["clId"], cjList[0]["studentNum"], cjList[0]["relation"]) //判断此家长是否有此学生，（userId , name , sex , birthday ）
		if err != nil {
			return err, ""
		}

		stFlag := 0  //学生表是否插入数据标识
		fcFlag := 0  //家长成员是否插入数据标识
		psFlag := 0  //家长学生关系表是否插入数据标识
		msFlag := 0  //我的校园是否插入数据标识
		sclFlag := 0 //学生班级表是否插入数据标识
		stId := ""   //最终要插入到各个表的stid
		if chs == true { //此班有此学生
			stId = chsStId
			if chf == true { //此班有此家长（此情况下临时解决中已决定此家长必定有此学生，所以只需要给家长补一个视图就好）
				msFlag = 1
			}
			if chf == false { //此班无此家长
				if fhs == true { //此家长有此学生
					fcFlag = 1
					msFlag = 1
				} else { //此家长无此学生
					fcFlag = 1
					msFlag = 1
					psFlag = 1
				}
				c := &SchoolYear{CjId: sy.CjId, Status: ""}
				cj, err := s_check_join(c)
				if err != nil {
					return err, ""
				}
				ll := len(cj)
				log.Debug("len::::::::::::::::::", fmt.Sprint(ll))
				//获取班级id
				clId := cj[0]["clId"]

				//用户id
				userId := cj[0]["userId"]
				userIdArrayu := []string{userId}
				cl := &chat.GroupChatInfo{UserIdArray: userIdArrayu, AcId: clId}
				response, err := chatClient.JoinGroupChat(cl)
				if err != nil {
					return err, ""
				}
				if response.Err != "" {
					return errors.New(response.Err), ""
				}
			}
		} else { //此班无此学生
			if chf == true { //此班有此家长
				if fhs == true { //此家长有此学生
					stId = fhsStId
					sclFlag = 1
					msFlag = 1
				} else { //此家长无此学生
					stFlag = 1
					msFlag = 1
					psFlag = 1
					sclFlag = 1
				}
			} else { //此班无此家长
				if fhs == true { //此家长有此学生
					stId = fhsStId
					sclFlag = 1
					msFlag = 1
					fcFlag = 1
				} else { //此家长无此学生
					stFlag = 1
					msFlag = 1
					psFlag = 1
					sclFlag = 1
					fcFlag = 1
				}
				c := &SchoolYear{CjId: sy.CjId, Status: ""}
				cj, err := s_check_join(c)
				if err != nil {
					return err, ""
				}
				ll := len(cj)
				log.Debug("len::::::::::::::::::", fmt.Sprint(ll))
				//获取班级id
				clId := cj[0]["clId"]

				//用户id
				userId := cj[0]["userId"]
				userIdArrayu := []string{userId}
				cl := &chat.GroupChatInfo{UserIdArray: userIdArrayu, AcId: clId}
				response, err := chatClient.JoinGroupChat(cl)
				if err != nil {
					return err, ""
				}
				if response.Err != "" {
					return errors.New(response.Err), ""
				}
			}
		}

		//状态为1的都插入到数据库
		id := idGenClient.GetUniqueId() //"1110000000" //
		if stId == "" { //如果学生id不存在，则创建一个
			stId = "st" + id
		}
		if stFlag == 1 {
			//插入student表
			sy.StId = stId
			sy.Name = cjList[0]["name"]
			sy.Birthday = cjList[0]["birthday"]
			sy.Sex = cjList[0]["sex"]
			sy.UserId = cjList[0]["userId"]
			code, err := createObjecte(sy) //给学生生成一个宝贝
			if err != nil || code != "" {
				return err, code
			}
			err = i_student(sy) // sy.StId, sy.Name, sy.FacePath, sy.Birthday, sy.Sex
			if err != nil {
				return err, ""
			}
		}
		if fcFlag == 1 {
			//插入家长表
			sy.FcId = "fc" + id
			sy.UserId = cjList[0]["userId"]
			sy.ClId = cjList[0]["clId"]
			err = i_family_classmember(sy) // sy.FcId, sy.ClId, sy.UserId,
			if err != nil {
				return err, ""
			}
		}

		if psFlag == 1 {
			//插入patriarch_student表
			sy.PsId = "ps" + id
			sy.UserId = cjList[0]["userId"]
			sy.Relation = cjList[0]["relation"]
			sy.StId = stId
			err = i_patriarch_student(sy) //sy.PsId, sy.UserId, sy.StId,sy.Relation
			if err != nil {
				return err, ""
			}

		}

		if sclFlag == 1 {
			//插入student_class表
			sy.SclId = "scl" + id
			sy.StId = stId
			sy.ClId = cjList[0]["clId"]
			sy.StudentNum = cjList[0]["studentNum"]
			err = i_student_class(sy) //sy.SclId, sy.StId, sy.ClId, sy.StudentNum
			if err != nil {
				return err, ""
			}
		}

		if msFlag == 1 {
			//插入myschool表
			sy.MsId = "ms" + id
			sy.UserId = cjList[0]["userId"]
			sy.Type = "2" //家长班级视图
			sy.Title = cjList[0]["name"]
			sy.ContentId = cjList[0]["clId"]
			sy.StId = stId
			sy.StudyYear = now2schoolYear()
			err = i_myschool(sy) //sy.MsId, sy.UserId, sy.Type, sy.Title, sy.ContentId,sy.StId, sy.StudyYear
			if err != nil {
				return err, ""
			}
		}
		err = u_check_join(&SchoolYear{CjId: sy.CjId, Ad: "d"}) //删除审核表中的此条数据
		if err != nil {
			return err, ""
		}
		//添加推送
		clId := cjList[0]["clId"]
		famUserId := cjList[0]["userId"]
		class, err := s_class(&SchoolYear{ClId: clId})
		if err != nil {
			return err, ""
		}
		//班级名字
		clName := class[0]["clName"]
		//年级id
		grId := class[0]["grId"]
		//用户id
		directoryUserId := class[0]["directorUserId"]
		gr, err := selectByGrid(&SchoolYear{GrId: grId})
		if err != nil {
			return err, ""
		}
		//年级名字
		grade := gr["grade"]
		tea, err := s_teacher_classmember(&SchoolYear{ClId: sy.ClId, UserId: directoryUserId})
		if err != nil {
			return err, ""
		}
		relation := tea[0]["relation"]
		text := relation + "已同意您加入" + grade + clName
		title := "同意加入班级申请"
		//根据clId查询class表获取syId
		cl := &SchoolYear{ClId: sy.ClId}
		sclass, err := s_class(cl)
		if err != nil {
			return err, ""
		}
		if len(sclass) != 0 {
			ssyId := &SchoolYear{SyId: sclass[0]["syId"]}
			schoolSyId, err := selectStudyYear(mysqlClient, ssyId)
			if err != nil {
				return err, ""
			}
			stu, err := s_student_class(&SchoolYear{StudentNum: cjList[0]["studentNum"], ClId: cjList[0]["clId"]})
			if err != nil {
				return err, ""
			}
			if len(stu) == 0 {
				push := &Push{Text: text, Title: title, UserId: famUserId, ClId: sy.ClId, StId: "st" + id, HaveRobot: schoolSyId["haveRobot"], Type: "AP", Name: cjList[0]["name"]}
				pushMessage(push)
			} else {
				push := &Push{Text: text, Title: title, UserId: famUserId, ClId: sy.ClId, StId: stu[0]["stId"], HaveRobot: schoolSyId["haveRobot"], Type: "AP", Name: cjList[0]["name"]}
				pushMessage(push)
			}

		}

		// push:=&Push{Text:text,Title:title,UserId:famUserId,ClId:sy.ClId,StId:"st" + id}
		//pushMessage(push)
	}
	log.Debug("method_end", "ApproveMembers", "status", "success", nil)
	return nil, "100"
}

//修改学生信息 操作者Userid,学生id，FacePath人脸图
func (service createClassService) UpdateStudent(sy *SchoolYear) (error, string) {
	log.Debug("method_start", "UpdateStudent", "input", "sy=="+fmt.Sprint(sy))
	psList, err := s_patriarch_student(&SchoolYear{UserId: sy.UserId}) //得到学生Id
	if err != nil {
		return err, ""
	}
	psFlag := 0
	for _, v := range psList {
		if v["stId"] == sy.StId {
			psFlag = 1
		}
	}
	if psFlag == 0 {
		return nil, flag3 //无权操作
	}
	err = u_student(&SchoolYear{FacePath: sy.FacePath, StId: sy.StId})
	if err != nil {
		return err, ""
	}
	//获取oss
	client, err := oss.New(oss_endpoint, key_id, key_secret)
	if err != nil {
		fmt.Println("err2:", err)
	}

	bucket, err := client.Bucket(bucket_name_face) //cketName,"qx-mtalk-test"
	if err != nil {
		fmt.Println("err3:", err)
	}
	//同步数据到人脸库
	acTk, err := getAccessToken()
	if err != nil {
		return err, ""
	}
	scInfo, err := SearchScIdByStId(mysqlClient, sy.StId)
	if err != nil {
		return err, ""
	}
	if len(scInfo) > 0 {
		red, err := bucket.GetObject(sy.FacePath)
		if err != nil {
			return err, ""
		}
		imgBytes, err := ioutil.ReadAll(red)
		if err != nil {
			return err, ""
		}
		baseImg := base64.StdEncoding.EncodeToString(imgBytes)
		url := "https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add?access_token=" + acTk
		param := "{\"image\":\"" + baseImg + "\",\"image_type\":\"BASE64\",\"group_id\":\"" + scInfo["scId"] + "\",\"user_id\":\"" + sy.StId + "\"}"
		resp, err := DoBytesPost(url, param)
		if err != nil {
			return err, ""
		}
		if len(resp) == 0 {
			return errors.New("同步人脸识别人脸库失败"), ""
		}
	}
	log.Debug("method_end", "UpdateStudent", "status", "success", nil)
	return nil, "100"
}

func (service createClassService) ManagerMember(gp *SchoolYear) ([]map[string]string, error, string) {
	log.Debug("method_start", "ManagerMember", "input", "=="+fmt.Sprint(gp))
	slice := []map[string]string{}

	//"0"表示学生
	if gp.Flag == "0" {
		stu, err := selectAllStu(mysqlClient, gp)
		if err != nil {
			return nil, err, ""
		}
		return stu, nil, "100"
	}

	//"1"表示家长
	if gp.Flag == "1" {
		//查询所有的家长
		family, err := selectAllFamlily(mysqlClient, gp)
		if err != nil {
			return nil, err, ""
		}
		//遍历所有家长信息
		for k, v := range family {
			sy := &SchoolYear{UserId: family[k]["userId"], ClId: gp.ClId}
			headPic, err := redisClient.Hget("U:"+v["userId"], "imageName")
			v["headPic"] = headPic
			pstu, err := patriarch_student(mysqlClient, sy)
			if err != nil {
				return nil, err, ""
			}

			relation := ""
			//if len(pstu) == 1 {
			//	stName := pstu[0]["name"]
			//	v["relation"] =stName+"的"+ pstu[0]["relation"]
			//	slice = append(slice, v)
			//}
			if len(pstu) == 1 {
				ad := pstu[0]["ad"]
				if ad == "d" {
					v["relation"] = ""
					slice = append(slice, v)
				} else {
					stName := pstu[0]["name"]
					v["relation"] = stName + "的" + pstu[0]["relation"]
					slice = append(slice, v)
				}

			} else {
				count := 0
				for s, stu := range pstu {
					count++
					if count > 2 {
						relation = relation[0:len(relation)-1] + "...."
						break
					}
					//relation = relation + stu["name"] + "的"+pstu[s]["relation"]+","
					if pstu[s]["ad"] == "a" {
						relation = relation + stu["name"] + "的" + pstu[s]["relation"] + ","
					}
					if pstu[s]["ad"] == "d" {
						count--
					}
				}
				if relation != "" {
					relation = relation[0 : len(relation)-1]
					v["relation"] = relation
					slice = append(slice, v)
				}
			}
		}
		return slice, nil, "100"
	}
	//"2"表示老师
	if gp.Flag == "2" {
		teacher, err := selectAllTeacher(mysqlClient, gp)
		if err != nil {
			return nil, err, ""
		}
		for _, v := range teacher {
			headPic, err := redisClient.Hget("U:"+v["userId"], "imageName")
			if err != nil {
				return nil, err, ""
			}
			v["headPic"] = headPic
		}
		return teacher, nil, "100"
	}
	log.Debug("method_end", "ManagerMember", "status", "success", "")
	return nil, nil, "100"
}

func (service createClassService) OperateMember(gp *SchoolYear) (error, string) {
	log.Debug("method_start", "OperateMember", "input", "=="+fmt.Sprint(gp))
	err, code := checkTeachManager(gp)
	if err != nil || code == flag3 {
		return err, code
	}
	//操作学生
	if gp.Flag == "0" {
		err := operateStu(mysqlClient, gp)
		if err != nil {
			return err, ""
		}
		//删除学生家长对应关系

		////如果是删除学生操作，则把他的家长移除群聊
		//if gp.Type == "0" {
		//	stu := &SchoolYear{StId: gp.UserId}
		//	s, err := s_patriarch_student(stu)
		//	log.Debug("sssssssssss",fmt.Sprint(s))
		//	if err != nil {
		//		return err, ""
		//	}
		//	if len(s)==0{
		//		return nil,flag10
		//	}
		//	 for k,_:=range s{
		//		 sc := &SchoolYear{UserId: s[k]["userId"]}
		//	 fam, err := s_patriarch_student(sc)
		//		 if err != nil {
		//			 return err, ""
		//		 }
		//		 //如果家长对应的孩子只有一个，那么说明此孩子就是家长的孩子
		//		 //删除学生时应该对应删除他的家长
		//		 log.Debug("ssssssssssssssssssssssssssssssss")
		//		 log.Debug("len",fmt.Sprint(len(fam)))
		//		 if len(fam) == 1 {
		//			 //groupchat := &chat.GroupChatInfo{AcId: gp.ClId}
		//			 ////获取群聊id
		//			 //response, err := chatClient.SearchChatInfo(groupchat)
		//			 //if err != nil {
		//			 //	return err, ""
		//			 //}
		//			 //if response.Err != "" {
		//			 //	return errors.New(response.Err), ""
		//			 //}
		//			 //data := response.Data
		//			 //groupChatId := data["groupChatId"]
		//			 //userId := []string{fam[0]["userId"]}
		//			 //userIdArray := userId
		//			 //chat := &chat.GroupChatInfo{GroupChatId: groupChatId, UserIdArray: userIdArray}
		//			 ////移除群聊
		//			 //result, err := chatClient.QuitGroupChat(chat)
		//			 //if err != nil {
		//			 //	return nil, ""
		//			 //}
		//			 //if result.Err != "" {
		//			 //	return errors.New(response.Err), ""
		//			 //}
		//			 //删除家长
		//			 f:=&SchoolYear{UserId:fam[0]["userId"],Type:"0",ContentId:gp.ClId}
		//			 log.Debug("f",fmt.Sprint(f))
		//			 err = operateFamily(mysqlClient, f)
		//			 if err != nil {
		//				 return err, ""
		//			 }
		//			 //删除家长视图
		//			 fam:=&SchoolYear{UserId:fam[0]["userId"],Type:"2",ContentId:gp.ClId}
		//			 log.Debug("fam",fmt.Sprint(fam))
		//			 err=myschool(fam)
		//			 if err!=nil{
		//				 return err,""
		//			 }
		//		 }
		//	 }
		//}

	}
	//操作家长
	if gp.Flag == "1" {
		err := operateFamily(mysqlClient, gp)
		if err != nil {
			return err, ""
		}
		teacher, err := s_teacher_classmember(gp)
		if err != nil {
			return err, ""
		}
		if len(teacher) == 0 {
			groupchat := &chat.GroupChatInfo{AcId: gp.ClId}
			//获取群聊id
			response, err := chatClient.SearchChatInfo(groupchat)
			if err != nil {
				return err, ""
			}
			if response.Err != "" {
				return errors.New(response.Err), ""
			}
			data := response.Data
			groupChatId := data["groupChatId"]
			userId := []string{gp.UserId}
			userIdArray := userId
			chat := &chat.GroupChatInfo{GroupChatId: groupChatId, UserIdArray: userIdArray}
			//移除群聊
			result, err := chatClient.QuitGroupChat(chat)
			if err != nil {
				return nil, ""
			}
			if result.Err != "" {
				return errors.New(response.Err), ""
			}
		}
		//删除家长视图
		fam := &SchoolYear{UserId: gp.UserId, Type: "2", ContentId: gp.ClId}
		err = myschool(fam)
		if err != nil {
			return err, ""
		}
	}
	//班主任操作老师
	if gp.Flag == "2" {
		err := operateTeacher(mysqlClient, gp)
		if err != nil {
			return err, ""
		}
		//将自身改变为老师，身份转移
		if gp.Type == "3" {
			//老师视图变班主任视图
			sy := &SchoolYear{ContentId: gp.ClId, UserId: gp.UserId, Type: "1"}
			err := updateMschool(sy)
			if err != nil {
				return err, ""
			}
			//班主任变老师
			err = operateClassOwner(mysqlClient, gp)
			if err != nil {
				return err, ""
			}
			//改变班主任视图变老师
			ss := &SchoolYear{ContentId: gp.ClId, UserId: gp.DisposeUserId, Type: "3"}
			err = ReupdateMschool(ss)
			if err != nil {
				return err, ""
			}
		}
		fan, err := s_family_classmember(gp)
		if err != nil {
			return err, ""
		}

		if gp.Type == "0" && len(fan) == 0 {
			groupchat := &chat.GroupChatInfo{AcId: gp.ClId}
			//获取群聊id
			response, err := chatClient.SearchChatInfo(groupchat)
			if err != nil {
				return err, ""
			}
			if response.Err != "" {
				return errors.New(response.Err), ""
			}
			data := response.Data
			groupChatId := data["groupChatId"]
			userId := []string{gp.UserId}
			userIdArray := userId
			chat := &chat.GroupChatInfo{GroupChatId: groupChatId, UserIdArray: userIdArray}
			//移除群聊
			result, err := chatClient.QuitGroupChat(chat)
			if err != nil {
				return nil, ""
			}
			if result.Err != "" {
				return errors.New(response.Err), ""
			}
		}
		//删除老师视图
		fam := &SchoolYear{UserId: gp.UserId, Type: "3", ContentId: gp.ClId}
		err = myschool(fam)
		if err != nil {
			return err, ""
		}
	}
	log.Debug("method_end", "OperateMember", "status", "success")
	return nil, "100"
}

func (service createClassService) ClassQrCode(gp *SchoolYear) (map[string]string, error, string) {
	log.Debug("method_start", "ClassQrCode", "input", "=="+fmt.Sprint(gp))
	data, err := selectByClid(gp)
	if gp.Flag == "0" {

		if err != nil {
			return data, err, ""
		}
		gp = &SchoolYear{ClId: gp.ClId, SyId: data["syId"], GrId: data["grId"]}
		sy, err := selectStudyYear(mysqlClient, gp)
		if err != nil {
			return nil, err, ""
		}
		count, err := countClid(gp)
		if err != nil {
			return nil, err, ""
		}
		grade, err := selectByGrid(gp)
		if err != nil {
			return nil, err, ""
		}
		data["studyYear"] = sy["studyYear"]
		data["scName"] = sy["scName"]
		data["count"] = strconv.Itoa(count)
		data["grade"] = grade["grade"]
	}
	if gp.Flag == "1" {
		//获取二维码
		clInviteCode, clInviteQRCode, err := InviteCode()
		if err != nil {
			return nil, err, ""
		}
		data, err := selectByClid(gp)
		if err != nil {
			return nil, err, ""
		}
		gp = &SchoolYear{ClInviteCode: clInviteCode, ClInviteQRCode: clInviteQRCode, ClId: data["clId"], SyId: data["syId"], GrId: data["grId"]}
		err = updateClass(gp)
		if err != nil {
			return nil, err, ""
		}
		gp = &SchoolYear{ClId: data["clId"], SyId: data["syId"], GrId: data["grId"]}
		data, err = selectByClid(gp)
		if err != nil {
			return nil, err, ""
		}
		sy, err := selectStudyYear(mysqlClient, gp)
		if err != nil {
			return nil, err, ""
		}
		count, err := countClid(gp)
		if err != nil {
			return nil, err, ""
		}
		grade, err := selectByGrid(gp)
		if err != nil {
			return nil, err, ""
		}
		data["studyYear"] = sy["studyYear"]
		data["scName"] = sy["scName"]
		data["count"] = strconv.Itoa(count)
		data["grade"] = grade["grade"]
		log.Debug("data:", fmt.Sprint(data))
		return data, err, "100"
	}
	log.Debug("method_end", "ClassQrCode", "status", "success", data)
	return data, nil, "100"
}

//返回班级所有成员（班主任，老师，家长）
func (service createClassService) FindAllMember(gp *SchoolYear) ([]map[string]string, error, string) {
	log.Debug("method_start", "FindAllMember", "input", "=="+fmt.Sprint(gp))
	data := []map[string]string{}
	teacherUserId := []string{} //去重用

	//查询所有老师
	teacher, err := s_teacher_classmember(gp)
	if err != nil {
		return nil, err, ""
	}
	for k, _ := range teacher {
		m := map[string]string{}
		headPic, err := redisClient.Hget("U:"+teacher[k]["userId"], "imageName")
		if err != nil {
			return nil, err, ""
		}
		teacher[k]["headPic"] = headPic
		m["userId"] = teacher[k]["userId"]
		m["headPic"] = teacher[k]["headPic"]
		m["nickName"] = teacher[k]["relation"]
		m["type"] = teacher[k]["type"]
		data = append(data, m)
		teacherUserId = append(teacherUserId, m["userId"]) //老师id，家长去重时用
	}

	//查询所有班级家长
	family, err := s_family_classmember(gp)
	if err != nil {
		return nil, err, ""
	}
	for k, _ := range family {
		//家长和老师去重
		isRedo := 0
		for _, userId := range teacherUserId {
			if family[k]["userId"] == userId { //家长和老师重复
				isRedo = 1
				break
			}
		}
		if isRedo == 1 { //如果重复则直接跳过本轮循环
			continue
		}
		m := map[string]string{}
		f := &SchoolYear{UserId: family[k]["userId"], ClId: gp.ClId}
		stuName := ""
		famName := ""
		pstu, err := patriarch_student(mysqlClient, f)
		log.Debug("fffffffffffffff", fmt.Sprint(pstu))
		if err != nil {
			return nil, err, ""
		}
		nickName := ""
		if len(pstu) == 1 {
			ad := pstu[0]["ad"]
			if ad == "a" {
				stuName = pstu[0]["name"]
				famName = pstu[0]["relation"]
				m["nickName"] = stuName + "的" + famName
			} else {
				m["nickName"] = ""
			}
		} else {
			count := 0
			for s, stu := range pstu {
				count++
				if count > 2 {
					nickName = nickName[0:len(nickName)-1] + "...."
					break
				}
				if pstu[s]["ad"] == "a" {
					nickName += stu["name"] + "的" + pstu[s]["relation"] + ","
				}
				if pstu[s]["ad"] == "d" {
					count--
				}
			}
			log.Debug("nickName", fmt.Sprint(nickName))
			if nickName != "" {
				nickName = nickName[0 : len(nickName)-1]
				m["nickName"] = nickName
			}
		}
		headPic, err := redisClient.Hget("U:"+family[k]["userId"], "imageName")
		if err != nil {
			return nil, err, ""
		}
		family[k]["headPic"] = headPic
		family[k]["type"] = "2"
		//组装map
		m["headPic"] = family[k]["headPic"]
		m["type"] = family[k]["type"]
		m["userId"] = family[k]["userId"]
		data = append(data, m)
	}
	log.Debug("method_end", "FindAllMember", "status", "success", data)
	return data, nil, "100"
}
func (service createClassService) FindTeacherMember(gp *SchoolYear) ([]map[string]string, error, string) {
	log.Debug("method_start", "FindTeacherMember", "input", "=="+fmt.Sprint(gp))
	data := []map[string]string{}
	tmemeber := gp.Tmember
	for k, _ := range tmemeber {
		clId := tmemeber[k]["clId"]
		userId := tmemeber[k]["userId"]
		teach := &SchoolYear{ClId: clId, UserId: userId}
		fam, err := s_family_classmember(teach)
		if err != nil {
			return nil, err, ""
		}
		member, err := s_teacher_classmember(teach)
		if err != nil {
			return nil, err, ""
		}
		//此人身份是老师
		log.Debug("fam", fmt.Sprint(len(fam)))
		log.Debug("member", fmt.Sprint(len(member)))

		if len(fam) == 0 {
			if len(member) == 0 {
				tmemeber[k]["relation"] = ""
			} else {
				relation := member[0]["relation"]
				tmemeber[k]["relation"] = relation
			}
			data = append(data, tmemeber[k])
		}
		//此人身份是家长
		if len(member) == 0 {
			pstu, err := patriarch_student(mysqlClient, teach)
			if err != nil {
				return nil, err, ""
			}
			if len(pstu) == 0 {
				tmemeber[k]["relation"] = ""
			}
			nickName := ""
			stuName := ""
			famName := ""
			if len(pstu) == 1 {
				ad := pstu[0]["ad"]
				if ad == "a" {
					stuName = pstu[0]["name"]
					famName = pstu[0]["relation"]
					nickName = stuName + "的" + famName
					tmemeber[k]["relation"] = nickName
				} else {
					tmemeber[k]["relation"] = ""
				}
			} else {
				count := 0
				for s, stu := range pstu {
					count++
					if count > 2 {
						nickName = nickName[0:len(nickName)-1] + "...."
						break
					}
					if pstu[s]["ad"] == "a" {
						nickName += stu["name"] + "的" + pstu[s]["relation"] + ","
					}
					if pstu[s]["ad"] == "d" {
						count--
					}
				}
				log.Debug("nickName", fmt.Sprint(nickName))
				if nickName != "" {
					nickName = nickName[0 : len(nickName)-1]
					tmemeber[k]["relation"] = nickName
				} else {
					tmemeber[k]["relation"] = nickName
				}
			}
			data = append(data, tmemeber[k])
		}
		if len(fam) == 1 && len(member) == 1 {
			if len(member) == 0 {
				tmemeber[k]["relation"] = ""
			} else {
				relation := member[0]["relation"]
				tmemeber[k]["relation"] = relation
			}
			data = append(data, tmemeber[k])
		}
	}
	log.Debug("method_end", "FindTeacherMember", "status", "success")
	return data, nil, "100"
}

func (service createClassService) UpdateGroupChatInfo(gp *SchoolYear) (error, string) {
	log.Debug("method_start", "UpdateGroupChatInfo", "input", "=="+fmt.Sprint(gp))
	groupchat := &chat.GroupChatInfo{GroupChatId: gp.GroupChatId, GroupChatUrl: gp.GroupChatUrl, UserId: gp.UserId}
	response, err := chatClient.UpdateGroupChat(groupchat)
	if err != nil {
		return err, ""
	}
	if response.Err != "" {
		return errors.New(response.Err), ""
	}
	log.Debug("method_end", "UpdateGroupChatInfo", "status", "success")
	return nil, response.Code
}
func (service createClassService) UpdateTeachInfo(gp *SchoolYear) (error, string) {
	log.Debug("method_start", "UpdateGroupChatInfo", "input", "=="+fmt.Sprint(gp))
	res, err := s_teacher_classmember(gp)
	if err != nil {
		return err, ""
	}
	if len(res) == 0 {
		return nil, flag3
	}
	err = update_teacherMember(gp)
	if err != nil {
		return err, ""
	}
	//更新班级视图title
	st, err := s_teacher_classmember(gp)
	if err != nil {
		return err, ""
	}
	if len(st) == 0 {
		return nil, ""
	}
	status := st[0]["type"]
	if status == "0" {
		status = "1"
	} else {
		status = "3"
	}

	log.Debug("status", fmt.Sprint(status))
	myschool := &SchoolYear{ContentId: gp.ClId, Title: gp.Relation, UserId: gp.UserId, Type: status}
	log.Debug("myschool", fmt.Sprint(myschool))
	err = update_myschool(myschool)
	if err != nil {
		return err, ""
	}
	log.Debug("method_end", "UpdateGroupChatInfo", "status", "success")
	return nil, "100"
}

//check用户是否是班主任
func checkTeachManager(sy *SchoolYear) (error, string) {
	data, err := checkTeacher(mysqlClient, sy)
	if err != nil {
		return err, ""
	}
	//0代表
	if data["type"] != "0" || data == nil {
		return nil, flag3 //无权限操作
	}
	return nil, "100"
}

//check 班级clId是否存在
func checkClssId(sy *SchoolYear) (int, error) {
	data, err := checkClId(mysqlClient, sy)
	if err != nil {
		return data, err
	}
	if data == 0 {
		return data, err
	}
	return data, err
}

//获取学校邀请码和二维码
func InviteCode() (inviteCode string, inviteQRCode string, err error) {

	req, err := redisClient.Incr("mschool:clInvite")
	if err != nil {
		return
	}
	inviteCode = fmt.Sprint(100000 + req)
	inviteQRCode = "mschool/clInvite/" + inviteCode + ".png"
	go writePng(inviteCode)
	return
}

//生成邀请码二维码
func writePng(inviteCode string) (string, error) {
	filename := "mschool/clInvite/" + inviteCode + ".png"
	img, err := qr.Encode(InviteUrl+"?clInvite="+inviteCode, qr.L, qr.Unicode)
	if err != nil {
		return "", err
	}

	img, err = barcode.Scale(img, 300, 300)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	out_buf := buf.Bytes()
	err = ossuplodding(filename, out_buf)
	if err != nil {
		return "", err
	}
	return filename, nil
}

//上传oss
func ossuplodding(name string, data []byte) error {
	client, err := oss.New(oss_endpoint, key_id, key_secret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucket_name_mtalk)
	if err != nil {
		return err
	}

	err = bucket.PutObject(name, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
} //获取学年
func getSYear() string {
	//获取当前时间
	t := time.Now()
	//获取当前时间的年份
	yNow, _ := strconv.Atoi(t.Format("2006"))
	yearNow := t.Format("2006") + "-06-01 00:00:00"
	//获取当前年份的6月1号0时
	the_time, _ := time.Parse("2006-01-02 15:04:05", yearNow)
	syear := ""
	if t.After(the_time) {
		syear = t.Format("2006") + "-" + strconv.Itoa(yNow+1)
	} else {
		syear = strconv.Itoa(yNow-1) + "-" + t.Format("2006")
	}
	return syear
}
func pushMessage(p *Push) {
	list := []push.PushInfo{}
	pushInfo0 := push.PushInfo{
		Title:    p.Title,
		Text:     p.Text,
		JsonInfo: "{\"clId\":\"" + p.ClId + "\",\"stId\":\"" + p.StId + "\",\"haveRobot\":\"" + p.HaveRobot + "\",\"type\":\"" + p.Type + "\",\"name\":\"" + p.Name + "\"}",
		Alias:    getAlias(p.UserId), //U14921352095531169 //U146970017247836
	} //6685c1b99792c68fcef07029f3ded9c2 可用的

	list = append(list, pushInfo0)
	req := &push.PushRequest{
		PushType:  "0", //0代表推送消息，1代表透传消息,all_0代表给所有用户推送通知
		PushInfos: list,
	}
	pushClient.PushErrorLog(req)
}
func getAlias(userId string) string {
	//推送的别名设置，由于推送别名不能超过40字节。
	//推送别名为：UUID（去掉横杆）后32位 + UserId的后8位。
	//目前UUID正好是32位，但是不确定以后会不会变动，所以判断超过32位的话取后32位
	UUID, _ := redisClient.Hget("U:"+userId, "UUID")
	if len(UUID) > 32 {
		UUID = UUID[len(UUID)-32 : len(UUID)]
	}
	returnString := ""
	UUIDs := strings.Split(UUID, "-")
	for _, v := range UUIDs {
		returnString += v
	}
	lenUId := len(userId)
	returnString += userId[lenUId-8 : lenUId]
	return returnString
}

func getAccessToken() (string, error) {
	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + baidu_face_API_Key + "&client_secret=" + baidu_face_Secret_Key
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}
	return result["access_token"].(string), nil
}

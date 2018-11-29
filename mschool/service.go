package main

import (
	"fmt"
	"time"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type MschoolService interface {
	SearchSchool(*School) (map[string]interface{}, string, error)     //通过userId获取其创建的学校信息
	CreateSchoolYear(*School) (map[string]interface{}, string, error) //创建学年
	UpSchoolYear(*School) (string, error)                             //修改学年
	MySchool(*School) ([][]map[string]string, string, error)          //我的校园主页
	SetWorkDay(*School) (string, error)                               //设置上学日期
	LookWorkDay(*School) ([]map[string]interface{}, string, error)    //查看日历
	WorkDay(*School) ([]map[string]string, string, error)             //单天上学休假情况

	//定时收集需要学校端拉取的数据
	FaceDataGather() (string, error)
	LabelGuidDataGather() (string, error)
	GetDataFileUrl(*School) (map[string]string, string, error) //根据学校Id获取学校数据文件临时Url列表,暂时无用
	DelDataFile(*School) (string, error)                       //根据学校Id删除学校数据文件列表
}

type mschoolService struct{}

type School struct {
	ScId            string   `json:"scId"`            //学校主键
	SyId            string   `json:"syId"`            //学年主键
	ScInviteCode    string   `json:"scInviteCode"`    //学校邀请码
	ScInviteQRCode  string   `json:"scInviteQRCode"`  //学校邀请码二维码
	Address         string   `json:"address"`         //学校地址
	BeginTime       string   `json:"beginTime"`       //学年开学时间
	EndTime         string   `json:"endTime"`         //学年结束时间
	Grades          []string `json:"grades"`          //年级数组
	Grade           string   `json:"grade"`           //年级
	Location        string   `json:"location"`        //定位
	Nature          string   `json:"nature"`          //学校性质（如 私立，公立）
	Principal       string   `json:"principal"`       //负责人姓名
	PrincipalPhone  string   `json:"principalPhone"`  //负责人电话
	PrincipalUserId string   `json:"principalUserId"` //负责人userid
	ScName          string   `json:"scName"`          //学校名字
	Type            string   `json:"type"`            //学校类型，如，小学，初中，高中
	MsType          string   `json:"msType"`          //我的校园 0:学校负责人视图，1老师班级视图，2家长班级视图
	MsId            string   `json:"msId"`            //我的校园id
	Title           string   `json:"title"`           //视图上显示的标题（明航小学）
	ContentId       string   `json:"contentId"`       //内容id,学年id,老师班级id，家长班级id
	StudyYear       string   `json:"studyYear"`       //学年（2018-2019）
	GrId            string   `json:"grId"`            //年级id
	UserId          string   `json:"userId"`          //用户id
	//设置上学日期
	WdId         string              `json:"wdId"`         //主键
	ExceptionDay string              `json:"exceptionDay"` //异常日期
	Work         string              `json:"work"`         //0休息，1上学,2半休半上学
	Gradess      []map[string]string `json:"gradess"`      //传进来的 年级和是否上学 数组

	//设置入院时间字段
	WtId       string `json:"wtId"`       //主键
	BeginDay   string `json:"beginDay"`   //生效日期
	EffectWeek string `json:"effectWeek"` //作用到的星期（0-6逗号隔开）
	EndDay     string `json:"endDay"`     //失效日期
	IntoTime   string `json:"intoTime"`   //入园时间
	OutTime    string `json:"outTime"`    //出园时间

	//节假日表
	Year string `json:"year"` //年（2018）

	ClId string `json:"clId"` //班级Id
}

//创建学年
func (service mschoolService) CreateSchoolYear(school *School) (map[string]interface{}, string, error) { //结果，code err
	log.Debug("method_start", "CreateSchoolYear", "input", fmt.Sprint(school))
	returnMap := map[string]interface{}{}
	//获取id
	id := idGenClient.GetUniqueId()

	//查询此用户创建的学校
	scList, err := s_school(&School{PrincipalUserId: school.PrincipalUserId})
	if err != nil {
		return returnMap, "", err
	}
	if len(scList) < 1 { //如果没有创建过
		//获取学校邀请码
		school.ScInviteCode, school.ScInviteQRCode, err = InviteCode()
		if err != nil {
			return returnMap, "", err
		}
		school.ScId = "sc" + id
		//插入数据到school
		err = i_school(school)
		if err != nil {
			return returnMap, "", err
		}
	} else { //如果已经创建过
		school.ScId = scList[0]["scId"]
		school.ScName = scList[0]["scName"]
	}
	//判断该学校有没有重复创建学年
	school.StudyYear = now2schoolYear()
	syList, err := s_schoolyear(&School{ScId: school.ScId, StudyYear: school.StudyYear})
	if err != nil {
		return returnMap, "", err
	}
	if len(syList) > 0 {
		return returnMap, code1, nil //该学校本学年已创建过
	}

	school.SyId = "sy" + id
	school.MsId = "ms" + id

	//设置Title,ContentId
	school.Title = school.ScName
	school.ContentId = school.SyId
	school.MsType = "0"

	//插入数据到schoolyear
	err = i_schoolyear(school)
	if err != nil {
		return returnMap, "", err
	}
	//插入数据到grade
	gradeList := []map[string]string{}
	for i := 0; i < len(school.Grades); i++ {
		gradeMap := map[string]string{}
		school.Grade = school.Grades[i]
		school.GrId = "gr" + id + fmt.Sprint(i+1)
		gradeMap["grade"] = school.Grade
		gradeMap["grId"] = school.GrId
		gradeList = append(gradeList, gradeMap)
		err = i_grade(school)
		if err != nil {
			return returnMap, "", err
		}
	}
	//插入数据到myschool
	err = i_myschool(school)
	if err != nil {
		return returnMap, "", err
	}
	returnMap["gradess"] = gradeList
	returnMap["scId"] = school.ScId
	returnMap["studyYear"] = school.StudyYear
	returnMap["syId"] = school.SyId
	returnMap["scInviteCode"] = school.ScInviteCode
	returnMap["scInviteQRCode"] = school.ScInviteQRCode
	log.Debug("method_end", "CreateSchoolYear", "status", "success")
	return returnMap, "100", nil
}

//通过学年id查找学校学年
func (service mschoolService) SearchSchool(school *School) (map[string]interface{}, string, error) { //结果，code err
	log.Debug("method_start", "SearchSchool", "input", fmt.Sprint(school))
	returnMap := map[string]interface{}{}
	//查询学年
	syList, err := s_schoolyear(&School{SyId: school.SyId})
	if err != nil || len(syList) != 1 {
		return returnMap, "", errors.New("数据错误" + err.Error())
	}
	//查询学校
	scList, err := s_school(&School{ScId: syList[0]["scId"]})
	if err != nil || len(scList) != 1 {
		return returnMap, "", errors.New("数据错误" + err.Error())
	}
	//查询年级
	grList, err := s_grade(&School{SyId: school.SyId})
	if err != nil {
		return returnMap, "", err
	}

	returnMap["gradess"] = grList
	returnMap["school"] = scList[0]
	returnMap["schoolYear"] = syList[0]
	log.Debug("method_end", "SearchSchool", "status", "success")
	return returnMap, "100", nil
}

//通过userid查找我的校园
func (service mschoolService) MySchool(school *School) ([][]map[string]string, string, error) { //结果，code err
	log.Debug("method_start", "MySchool", "input", fmt.Sprint(school))
	returnList := [][]map[string]string{} //返回值
	syList, err := s_myschool(&School{UserId: school.UserId})
	if err != nil {
		return returnList, "", err
	}
	//遍历封装数据,一个学年封装成一个list,所有学年封装成一个大list
	studyYear := ""
	for _, syMap1 := range syList {
		if studyYear == syMap1["studyYear"] {
			continue
		}
		studyYear = syMap1["studyYear"]
		list := []map[string]string{}
		for _, syMap2 := range syList {
			if syMap1["studyYear"] == syMap2["studyYear"] {
				//插入是否有设备字段
				if syMap2["type"] == "0" { //学校负责人（contentId是学年id）
					syId := syMap2["contentId"]
					syList1, err := s_schoolyear(&School{SyId: syId})
					if err != nil {
						return returnList, "", err
					}
					if len(syList1) != 1 {
						return returnList, code3, nil
					}
					scId := syList1[0]["scId"]
					scList, err := s_school(&School{ScId: scId})
					if err != nil {
						return returnList, "", err
					}
					if len(scList) != 1 {
						return returnList, code3, nil
					}
					syMap2["gradeClass"] = "" //后端没用，前端解析需要
					syMap2["haveRobot"] = scList[0]["haveRobot"]
				} else { //家长或者班主任（contentId是班级id）
					clId := syMap2["contentId"]
					clList, err := s_class(&School{ClId: clId}) //查询班级取syId,clName,grId
					if err != nil {
						return returnList, "", err
					}
					if len(clList) != 1 {
						return returnList, code3, nil
					}
					grId := clList[0]["grId"]
					grList, err := s_grade(&School{GrId: grId}) //查年级名字
					if err != nil {
						return returnList, "", err
					}
					if len(grList) != 1 {
						return returnList, code3, nil
					}
					syId := clList[0]["syId"]
					syList1, err := s_schoolyear(&School{SyId: syId}) //查学年取scId，
					if err != nil {
						return returnList, "", err
					}
					if len(syList1) != 1 {
						return returnList, code3, nil
					}
					scId := syList1[0]["scId"]
					scList, err := s_school(&School{ScId: scId}) //查学校
					if err != nil {
						return returnList, "", err
					}
					if len(scList) != 1 {
						return returnList, code3, nil
					}
					syMap2["gradeClass"] = grList[0]["grade"] + clList[0]["clName"] //年级+班级名字
					syMap2["haveRobot"] = scList[0]["haveRobot"]
				}
				list = append(list, syMap2)
			}
		}
		returnList = append(returnList, list)
	}
	log.Debug("method_end", "MySchool", "status", "success")
	return returnList, "100", nil
}

//设置上学日期
func (service mschoolService) SetWorkDay(school *School) (string, error) { //结果，code err
	log.Debug("method_start", "SetWorkDay", "input", fmt.Sprint(school))
	//查看此人是否有权限操作此方法
	isAdmin, err := isAdmin(school)
	if err != nil || isAdmin == false { //无权操作
		return code2, err
	}

	//查询此学年这一天有没有设置
	wdList, err := s_workday(&School{SyId: school.SyId, ExceptionDay: school.ExceptionDay})
	if err != nil {
		return "", err
	}
	week := day2week(school.ExceptionDay) //获取这一天的周
	isWork, err := week2work(week)        //获取这一天的正常上学双休情况
	if err != nil {
		return "", err
	}

	//查看今天是否是法定安排
	hoList, err := s_holiday(&School{ExceptionDay: school.ExceptionDay})
	if err != nil {
		return "", err
	}
	status := ""
	if len(hoList) > 0 {
		status = hoList[0]["status"]
	}
	//查看各年级是否都设置得一样
	work := ""
	flag := 0
	for i, v := range school.Gradess {
		if i == 0 {
			work = v["work"]
		}

		if v["work"] != work { //如果传进来的参数中，所有年级有一个设置得不一样就都修改进数据库
			flag = 1
			break
		}
	}

	//如果这一天没有被设置
	if len(wdList) == 0 {
		inputFlag := 0 //是否存入数据库标识
		if flag == 1 { //如果传进来的参数中，所有年级有一个设置得不一样就都插入进数据库
			inputFlag = 1
		} else {
			if work != status { //与法定安排不一样，
				inputFlag = 1
			} else if status == "" && work != isWork { //法定假日没有设置，并且与正常周5，双休不一样
				inputFlag = 1
			}
		}
		if inputFlag == 1 {
			id := idGenClient.GetUniqueId()
			for i, _ := range school.Gradess {
				//赋予一个id,然后存到数据库
				school.Gradess[i]["wdId"] = "wd" + id + fmt.Sprint(i)
			}
			err := i_workday(school)
			if err != nil {
				return "", err
			}
		}
	} else { //如果这一天有设置

		if flag == 1 { //如果传进来的参数中，所有年级有一个设置得不一样就都修改进数据库
			for _, v := range school.Gradess {
				school.Work = v["work"]
				school.GrId = v["grId"]
				err := u_workday(school)
				if err != nil {
					return "", err
				}
			}
		} else { //如果所有年级都设置得一样，则判断设置得是否跟法定安排一样，是否与正常周5，双休一样
			if work == status { //与法定安排一样
				err := d_workday(school)
				if err != nil {
					return "", err
				}
			} else if status == "" && work == isWork { //没有设置法定假日并且与正常周5，双休一样   就删掉数据库中当天的数据
				err := d_workday(school)
				if err != nil {
					return "", err
				}
			} else { //否则修改数据库中这天的所有
				for _, v := range school.Gradess {
					school.Work = v["work"]
					school.GrId = v["grId"]
					err := u_workday(school)
					if err != nil {
						return "", err
					}
				}
			}
		}
	}
	//通过管道删除redis日期缓存worktime:WT:${syId}:${grId}
	err = DelCache(school.SyId)
	if err != nil {
		return "", err
	}
	log.Debug("method_end", "SetWorkDay", "status", "success")
	return "100", nil
}

//查看学年日历
func (service mschoolService) LookWorkDay(school *School) ([]map[string]interface{}, string, error) { //结果，code err
	log.Debug("method_start", "LookWorkDay", "input", fmt.Sprint(school))
	returnList := []map[string]interface{}{}

	//查询此学年所有设置过的日期
	wdList, err := s_workday(&School{SyId: school.SyId})
	if err != nil {
		return returnList, "", err
	}
	//查询今年所有节假日
	hoList, err := s_holiday(&School{Year: fmt.Sprint(time.Now().Year())})
	if err != nil {
		return returnList, "", err
	}

	month := ""
	for _, sdMap1 := range wdList {
		if month == day2month(sdMap1["exceptionDay"]) { //相同的月只遍历一次
			continue
		}
		monthMap := map[string]interface{}{}
		month = day2month(sdMap1["exceptionDay"])
		days := []map[string]string{}

		list := []map[string]string{} //单个月的数据数组
		for _, sdMap2 := range wdList { //把单个月的数据放进单个月的数组
			if month == day2month(sdMap2["exceptionDay"]) {
				list = append(list, sdMap2)
			}
		}

		//遍历单个月的数据
		day := ""
		for _, v1 := range list {
			if day == v1["exceptionDay"] { //相同的日期只遍历一次
				continue
			}
			day = v1["exceptionDay"]
			flag := 0
			for _, v2 := range list {
				if day == v2["exceptionDay"] { //如果是同一天
					if v1["work"] != v2["work"] { //如果年级设置得不一致
						flag = 1
						break
					}
				}
			}
			dayMap := map[string]string{}
			dayMap["d"] = day
			dayMap["f"] = "0" //全部定义为  0调整
			if flag == 1 {
				dayMap["w"] = "2" //一半上课，一半不上课
			} else {
				dayMap["w"] = v1["work"] //数据库中的值
			}
			days = append(days, dayMap) //一天的map封装进天的数组
		}

		monthMap["month"] = month
		monthMap["days"] = days
		returnList = append(returnList, monthMap)
	}

	//把缺失的节假日再补进去
	for _, hoMap := range hoList { //遍历节假日
		hoMonth := day2month(hoMap["hoDay"]) //月份
		monthFlag := 0
		dayFlag := 0
	a:
		for _, reL := range returnList { //遍历已设置的数据
			daysList := reL["days"].([]map[string]string)
			month := reL["month"].(string)
			if hoMonth == month { //节假日月份中有被设置
				monthFlag = 1
				for _, day := range daysList {
					if day["d"] == hoMap["hoDay"] { //节假日中此天被设置
						dayFlag = 1
						break a
					}
				}
			}
		}
		if monthFlag == 0 && dayFlag == 0 { //节假日这个月这个天都没有被负责人设置过,则把节假日连同月也添加进返回数组
			monthMap := map[string]interface{}{}
			days := []map[string]string{}
			dayMap := map[string]string{}
			dayMap["d"] = hoMap["hoDay"]  //日期等于节假日的日期
			dayMap["f"] = "1"             //全部定义为  1法定
			dayMap["w"] = hoMap["status"] //法定是否上课
			days = append(days, dayMap)   //map封装进天的数组
			monthMap["month"] = hoMonth
			monthMap["days"] = days
			returnList = append(returnList, monthMap)
		} else if monthFlag == 1 && dayFlag == 0 { //这个月中某一天被设置过，且与此节假日是不重复的天，则把此天插入到返回数组对应的月中
			for i, v := range returnList {
				daysList := v["days"].([]map[string]string)
				month := v["month"].(string)
				if hoMonth == month {
					dayMap := map[string]string{}
					dayMap["d"] = hoMap["hoDay"]        //日期等于节假日的日期
					dayMap["f"] = "1"                   //全部定义为  1法定
					dayMap["w"] = hoMap["status"]       //法定是否上课
					daysList = append(daysList, dayMap) //map封装进天的数组
					returnList[i]["days"] = daysList
				}
			}
		}

	}
	log.Debug("method_end", "LookWorkDay", "status", "success")
	return returnList, "100", nil
}

//单天上学放假情况  传入天，学年
func (service mschoolService) WorkDay(school *School) ([]map[string]string, string, error) { //结果，code err
	log.Debug("method_start", "WorkDay", "input", fmt.Sprint(school))
	returnList := []map[string]string{}
	//查询这天是否设置
	wdList, err := s_workday(&School{ExceptionDay: school.ExceptionDay, SyId: school.SyId})
	if err != nil {
		return returnList, code2, err
	}

	if len(wdList) > 0 { //如果这一天被设置过，则直接取这一天的数据
		for _, v := range wdList {
			wdMap := map[string]string{}
			wdMap["grade"] = v["grade"]
			wdMap["grId"] = v["grId"]
			wdMap["work"] = v["work"]
			returnList = append(returnList, wdMap)
		}
	} else {

		//查询这天是否是节假日安排
		hoList, err := s_holiday(&School{ExceptionDay: school.ExceptionDay})
		if err != nil {
			return returnList, code2, err
		}
		//查询此学年有多少年级
		grList, err := s_grade(&School{SyId: school.SyId})

		if len(hoList) > 0 { //如果这一天是节假日
			for _, v := range grList {
				wdMap := map[string]string{}
				wdMap["grade"] = v["grade"]
				wdMap["grId"] = v["grId"]
				wdMap["work"] = hoList[0]["status"]
				returnList = append(returnList, wdMap)
			}
		} else { //如果不是节假日，也没有设置
			for _, v := range grList {
				wdMap := map[string]string{}
				wdMap["grade"] = v["grade"]
				wdMap["grId"] = v["grId"]
				work, err := week2work(day2week(school.ExceptionDay))
				if err != nil {
					return returnList, code2, err
				}
				wdMap["work"] = work
				returnList = append(returnList, wdMap)
			}
		}

	}

	log.Debug("method_end", "WorkDay", "status", "success")
	return returnList, "100", nil
}

//修改学年信息
func (service mschoolService) UpSchoolYear(school *School) (string, error) { //结果，code err
	log.Debug("method_start", "UpSchoolYear", "input", fmt.Sprint(school))
	isAdmin, err := isAdmin(school)
	if err != nil || isAdmin == false { //无权操作
		return code2, err
	}
	if school.BeginTime == "" && school.EndTime == "" { //修改邀请码
		syList, err := s_schoolyear(&School{SyId: school.SyId})
		if err != nil {
			return "", err
		}
		if len(syList) != 1 {
			return code3, err
		}
		school.ScId = syList[0]["scId"]
		school.ScInviteCode, school.ScInviteQRCode, err = InviteCode() //生成新的邀请码
		if err != nil {
			return "", err
		}
		err = u_school(school)
		if err != nil {
			return "", err
		}
	} else { //修改学年日期
		err := u_schoolyear(school)
		if err != nil {
			return "", err
		}
		//通过管道删除redis日期缓存worktime:WT:${syId}:${grId}
		err = DelCache(school.SyId)
		if err != nil {
			return "", err
		}
	}
	log.Debug("method_end", "UpSchoolYear", "status", "success")
	return "100", nil
}

//收集需要学校端拉去的人脸数据
func (service mschoolService) FaceDataGather() (string, error) { //code err
	log.Debug("method_start", "FaceDataGather", "input", "")
	scList, err := s_school(&School{})
	if err != nil {
		return "", err
	}
	for _, v := range scList {
		syIdMap, err := s_syIdByScIdOneData(v["scId"]) //根据学校id查出此学校最新学年id
		if err != nil {
			return "", err
		}
		if len(syIdMap) < 1 { //代表此学校还没创建学年
			continue
		}
		clIdList, err := s_clIdBySyId(syIdMap["syId"]) //根据学年Id查询班级Id
		if err != nil {
			return "", err
		}
		if len(clIdList) < 1 { //代表此学年还没创建班级
			continue
		}
		stIdList, err := s_stIdByclIds(clIdList) //根据班级Id批量查询学生id (本学年此学校所有学生id)
		if err != nil {
			return "", err
		}
		if len(stIdList) < 1 { //代表这些班级还没加入学生
			continue
		}
		studentList, err := s_studentByFileAndstIds(stIdList) //查询这一批学生中人脸被更新过的数据
		if err != nil {
			return "", err
		}
		if len(studentList) < 1 { //代表这些学生没有更新过信息
			continue
		}

		err = studentListFaceOss(v["scId"], studentList) //把这个学校所有的学生的人脸图丢到oss上一个文件夹中
		if err != nil {
			return "", err
		}
		//修改数据库
		err = u_studentByFileAndstIds(studentList) //修改这一批学生中人脸被更新过的数据
		if err != nil {
			return "", err
		}

	}

	log.Debug("method_end", "FaceDataGather", "status", "success")
	return "100", nil
}

//收集需要学校端拉去的LabelGuid数据
func (service mschoolService) LabelGuidDataGather() (string, error) { //code err
	log.Debug("method_start", "LabelGuidDataGather", "input", "")
	scList, err := s_school(&School{})
	if err != nil {
		return "", err
	}
	for _, v := range scList {
		syIdMap, err := s_syIdByScIdOneData(v["scId"]) //根据学校id查出此学校最新学年id
		if err != nil {
			return "", err
		}
		if len(syIdMap) < 1 { //代表此学校还没创建学年
			continue
		}
		clIdList, err := s_clIdBySyId(syIdMap["syId"]) //根据学年Id查询班级Id
		if err != nil {
			return "", err
		}
		if len(clIdList) < 1 { //代表此学年还没创建班级
			continue
		}
		stIdList, err := s_stIdByclIds(clIdList) //根据班级Id批量查询学生id (本学年此学校所有学生id)
		if err != nil {
			return "", err
		}
		if len(stIdList) < 1 { //代表这些班级还没加入学生
			continue
		}
		guIdList, err := s_guIdByStIds(stIdList) //查询一批学生中的Guid没同步过的数据 按stid排序
		if err != nil {
			return "", err
		}
		if len(guIdList) < 1 { //代表这些学生还没新绑定定位贴
			continue
		}
		err = studentListLabelGuidOss(v["scId"], guIdList) //把这个学校所有的学生的人脸图和绑定的定位贴丢到oss上一个文件夹中
		if err != nil {
			return "", err
		}
		err = u_guIdByStIds(stIdList) //修改一批学生中的Guid没同步过的数据为已同步
		if err != nil {
			return "", err
		}
	}

	log.Debug("method_end", "LabelGuidDataGather", "status", "success")
	return "100", nil
}

////根据学校Id获取学校数据文件临时Url列表
func (service mschoolService) GetDataFileUrl(school *School) (map[string]string, string, error) { //code err
	log.Debug("method_start", "GetDataFileUrl", "input", fmt.Sprint(school))
	returnMap := map[string]string{}
	client, err := oss.New(oss_endpoint, key_id, key_secret) //初始化账号密码
	if err != nil {
		return nil, "", err
	}
	bucket, err := client.Bucket(bucket_name_face) //初始化bucket为qx-face
	if err != nil {
		return nil, "", err
	}
	faceFile := school.ScId + "/FaceData.txt"
	isHave, err := bucket.IsObjectExist(faceFile) //判断oss上目前是否存在此文件
	if err != nil {
		return nil, "", err
	}
	if isHave == true {
		fileUrl, err := bucket.SignURL(faceFile, oss.HTTPGet, 3600) //存在就获取他的临时url
		if err != nil {
			return nil, "", err
		}
		returnMap["faceFile"] = fileUrl
	}

	labelGuidFile := school.ScId + "/LabelGuidData.txt"
	isHave, err = bucket.IsObjectExist(labelGuidFile) //判断oss上目前是否存在此文件
	if err != nil {
		return nil, "", err
	}
	if isHave == true {
		fileUrl, err := bucket.SignURL(labelGuidFile, oss.HTTPGet, 3600) //存在就获取他的临时url
		if err != nil {
			return nil, "", err
		}
		returnMap["labelGuidFile"] = fileUrl
	}
	log.Debug("method_end", "GetDataFileUrl", "status", "success")
	return returnMap, "100", nil
}

//根据学校Id删除学校数据文件列表
func (service mschoolService) DelDataFile(school *School) (string, error) { //code err
	log.Debug("method_start", "DelDataFile", "input", fmt.Sprint(school))
	faceFile := school.ScId + "/FaceData.txt"
	labelGuidFile := school.ScId + "/LabelGuidData.txt"
	client, err := oss.New(oss_endpoint, key_id, key_secret) //初始化账号密码
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(bucket_name_face) //初始化bucket为qx-face
	if err != nil {
		return "", err
	}
	//删除文件
	err = bucket.DeleteObject(faceFile)
	if err != nil {
		return "", err
	}
	err = bucket.DeleteObject(labelGuidFile)
	if err != nil {
		return "", err
	}
	log.Debug("method_end", "DelDataFile", "status", "success")
	return "100", nil
}

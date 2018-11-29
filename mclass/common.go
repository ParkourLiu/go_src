package main

import (
	"fmt"
	"time"
	"strconv"
	"errors"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//当前时间换算成学年 返回2018-2019
func now2schoolYear() string {
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

//日期转星期,传入2018-09-08返回（0-6）
func day2week(datetime string) string {
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation("2006-01-02", datetime, loc) //使用模板在对应时区转化为time.time类型
	return fmt.Sprint(int(theTime.Weekday()))
}

//判断星期是周一到周五工作日还是周6周7双休日，传入0-6 返回 0-1
func week2work(week string) (string, error) {
	if week == "1" || week == "2" || week == "3" || week == "4" || week == "5" {
		return "1", nil
	} else if week == "6" || week == "0" {
		return "0", nil
	} else {
		return "", errors.New("星期参数错误" + week)
	}
}

//2018-09-30  变成 2018-09
func day2month(day string) string {
	month := string([]rune(day)[:7])
	return month
}

//班级是否有此学生,有的话返回此学生id
func classHaveStudent(clId string, studentNum string, name string, sex string, birthday string) (string, bool, error) {
	haveSt := false
	sclList, err := s_student_class(&SchoolYear{ClId: clId, StudentNum: studentNum})
	if err != nil {
		return "", haveSt, err
	}
	for _, sclMap := range sclList {
		stList, err := s_student(&SchoolYear{StId: sclMap["stId"]})
		if err != nil {
			return "", haveSt, err
		}
		if len(stList) < 1 {
			continue
		}
		if name == stList[0]["name"] && sex == stList[0]["sex"] && birthday == stList[0]["birthday"] { //判断传进来的孩子是一致
			haveSt = true
			return sclMap["stId"], haveSt, nil
		}
	}
	return "", haveSt, nil
}

//班级是否有此家长
func classHaveFamily(clId string, userId string) (bool, error) {
	fcList, err := s_family_classmember(&SchoolYear{ClId: clId, UserId: userId})
	if err != nil {
		return false, err
	}
	if len(fcList) != 1 {
		return false, nil
	}
	return true, nil
}

//家长是否有此学生,有的话返回学生id，（此方法由于产品设计缺陷，重新加入临时处理机制）
func familyHaveStudent(userId string, name string, sex string, birthday string, clId string, studentNum string, relation string) (string, bool, error) {
	stId := ""
	psId := ""
	chsStId, chs, err := classHaveStudent(clId, studentNum, name, sex, birthday) //查询出此班级和此学生信息相同的学生Id
	if err != nil {
		return stId, false, err
	}
	psList, err := s_patriarch_student(&SchoolYear{UserId: userId}) //得到学生Id 数组
	if err != nil {
		return stId, false, err
	}

	for _, psMap := range psList { //遍历此家长的所有学生Id
		stList, err := s_student(&SchoolYear{StId: psMap["stId"]}) //通过id查询此学生信息
		if err != nil {
			return stId, false, err
		}
		if len(stList) != 1 {
			continue
		}
		if name == stList[0]["name"] && sex == stList[0]["sex"] && birthday == stList[0]["birthday"] { //判断传进来的孩子跟此家长已有的孩子是一致
			stId = psMap["stId"]
			psId = psMap["psId"]
			break
		}
	}

	if chs == true { //此班级有此学生
			if psId != "" { //此家长有此学生
				if chsStId != stId { //此班级的学生和此家长的学生不是同一个学生
					err := u_patriarch_student(&SchoolYear{StId: chsStId, PsId: psId, Relation: relation}) //修改家长与学生的关系指向班级中已有的学生
					if err != nil {
						return stId, false, err
					}
				}
			} else { //此家长没有此学生,新建一个与此学生的关系
				psId = "ps" + idGenClient.GetUniqueId()
				err := i_patriarch_student(&SchoolYear{PsId: psId, UserId: userId, StId: chsStId, Relation: relation})
				if err != nil {
					return stId, false, err
				}
			}
			return chsStId, true, err
		} else { //此班级没有此学生
			if psId != "" { //此家长有此学生
				err := u_patriarch_student(&SchoolYear{PsId: psId, Relation: relation}) //修改家长称呼为最新称呼
				if err != nil {
					return stId, false, err
				}
				return stId, true, err
			}
	}
	return stId, false, nil
}

//ctId := "CT" + idGenClient.GetUniqueId()
//ctMap := map[string]string{}
//ctMap["phoneNo"] = user.PhoneNo
//ctMap["trueName"] = "钦家用户"
//ctMap["isDef"] = "1" //设置成默认联系人
//err := redisClient.Hmset("CTR:"+ctId, ctMap)
//if err != nil {
//return err
//}
//return redisClient.LpushFromHead("U2CTR:"+user.UserId, ctId)

func createObjecte(sy *SchoolYear) (string, error) {
	scIdList, err := redisClient.LrangeAll("U2CTR:" + sy.UserId) //取出此用户所有紧急联系人
	if err != nil {
		return "", err
	}
	contactorId1 := ""
	if len(scIdList) < 1 { //此用户没有紧急联系人则创建
		contactorId1 = "CT" + idGenClient.GetUniqueId()
		phoneNo, err := redisClient.Hget("U:"+sy.UserId, "phoneNo") //获取此人手机号
		if err != nil {
			return "", err
		}
		go CTR(contactorId1, phoneNo, sy.UserId)
	} else {
		contactorId1 = scIdList[0]
	}
	url := mgtalk2Url + "Label/golangCreateObject.do"
	postStr := "{\"Head\": {},\"Para\": {\"trueName\":\"" + sy.Name + "\",\"userId\":\"" + sy.UserId + "\",\"stId\":\"" + sy.StId + "\",\"contactorId1\":\"" + contactorId1 + "\",\"birthDay\":\"" + sy.Birthday + "\",\"sex\":\"" + sy.Sex + "\"}}"
	req, err := DoBytesPost(url, postStr) //调用java生成对象的接口
	log.Debug("req", fmt.Sprint(req))
	if err != nil {
		return "", err
	}
	reqCode := req["code"].(string)
	if reqCode != "100" {
		return "JAVA" + reqCode, nil
	}
	return "", nil
}

//新增紧急联系人
func CTR(ctId string, phoneNo string, userId string) error {
	ctMap := map[string]string{}
	ctMap["phoneNo"] = phoneNo
	ctMap["trueName"] = "钦家用户"
	ctMap["isDef"] = "1" //设置成默认联系人
	err := redisClient.Hmset("CTR:"+ctId, ctMap)
	if err != nil {
		return err
	}
	return redisClient.LpushFromHead("U2CTR:"+userId, ctId)
}

//body提交二进制数据
func DoBytesPost(url string, postStr string) (map[string]interface{}, error) {
	data := []byte(postStr)
	fmt.Println(string(data))
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if resp != nil {
		defer resp.Body.Close() // ok，即使不读取Body中的数据，即使Body是空的，也要调用close方法,一定要判断resp是否为空
	}
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reqMap, err := jsonMap(string(b))
	if err != nil {
		return nil, err
	}
	return reqMap, err
}

func jsonMap(str string) (map[string]interface{}, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(str), &dat); err == nil {
		return dat, nil
	} else {
		return make(map[string]interface{}), err
	}
}

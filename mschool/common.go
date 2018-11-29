package main

import (
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"fmt"
	"bytes"
	"image/jpeg"
	"time"
	"strconv"
	"errors"
	"io/ioutil"
)

//获取学校邀请码和二维码
func InviteCode() (inviteCode string, inviteQRCode string, err error) {

	req, err := redisClient.Incr("mschool:scInvite")
	if err != nil {
		return
	}
	inviteCode = fmt.Sprint(100000 + req)
	inviteQRCode = "mschool/scInvite/" + inviteCode + ".png"
	go writePng(inviteCode)
	return
}

//生成邀请码二维码
func writePng(inviteCode string) (string, error) {
	filename := "mschool/scInvite/" + inviteCode + ".png"
	img, err := qr.Encode(InviteUrl+"?scInvite="+inviteCode, qr.L, qr.Unicode)
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

//上传oss到mtalk2
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
}

//当前时间换算成学年
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

//学校创建人判断
func isAdmin(school *School) (bool, error) {
	scList, err := s_school(&School{PrincipalUserId: school.UserId})
	if err != nil {
		return false, err
	}

	if len(scList) == 1 {
		school.StudyYear = now2schoolYear()
		school.ScId = scList[0]["scId"]
		syList, err := s_schoolyear(&School{ScId: school.ScId, StudyYear: school.StudyYear})
		if err != nil {
			return false, err
		}
		if len(syList) == 1 {
			if syList[0]["syId"] == school.SyId {
				return true, nil
			}
		}
	}
	return false, nil
}

//2018-09-30  变成 2018-09
func day2month(day string) string {
	month := string([]rune(day)[:7])
	return month
}

//一个学校的学生整理到一个Oss文件夹 ${scId}/FaceData.txt
func studentListFaceOss(scId string, studentList []map[string]string) error {
	FaceFile := scId + "/FaceData.txt"
	client, err := oss.New(oss_endpoint, key_id, key_secret) //初始化账号密码
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucket_name_face) //初始化bucket为qx-face
	if err != nil {
		return err
	}
	isHave, err := bucket.IsObjectExist(FaceFile) //判断oss上目前是否存在此文件
	if err != nil {
		return err
	}
	dataBuffer := bytes.Buffer{} //创建文本文件字符流
	if isHave == true { //oss上有此文件
		body, err := bucket.GetObject(FaceFile) //下载此文件,返回流文件
		if body != nil {
			defer body.Close() //如果不为空的话最终需要关闭掉此流
		}
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		dataBuffer.Write(data) //添加到总流文件
	}
	//遍历新的需要传输的学生人脸图并加入到总文件流
	for _, studentMap := range studentList {
		dataBuffer.WriteString(studentMap["stId"])     //加入学生id
		dataBuffer.WriteString(",")                    //逗号隔开
		dataBuffer.WriteString(studentMap["facePath"]) //人脸图片文件路径
		dataBuffer.WriteString("\r\n")                 //换行
	}

	dataBytes := dataBuffer.Bytes()
	return bucket.PutObject(FaceFile, bytes.NewReader(dataBytes)) //上传文件
}

//一个学校的学生整理到一个Oss文件夹 ${scId}/LabelGuidData.txt
func studentListLabelGuidOss(scId string, guIdList []map[string]string) error {
	fileName := scId + "/LabelGuidData.txt"
	client, err := oss.New(oss_endpoint, key_id, key_secret) //初始化账号密码
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucket_name_face) //初始化bucket为qx-face
	if err != nil {
		return err
	}
	isHave, err := bucket.IsObjectExist(fileName) //判断oss上目前是否存在此文件
	if err != nil {
		return err
	}
	dataBuffer := bytes.Buffer{} //创建文本文件字符流
	if isHave == true { //oss上有此文件
		body, err := bucket.GetObject(fileName) //下载此文件,返回流文件
		if body != nil {
			defer body.Close() //如果不为空的话最终需要关闭掉此流
		}
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		dataBuffer.Write(data) //添加到总流文件
	}
	//遍历新的需要传输的学生人脸图并加入到总文件流
	for _, studentMap := range guIdList {
		dataBuffer.WriteString(studentMap["stId"])      //加入学生id
		dataBuffer.WriteString(",")                     //逗号隔开
		dataBuffer.WriteString(studentMap["labelGuid"]) //LabelGuid
		dataBuffer.WriteString("\r\n")                  //换行
	}

	dataBytes := dataBuffer.Bytes()
	return bucket.PutObject(fileName, bytes.NewReader(dataBytes)) //上传文件
}

//通过管道删除redis日期缓存worktime:WT:${syId}:${grId}
func DelCache(syId string) error {
	keys, err := redisClient.Keys("worktime:WT:" + syId + "*")
	if err != nil {
		return err
	}
	c, err := redisClient.GetPipeline()
	if err != nil {
		return err
	}
	for _, v := range keys {
		redisClient.PipelineDel(c, v)
	}
	return redisClient.ExecutePipeline(c)

}



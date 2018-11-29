package main

import (
	//"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"net/http"

	"github.com/golang/go/src/pkg/errors"
)

type AutismService interface {
	StarDetails(*Autism) (map[string]interface{}, string, error) //星星详情
	SaveComment(*Autism) (string, error)

	StarList(*Autism) (string, string, []map[string]string, []map[string]string, error) //获取首页的数据
	Likes(*Autism) error                                                                //点赞
	GetUnionid(*Autism) (Un_ionid, error)                                               //获取unionid
}

type autismService struct{}

type Autism struct {
	St_id        string `json:"st_id"`
	St_name      string `json:"st_name"`
	St_head      string `json:"st_head"`
	St_type      string `json:"st_type"`
	St_content   string `json:"st_content"`
	St_likeCount string `json:"st_likeCount"`
	ClickUserId  string `json:"clickUserId"`
	LastCoid     string `json:"lastCoid"`
	Comment      string `json:"comment"`
	User_id      string `json:"user_id"`
	Code         string `json:"code"`
}

//单个星星详情，包括评论，点赞,评论，捐款,传入lastcoid,st_id,clickuserid
func (service autismService) StarDetails(autism *Autism) (map[string]interface{}, string, error) { //返回值，code,error
	returnMap := map[string]interface{}{}

	//查询出星星详情
	star, err := s_starById(mysqlClient, autism)
	if err != nil {
		return nil, "", err
	}
	if len(star) < 1 {
		return nil, flag1, nil
	}
	likeCount, err0 := redisClient.Slen(serviceName + ":likeFlag:" + autism.St_id)
	if err0 != nil {
		return nil, "", err0
	}

	likeCount0, _ := strconv.Atoi(star["st_likeCount"]) //假数据点赞数
	fmt.Println("likeCount", likeCount)
	fmt.Println("likeCount0", likeCount0)
	likeCount0 = int(likeCount) + likeCount0
	star["st_likeCount"] = fmt.Sprint(likeCount0) //查询出点赞数

	//查询评论
	commentList, err1 := s_commentsByStid(mysqlClient, autism)
	if err1 != nil {
		return nil, "", err1
	}
	for i := 0; i < len(commentList); i++ {
		photo, err2 := redisClient.Hget("U:"+commentList[i]["userId"], "imageName")
		if err2 != nil {
			return nil, "", err2
		}
		commentList[i]["photo"] = photo
	}

	//查询捐款
	donation, err3 := s_donation(mysqlClient)
	if err3 != nil {
		return nil, "", err3
	}
	//查询出评论数
	commentCount, err4 := s_commentCount(mysqlClient, autism)
	if err4 != nil {
		return nil, "", err4
	}
	//是否点赞
	likeFlag, err5 := redisClient.Sismember(serviceName+":likeFlag:"+autism.St_id, autism.ClickUserId) //是否点赞
	if err5 != nil {
		return nil, "", err5
	}

	//点亮数
	brightNum, err6 := redisClient.ZRANK(serviceName+":brightCount1", autism.ClickUserId)
	if err6 != nil {
		return nil, "", err6
	}
	if brightNum <= 0 { //存储点亮
		timeNow := time.Now().Unix()
		redisClient.Zadd(serviceName+":brightCount1", timeNow, autism.ClickUserId)
		brightCount0, _ := redisClient.Get(serviceName + ":brightCount0")
		brightCount0Int, _ := strconv.Atoi(brightCount0)
		brightCount1, _ := redisClient.Zlen(serviceName + ":brightCount1")
		redisClient.Set(serviceName+":brightNum:"+autism.ClickUserId, fmt.Sprint(int(brightCount1)+brightCount0Int))
	}

	//此人是第多少个点亮的
	bright, _ := redisClient.Get(serviceName + ":brightNum:" + autism.ClickUserId) //此人是第多少个点亮的

	if likeFlag {
		star["likeFlag"] = "1"
	} else {
		star["likeFlag"] = "0"
	}
	returnMap["star"] = star
	returnMap["commentList"] = commentList
	returnMap["donation"] = donation
	returnMap["commentCount"] = commentCount
	returnMap["brightNum"] = bright

	return returnMap, "100", nil
}

//传入userid,stid,comment
func (service autismService) SaveComment(autism *Autism) (string, error) { //返回值，code,error
	co_id := idGenClient.GetUniqueId()
	autism.LastCoid = co_id
	err := i_comment(mysqlClient, autism)

	return "100", err
}

//=====================================================================================================================
func (service autismService) StarList(autism *Autism) (string, string, []map[string]string, []map[string]string, error) {
	log.Debug("method_start", "StarList")
	dao := &AutismDao{}
	brightCount, likeCount, rolls, starList, err := dao.StarList(redisClient, mysqlClient, autism)
	if err != nil {
		log.Debug("method_end", "StarList", "status", "error")
		return brightCount, likeCount, rolls, starList, err
	}
	log.Debug("method_end", "StarList", "status", "success")
	return brightCount, likeCount, rolls, starList, err
}
func (service autismService) Likes(autism *Autism) error {
	log.Debug("method_start", "Likes")
	dao := &AutismDao{}
	err := dao.Likes(redisClient, mysqlClient, autism)
	if err != nil {
		return err
	}
	log.Debug("method_end", "Likes", "status", "success")
	return nil
}
func (service autismService) GetUnionid(autism *Autism) (Un_ionid, error) {
	log.Debug("method_start", "GetUnionid")
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid + "&secret=" + secret + "&code=" + autism.Code + "&grant_type=authorization_code"

	data, err := http.Get(url)

	if err != nil {
		return Un_ionid{}, err
	}
	defer data.Body.Close()
	getRes, gerErr := ioutil.ReadAll(data.Body)
	fmt.Println(string(getRes))
	if gerErr != nil {
		return Un_ionid{}, gerErr
	}
	resp := json2AccessToken(string(getRes))
	if len(resp.Errmsg) != 0 {
		return Un_ionid{}, errors.New(resp.Errmsg)
	}
	data2, err2 := http.Get("https://api.weixin.qq.com/sns/userinfo?access_token=" + resp.Access_token + "&openid=" + resp.Openid + "&lang=zh_CN")
	if err2 != nil {
		return Un_ionid{}, err2
	}
	defer data2.Body.Close()
	getRes2, gerErr2 := ioutil.ReadAll(data2.Body)
	fmt.Println(string(getRes2))
	if gerErr2 != nil {
		return Un_ionid{}, gerErr2
	}
	unionid := json2Unionid(string(getRes2))
	if len(unionid.Errmsg) != 0 {
		return Un_ionid{}, errors.New(unionid.Errmsg)
	}
	fmt.Println(unionid.Unionid)
	if 0 == unionid.Sex {
		unionid.Gender = "女"
	} else {
		unionid.Gender = "男"
	}

	//存储用户的union_id
	redisClient.Sismember(serviceName+":name:"+unionid.Unionid, autism.User_id)
	log.Debug("method_end", "GetUnionid", "status", "success")
	return unionid, nil
}
func json2AccessToken(str string) Access_Token {
	var token Access_Token
	json.Unmarshal([]byte(str), &token)
	return token
}
func json2Unionid(str string) Un_ionid {
	var unionid Un_ionid
	json.Unmarshal([]byte(str), &unionid)
	return unionid
}

type Access_Token struct {
	Access_token  string `json:"access_token"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
	Openid        string `json:"openid"`
	Scope         string `json:"scope"`
	Errcode       int    `json:"errcode"`
	Errmsg        string `json:"errmsg"`
}

type Un_ionid struct {
	Unionid    string   `json:"unionid"`
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Errcode    int      `json:"errcode"`
	Errmsg     string   `json:"errmsg"`
	Sex        int      `json:"sex"`
	Gender     string   `json:"gender"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
}

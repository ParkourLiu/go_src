package main

import (
	"encoding/json"
	"errors"
	"fmt"
	grpdynamic "grpdynamic/client"
	mactivity "mactivity/client"
	"math/rand"
	mgrpmgr "mgrpmgr/client"
	push "push/client"
	"strconv"
	"strings"
	"time"
	user "user/client"
)

type Recommend interface {
	PopuserUsers(*PageInfo) ([]map[string]string, string, error)            //活跃用户推荐
	RmduserUsers() ([]map[string]string, string, error)                     //分享中用户推荐
	FriendRecommend(*PageInfo) ([]map[string]interface{}, string, error)    //添加好友页面的推荐用户和动态（周动态发布最多的用户）
	HotDynamic(*PageInfo) ([]map[string]interface{}, string, string, error) //大家的动态，最热推荐
	NewDynamic(*PageInfo) ([]map[string]interface{}, string, string, error) //大家的动态，最新推荐

	ReverseHotDynamic(*PageInfo) ([]map[string]interface{}, string, string, error) //倒叙分页大家的动态，最热推荐
	ReverseNewDynamic(*PageInfo) ([]map[string]interface{}, string, string, error) //倒叙分页大家的动态，最新推荐
	SearchRecommend() ([]string, string, error)                                    //搜索推荐

	//存数据
	PushFavour() (string, error)          //给今天的动态推送通知
	PushStartActivity() (string, error)   //给明天要开始的活动参与者推送通知
	SavePopuserUsers() (string, error)    //存储活跃用户（周被赞最多）
	SaveRmduserUsers() (string, error)    //存储分享中用户(周发布个人动态+被赞数最多)
	SaveHotDynamic() (string, error)      //存储最热动态（赞最多）
	SaveFriendRecommend() (string, error) //存储推荐好友（周发布动态最多）
	SaveHomePageCache() (string, error)   //首页缓存
	AddDynamicFans() (string, error)      //增加动态的点赞数
}

type recommend struct{}

type PageInfo struct {
	LastId     string `json:"lastId"`
	PageSize   string `json:"pageSize"`
	LookUserId string `json:"lookUserId"`
}

//获取一周前的时间戳
func WeekTime() (int64, error) {
	log.Debug("████ ", fmt.Sprint("----------------------------------------"))
	now := time.Now()
	weekday := now.AddDate(0, 0, -7)
	loc, err1 := time.LoadLocation("Local")
	if err1 != nil {
		return 0, err1
	}
	theTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", weekday.Format("2006-01-02 15:04:05"), loc)
	if err2 != nil {
		return 0, err2
	}
	weekTime := theTime.Unix()
	return weekTime, nil
}

//获取制定天前或者后，传负数就是多少天前，正数就是多少天后的时间戳
func YesDayTime(days int) (int64, error) {
	now := time.Now()
	weekday := now.AddDate(0, 0, days)
	loc, err1 := time.LoadLocation("Local")
	if err1 != nil {
		return 0, err1
	}
	theTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", weekday.Format("2006-01-02 15:04:05"), loc)
	if err2 != nil {
		return 0, err2
	}
	yesDayTime := theTime.Unix()
	return yesDayTime, nil
}

func DifferenceTime(tm time.Duration) (int64, error) {
	now := time.Now()
	dTm := now.Add(tm)
	loc, err1 := time.LoadLocation("Local")
	if err1 != nil {
		return 0, err1
	}
	theTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", dTm.Format("2006-01-02 15:04:05"), loc)
	if err2 != nil {
		return 0, err2
	}
	DifferenceTime := theTime.Unix()
	return DifferenceTime, nil
}

//存储活跃用户
func (service recommend) SavePopuserUsers() (string, error) {
	log.Debug("method_start", "SavePopuserUsers", "input", "")
	allLikes, err := redisClient.Keys("grpdynamic:USERALLLIKE:*")
	if err != nil {
		return "", err
	}
	weekTime, err0 := WeekTime() //获取上周时间戳
	log.Debug("████", fmt.Sprint(err0))
	if err0 != nil {
		return "", err0
	}
	params := [][]interface{}{}
	for _, v := range allLikes {
		param := []interface{}{}
		userId := strings.Split(v, ":")[2]
		if userId != "" {
			fmt.Println("████userId", userId)
			weekLikeList, err1 := redisClient.ZRANGEBYSCORE(v, weekTime) //所有赞我的
			if err1 != nil {
				return "", err1
			}
			weekLikeCount := len(weekLikeList) //所有赞我的

			log.Debug("████weekLikeCount", fmt.Sprint(weekLikeCount))
			//封装数据，通过redis管道一次性插入
			param = append(param, "temporary_SavePopuserUsers")
			param = append(param, weekLikeCount)
			param = append(param, userId)
			params = append(params, param)
		}
	}

	err4 := redisClient.Del("temporary_SavePopuserUsers") //插入前先删除临时表
	if err4 != nil {
		return "", err4
	}
	_, err2 := redisClient.PipelineZadd(params) //临时排序表
	if err2 != nil {
		return "", err2
	}
	temporarylist, err3 := redisClient.Zrange("temporary_SavePopuserUsers", 0, 9)
	if err3 != nil {
		return "", err3
	}
	err4 = redisClient.Del("popuser:uids")
	log.Debug("████", fmt.Sprint("删除数据"))
	if err4 != nil {
		return "", err4
	}
	log.Debug("████", fmt.Sprint("444"))
	for _, v := range temporarylist {
		err5 := redisClient.LpushFromTail("popuser:uids", v)
		if err5 != nil {
			return "", err5
		}
	}
	redisClient.Del("temporary_SavePopuserUsers") //操作后删除临时表
	log.Debug("method_end", "SavePopuserUsers", "status", "success")
	return "100", nil
}

//存储分享中推荐用户
func (service recommend) SaveRmduserUsers() (string, error) {
	log.Debug("method_start", "SaveRmduserUsers", "input", "")
	weekTime, err0 := WeekTime() //获取上周时间戳
	if err0 != nil {
		return "", err0
	}
	weekNewDynamic, err := redisClient.ZRANGEBYSCORE("newDynamic", weekTime) //这一周发布的所有动态
	log.Debug("████weekNewDynamic", fmt.Sprint(weekNewDynamic))
	if err != nil {
		return "", err
	}

	err4 := redisClient.Del("temporary_SaveRmduserUsers") //插入前先删除临时排序表
	if err4 != nil {
		return "", err4
	}
	for _, v := range weekNewDynamic {
		userId, err1 := redisClient.Get("grpdynamic:D2U:" + v)
		log.Debug("████err1", fmt.Sprint(err1))
		if err1 != nil {
			if "redigo: nil returned" != err1.Error() {
				return "", err1
			}
		}
		if userId == "" {
			continue
		}
		err2 := redisClient.ZincrScore("temporary_SaveRmduserUsers", 1, userId)
		if err2 != nil {
			return "", err2
		}
	}
	log.Debug("████", fmt.Sprint("222222222222222222222222222222"))
	temporarylist, err3 := redisClient.Zrange("temporary_SaveRmduserUsers", 0, 9)
	if err3 != nil {
		return "", err3
	}
	log.Debug("████", fmt.Sprint("222222222222222222222222222222"))
	err4 = redisClient.Del("rmduser:uids")
	log.Debug("████", fmt.Sprint("删除数据"))
	if err4 != nil {
		return "", err4
	}
	for _, v := range temporarylist {
		err5 := redisClient.LpushFromTail("rmduser:uids", v)
		if err5 != nil {
			return "", err5
		}
	}
	redisClient.Del("temporary_SaveRmduserUsers") //操作后删除临时表
	log.Debug("method_end", "SaveRmduserUsers", "status", "success")
	return "100", nil
}

//最热动态存储
func (service recommend) SaveHotDynamic() (string, error) {
	log.Debug("method_start", "SaveHotDynamic", "input", "")
	weekTime, err := WeekTime()
	if err != nil {
		return "", err
	}
	dynamicIds, err := redisClient.ZRANGEBYSCORE("newDynamic", weekTime)
	for _, v := range dynamicIds {
		likeNum, err1 := redisClient.Zlen("grpdynamic:UL:" + v) //哪些人攒过我的动态
		fmt.Println("likeNum=============", likeNum)
		if err1 != nil {
			if err1.Error() != "redigo: nil returned" {
				return "", err1
			}
		}
		commentIds, err2 := redisClient.LrangeAll("grpdynamic:D2C:" + v) //动态的一级评论，不包含二级回复
		if err2 != nil {
			return "", err2
		}
		commentNum := int64(len(commentIds)) //此动态有多少一级评论
		fmt.Println("commentNum=============", commentNum)
		for _, va := range commentIds {
			replyNum, err3 := redisClient.Llen("grpdynamic:C2R:" + va) ////动态的二级回复
			fmt.Println("replyNum=============", replyNum)
			if err3 != nil {
				if err3.Error() != "redigo: nil returned" {
					return "", err3
				}
			}
			commentNum = commentNum + replyNum
		}
		_, err4 := redisClient.Zadd("hotdynamic:dids", likeNum+commentNum, v)
		if err4 != nil {
			return "", err4
		}
	}
	pushIds, err5 := redisClient.ZrangeAll("hotdynamic:push")
	if err5 != nil {
		return "", err5
	}
	nowTime := time.Now().Unix()
	if len(pushIds) > 0 {
		for _, v := range pushIds {
			_, err6 := redisClient.Zadd("hotdynamic:dids", int64(nowTime), v)
			if err6 != nil {
				return "", err6
			}
		}
	}
	log.Debug("method_end", "SaveHotDynamic", "status", "success")
	return "100", nil
}

//活跃用户推荐
func (service recommend) PopuserUsers(pageInfo *PageInfo) ([]map[string]string, string, error) { //返回值，code，err
	log.Debug("method_start", "PopuserUsers", "input", fmt.Sprint(pageInfo))
	request, err := redisClient.Lrange("popuser:uids", 0, 5)
	if err != nil {
		return nil, "", err
	}
	returnList := []map[string]string{}
	for _, v := range request {
		mysqlUser, err0 := userClient.SearchUser(&user.UserRequest{UserId: v})
		if err0 != nil {
			return nil, "", err0
		}
		if mysqlUser.Err != "" {
			return nil, "", errors.New(mysqlUser.Err)
		}
		requestUser := map[string]string{}
		requestUser["userId"] = mysqlUser.User.UserId
		requestUser["imageName"] = mysqlUser.User.ImageName
		requestUser["nickName"] = mysqlUser.User.NickName
		//查看此人发了多少动态
		dyList, err2 := redisClient.LrangeAll("U2D:" + v)
		if err2 != nil {
			return nil, "", err2
		}
		likeCount := int64(0) //总喜欢数
		for _, vdy := range dyList {
			likes, err3 := redisClient.Zlen("UL:" + vdy) //查看此动态被点了多少赞
			if err3 != nil {
				if "redigo: nil returned" != err3.Error() {
					return nil, "", err3
				}
			}
			likeCount += likes
		}

		//查看此人粉丝数
		fansCount, err4 := redisClient.Zlen("grpfans:uIdFansList:" + v)
		if err4 != nil {
			if "redigo: nil returned" != err4.Error() {
				return nil, "", err4
			}
		}

		//查看我是否关注此人
		score, err5 := redisClient.Zscore("grpfans:uIdFansList:"+v, pageInfo.LookUserId)
		if err5 != nil {
			if "redigo: nil returned" != err5.Error() {
				return nil, "", err5
			}

		}
		likeFlag := "0"
		if score != 0 {
			likeFlag = "1"
		}
		requestUser["likeCount"] = fmt.Sprint(likeCount)
		requestUser["fansCount"] = fmt.Sprint(fansCount)
		requestUser["likeFlag"] = fmt.Sprint(likeFlag)
		returnList = append(returnList, requestUser)
	}
	log.Debug("method_end", "PopuserUsers", "status", "success")
	return returnList, "100", nil
}

//分享中用户推荐
func (service recommend) RmduserUsers() ([]map[string]string, string, error) { //返回值，code，err
	log.Debug("method_start", "RmduserUsers", "input", "nil")
	request, err := redisClient.Lrange("rmduser:uids", 0, 11)
	if err != nil {
		return nil, "", err
	}
	returnList := []map[string]string{}
	for _, v := range request {
		mysqlUser, err0 := userClient.SearchUser(&user.UserRequest{UserId: v})
		if err0 != nil {
			return nil, "", err0
		}
		if mysqlUser.Err != "" {
			return nil, "", errors.New(mysqlUser.Err)
		}
		requestUser := map[string]string{}
		requestUser["userId"] = mysqlUser.User.UserId
		requestUser["imageName"] = mysqlUser.User.ImageName
		requestUser["nickName"] = mysqlUser.User.NickName
		returnList = append(returnList, requestUser)
	}
	log.Debug("method_end", "RmduserUsers", "status", "success")
	return returnList, "100", nil
}

//分享中最热动态推荐
func (service recommend) HotDynamic(pageInfo *PageInfo) ([]map[string]interface{}, string, string, error) { //返回值，最后一页标识，code，err
	log.Debug("method_start", "HotDynamic", "input", fmt.Sprint(pageInfo))
	returnList, pageFlag, code, err := dy(pageInfo, "0")
	if err != nil {
		return nil, "", "", err
	}
	log.Debug("method_end", "HotDynamic", "status", "success")
	return returnList, pageFlag, code, nil
}

//分享中新动态推荐
func (service recommend) NewDynamic(pageInfo *PageInfo) ([]map[string]interface{}, string, string, error) { //返回值，最后一页标识，code，err
	log.Debug("method_start", "HotDynamic", "input", fmt.Sprint(pageInfo))
	returnList, pageFlag, code, err := dy(pageInfo, "1")
	if err != nil {
		return nil, "", "", err
	}
	log.Debug("method_end", "HotDynamic", "status", "success")
	return returnList, pageFlag, code, nil
}

//倒叙分页分享中最热动态推荐
func (service recommend) ReverseHotDynamic(pageInfo *PageInfo) ([]map[string]interface{}, string, string, error) { //返回值，最后一页标识，code，err
	log.Debug("method_start", "HotDynamic", "input", fmt.Sprint(pageInfo))
	returnList, pageFlag, code, err := Reversedy(pageInfo, "0")
	if err != nil {
		return nil, "", "", err
	}
	log.Debug("method_end", "HotDynamic", "status", "success")
	return returnList, pageFlag, code, nil
}

//倒叙分页分享中新动态推荐
func (service recommend) ReverseNewDynamic(pageInfo *PageInfo) ([]map[string]interface{}, string, string, error) { //返回值，最后一页标识，code，err
	log.Debug("method_start", "HotDynamic", "input", fmt.Sprint(pageInfo))
	returnList, pageFlag, code, err := Reversedy(pageInfo, "1")
	if err != nil {
		return nil, "", "", err
	}
	log.Debug("method_end", "HotDynamic", "status", "success")
	return returnList, pageFlag, code, nil
}

func dy(pageInfo *PageInfo, s string) ([]map[string]interface{}, string, string, error) { //返回值，最后一页标识，code，err
	mKey := ""
	if "0" == s {
		mKey = "hotdynamic:dids"
	} else if "1" == s {
		mKey = "newDynamic"
	}
	log.Debug("mKey", mKey)
	pageSizeInt, _ := strconv.Atoi(pageInfo.PageSize)
	pageSize := int64(pageSizeInt)

	var pagebegin int64
	var lenlist int64
	pageFlag := "0"
	if "0" != pageInfo.LastId {
		scorre, err := redisClient.ZRANK(mKey, pageInfo.LastId)
		if err != nil {
			if "redigo: nil returned" != err.Error() {
				return nil, "", "", err
			} else if "redigo: nil returned" == err.Error() {
				return nil, pageFlag, flag1, nil
			}

		}
		if scorre != 0 {
			pagebegin = scorre

		} else {
			return nil, "", flag1, nil
		}
		lenlist, err = redisClient.Zlen(mKey)
		if err != nil {
			return nil, "", "", err
		}
		pagebegin = lenlist - pagebegin //zrank是正序，而列表查找是倒叙，所以相减一下
	}
	lenlist, err := redisClient.Zlen(mKey)
	if err != nil {
		return nil, "", "", err
	}
	if pagebegin+pageSize >= lenlist {
		pageFlag = "1"
	}
	log.Debug("████████pagebegin", pagebegin, "████████pageSize", pageSize, "████████lenlist", lenlist)
	reqList, err0 := redisClient.Zrange(mKey, pagebegin, pageSize+pagebegin-1)
	log.Debug("████████reqList", fmt.Sprint(reqList))
	if err0 != nil {
		return nil, "", "", err0
	}
	//把此切片调用动态模块获取动态list
	g := &grpdynamic.ListDynamicRequest{UserId: pageInfo.LookUserId, Ids: reqList}
	AllDynamics, err1 := grpdynamicClient.AllDynamics(g)
	if err1 != nil {
		return nil, "", "", err1
	}
	if AllDynamics.Err != "" {
		return nil, "", "", errors.New(AllDynamics.Err)
	}
	if AllDynamics.Code != "100" {
		return nil, "", AllDynamics.Code, nil
	}

	log.Debug("████████请求成功:", fmt.Sprint(AllDynamics.Dynamics))
	return AllDynamics.Dynamics, pageFlag, "100", nil
}

func Reversedy(pageInfo *PageInfo, s string) ([]map[string]interface{}, string, string, error) { //返回值，最后一页标识，code，err
	mKey := ""
	if "0" == s {
		mKey = "hotdynamic:dids"
	} else if "1" == s {
		mKey = "newDynamic"
	}
	log.Debug("mKey", mKey)
	pageSizeInt, _ := strconv.Atoi(pageInfo.PageSize)
	pageSize := int64(pageSizeInt)

	var pagebegin int64
	var lenlist int64
	pageFlag := "0"
	if "0" != pageInfo.LastId {
		scorre, err := redisClient.ZRANK(mKey, pageInfo.LastId)
		if err != nil {
			if "redigo: nil returned" != err.Error() {
				return nil, "", "", err
			} else if "redigo: nil returned" == err.Error() {

				return nil, pageFlag, flag1, nil
			}

		}
		pagebegin = scorre
		lenlist, err = redisClient.Zlen(mKey)
		if err != nil {
			return nil, "", "", err
		}
		pagebegin = lenlist - pagebegin //zrank是正序，而列表查找是倒叙，所以相减一下
	}
	lenlist, err := redisClient.Zlen(mKey)
	if err != nil {
		return nil, "", "", err
	}
	if pagebegin-pageSize <= 0 {
		pageSize = pagebegin
		pageFlag = "1"
	}
	log.Debug("████████pagebegin", pagebegin, "████████pageSize", pageSize, "████████lenlist", lenlist)
	reqList, err0 := redisClient.Zrange(mKey, pagebegin-pageSize, pagebegin-1)
	log.Debug("reqList", fmt.Sprint(reqList))
	if err0 != nil {
		if err0.Error() != "redigo: nil returned" {
			return nil, "", "", err0
		}
	}
	//把此切片调用动态模块获取动态list
	g := &grpdynamic.ListDynamicRequest{UserId: pageInfo.LookUserId, Ids: reqList}
	AllDynamics, err1 := grpdynamicClient.AllDynamics(g)
	if err1 != nil {
		return nil, "", "", err1
	}
	if AllDynamics.Err != "" {
		return nil, "", "", errors.New(AllDynamics.Err)
	}
	if AllDynamics.Code != "100" {
		return nil, "", AllDynamics.Code, nil
	}

	log.Debug("████████请求成功:", fmt.Sprint(AllDynamics.Dynamics))
	return AllDynamics.Dynamics, pageFlag, "100", nil

}

//添加好友推荐
func (service recommend) FriendRecommend(pageInfo *PageInfo) ([]map[string]interface{}, string, error) { //返回值，code，err
	log.Debug("method_start", "FriendRecommend", "input", "nil")
	redisUserIds, err := redisClient.Zrange("FriendRecommend", 0, 100)
	if err != nil {
		return nil, "", err
	}
	userIds := []string{}
	for _, v := range redisUserIds {
		if v != pageInfo.LookUserId { //排除自己
			//查看我是否关注此人
			score, err5 := redisClient.Zscore("grpfans:uIdFansList:"+v, pageInfo.LookUserId)
			if err5 != nil {
				if "redigo: nil returned" != err5.Error() {
					return nil, "", err5
				}

			}
			if score == 0 { //代表未关注过这个人
				userIds = append(userIds, v)
				if len(userIds) >= 10 {
					break
				}
			}
		}
	}
	returnList := []map[string]interface{}{}
	if len(userIds) == 0 {
		return returnList, "100", nil
	}
	fmt.Println("███userIds:", userIds)

	request, err1 := userClient.SearchUsersByUserIds(&user.UserRequest{UserIds: userIds})

	fmt.Println("███request:", request)
	if err1 != nil {
		return nil, "", err1
	}
	if request.Err != "" {
		return nil, "", errors.New(request.Err)
	}

	for _, v := range request.Users {
		v["likeFlag"] = "0"
		returnMap := map[string]interface{}{}
		returnMap["user"] = v
		dyIds, err2 := user4Dynamic(v["userId"])
		if err2 != nil {
			return nil, "", err2
		}
		//把此切片调用动态模块获取动态list
		g := &grpdynamic.ListDynamicRequest{UserId: pageInfo.LookUserId, Ids: dyIds}
		AllDynamics, err1 := grpdynamicClient.AllDynamics(g)
		if err1 != nil {
			return nil, "", err1
		}
		if AllDynamics.Err != "" {
			return nil, "", errors.New(AllDynamics.Err)
		}
		if AllDynamics.Code != "100" {
			return nil, AllDynamics.Code, nil
		}
		returnMap["dynamics"] = AllDynamics.Dynamics
		if len(AllDynamics.Dynamics) > 4 {
			returnMap["dynamics"] = AllDynamics.Dynamics[:4]
		}
		returnList = append(returnList, returnMap)
	}
	log.Debug("method_end", "FriendRecommend", "status", "success")
	return returnList, "100", nil
}

func user4Dynamic(userId string) ([]string, error) {
	result := []string{}
	uIds, err := redisClient.LrangeAll("grpdynamic:U2D:" + userId)
	if err != nil {
		return nil, err
	}
	result = uIds
	index := 0
	upIds, err := redisClient.LrangeAll("grpdynamic:UP2D:" + userId)
	if err != nil {
		return nil, err
	}
	for i, v := range uIds {
		for _, va := range upIds {
			if v == va {
				result = append(result[:i-index], result[i-index+1:]...)
				index = index + 1
			}
		}
	}
	if len(result) > 3 {
		result = result[0:3]
	}
	return result, nil
}

//存储推荐好友（周发布动态最多）
func (service recommend) SaveFriendRecommend() (string, error) {
	log.Debug("method_start", "SaveFriendRecommend", "input", "")
	weekTime, err0 := WeekTime() //获取上周时间戳
	if err0 != nil {
		return "", err0
	}
	weekNewDynamic, err := redisClient.ZRANGEBYSCORE("newDynamic", weekTime) //这一周发布的所有动态
	if err != nil {
		return "", err
	}

	err4 := redisClient.Del("temporary_SaveFriendRecommend") //插入前先删除临时排序表
	if err4 != nil {
		return "", err4
	}
	for _, v := range weekNewDynamic {
		userId, err1 := redisClient.Get("grpdynamic:D2U:" + v)
		if err1 != nil {
			if "redigo: nil returned" != err1.Error() {
				return "", err1
			}
		}
		if userId == "" {
			continue
		}
		err2 := redisClient.ZincrScore("temporary_SaveFriendRecommend", 1, userId)
		if err2 != nil {
			return "", err2
		}
	}
	temporarylist, err3 := redisClient.Zrange("temporary_SaveFriendRecommend", 0, 99) //取100个存储起来
	if err3 != nil {
		return "", err3
	}
	_, err4 = redisClient.Zremrangebyrank("FriendRecommend", 0, -1)
	log.Debug("████", fmt.Sprint("清空数据"))
	if err4 != nil {
		return "", err4
	}
	for _, v := range temporarylist {
		dycount, _ := redisClient.Zscore("temporary_SaveFriendRecommend", v)
		_, err5 := redisClient.Zadd("FriendRecommend", dycount, v)
		if err5 != nil {
			return "", err5
		}
	}
	redisClient.Del("temporary_SaveRmduserUsers") //操作后删除临时表
	log.Debug("method_end", "SaveFriendRecommend", "status", "success")
	return "100", nil
}

//搜索关键字推荐
func (service recommend) SearchRecommend() ([]string, string, error) { //返回值，code，err
	log.Debug("method_start", "SearchRecommend", "input", "nil")
	returnList, err := redisClient.Zrange("SearchRecommend", 0, -1)
	if err != nil {
		return returnList, "", err
	}
	log.Debug("method_end", "SearchRecommend", "status", "success")
	return returnList, "100", nil
}

//给今天发布的动态推送点赞通知
func (service recommend) PushFavour() (string, error) {
	log.Debug("method_start", "PushFavour", "input", "")
	yesDayTime2, err0 := YesDayTime(-2) //获取前天此时的时间戳，参数代表天数正负代表加减
	if err0 != nil {
		return "", err0
	}
	yesDayTime1, err0 := YesDayTime(-1) //获取昨天此时的时间戳，参数代表天数正负代表加减
	if err0 != nil {
		return "", err0
	}
	yesDayNewDynamic, err := redisClient.ZRANGEBYSCORE("newDynamic", yesDayTime2) //这2天发布的所有动态
	if err != nil {
		return "", err
	}
	pushMap := map[string]string{}
	for _, v := range yesDayNewDynamic {
		likes, _ := redisClient.ZRANGEBYSCORE("grpdynamic:UL:"+v, yesDayTime1) //此动态今天有没有被点赞
		if len(likes) > 0 {
			userId, _ := redisClient.Hget("grpdynamic:D:"+v, "UserId")
			pushMap[userId] = v //把用户id装到map的key中，模拟set去重
		}
	}
	//开始推送
	list1 := []push.PushInfo{}
	for k, v := range pushMap {
		title := "今天你的动态被人点赞咯"
		text2 := "{\"alert\":\"" + title + "\",\"extras\":{\"type\":\"DY1\",\"DynamicId\":\"" + v + "\"}}"
		pushInfo1 := push.PushInfo{
			Title:    title,
			Text:     " ",
			JsonInfo: text2,
			Alias:    getAlias(k),
		}

		list1 = append(list1, pushInfo1)

	}
	req1 := &push.PushRequest{
		PushType:  "0", //0代表推送消息，1代表透传消息
		PushInfos: list1,
	}
	pushClient.PushErrorLog(req1)
	log.Debug("method_end", "SaveFriendRecommend", "status", "success")
	return "100", nil
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

////给明天要开始的活动参与者推送通知
func (service recommend) PushStartActivity() (string, error) {
	log.Debug("method_start", "PushStartActivity", "input", "")
	fmt.Println("-------------------------------------------------")
	acReq, _ := activityClient.StartAc()
	fmt.Println("acReq", acReq)
	acIdList := acReq.Data
	//acIdList := []map[string]string{}
	list1 := []push.PushInfo{}
	for _, v := range acIdList {
		acId := v["acId"]
		menbers, _ := redisClient.Zrange("mactivity:ActivityMemberList:"+acId, 0, -1) //获取此活动的所有成员Id
		for _, mv := range menbers {
			title := "你报名的活动就要开始了，记得来哦"
			text2 := "{\"alert\":\"" + title + "\",\"extras\":{\"type\":\"ACST1\",\"acId\":\"" + acId + "\"}}"
			pushInfo1 := push.PushInfo{
				Title:    title,
				Text:     " ",
				JsonInfo: text2,
				Alias:    getAlias(mv),
			}
			list1 = append(list1, pushInfo1)
		}
	}
	req1 := &push.PushRequest{
		PushType:  "0", //0代表推送消息，1代表透传消息
		PushInfos: list1,
	}
	pushClient.PushErrorLog(req1)
	log.Debug("method_end", "PushStartActivity", "status", "success")
	return "100", nil
}

func (service recommend) SaveHomePageCache() (string, error) {
	log.Debug("method_start", "SaveHomePageCache", "input", "")
	fmt.Println("-------------------------------------------------")
	cache := make(map[string]string, 0)
	grpResp, gerr := mgrpmgrClient.ListGrps(&mgrpmgr.GrpRequest{GtId: "Gt1", LastId: "0", PageSize: "6"})
	if gerr != nil {
		log.Debug("█gerr", gerr.Error())
	} else {
		log.Debug("█grpResp.Code", grpResp.Code)
		log.Debug("█grpResp.err", grpResp.Err)
	}

	grpsresponse, _ := json.Marshal(grpResp.Grops)
	cache["grpsresponse"] = string(grpsresponse)
	actResp, aerr := mactivityClient.SearchAllActivity(&mactivity.SeachActivityInfo{GongYiActivityID: "1", Ad: "a"})
	fmt.Println("-------------------------------------------------actResp,aerr", actResp, aerr)
	if aerr != nil {
		log.Debug("█aerr", aerr.Error())
	} else {
		log.Debug("█actResp.Code", actResp.Code)
		log.Debug("█actResp.err", actResp.Err)
	}
	activitysresponse, _ := json.Marshal(actResp.Data)
	cache["activitysresponse"] = string(activitysresponse)
	returnList := []map[string]string{}
	helpSlogan, err := redisClient.HgetAllMap("helpSlogan") //求助
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			log.Debug("█err", err.Error())
		}
	}
	publicBenefitSlogan, err1 := redisClient.HgetAllMap("publicBenefitSlogan") //公益
	if err1 != nil {
		if err1.Error() != "redigo: nil returned" {
			log.Debug("█err1", err1.Error())
		}
	}
	returnList = append(returnList, helpSlogan)
	returnList = append(returnList, publicBenefitSlogan)
	advertisement, _ := json.Marshal(returnList)
	cache["advertisement"] = string(advertisement)

	err = redisClient.Hmset("cache:homepage", cache)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	log.Debug("method_end", "SaveHomePageCache", "status", "success")
	return "100", nil
}

func (service recommend) AddDynamicFans() (string, error) {
	log.Debug("method_start", "AddDynamicFans", "input", "")

	t2, err := DifferenceTime(-2 * time.Hour)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	t4, err := DifferenceTime(-4 * time.Hour)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	t6, err := DifferenceTime(-6 * time.Hour)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	t8, err := DifferenceTime(-8 * time.Hour)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	t10, err := DifferenceTime(-10 * time.Hour)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	t12, err := DifferenceTime(-12 * time.Hour)
	if err != nil {
		log.Debug("█err", err.Error())
	}
	dyIds1, err1 := redisClient.ZRANGEBYSCORE2("newDynamic", t4, t2)

	if err1 != nil {
		log.Debug("█err1", err1.Error())
	}
	addFans(dyIds1, 0, 80)
	dyIds2, err2 := redisClient.ZRANGEBYSCORE2("newDynamic", t6, t4)
	if err2 != nil {
		log.Debug("█err2", err2.Error())
	}
	addFans(dyIds2, 100, 180)
	dyIds3, err3 := redisClient.ZRANGEBYSCORE2("newDynamic", t8, t6)
	if err3 != nil {
		log.Debug("█err3", err3.Error())
	}
	addFans(dyIds3, 200, 280)
	dyIds4, err4 := redisClient.ZRANGEBYSCORE2("newDynamic", t10, t8)
	if err4 != nil {
		log.Debug("█err4", err4.Error())
	}
	addFans(dyIds4, 300, 380)
	dyIds5, err5 := redisClient.ZRANGEBYSCORE2("newDynamic", t12, t10)
	if err5 != nil {
		log.Debug("█err5", err5.Error())
	}
	addFans(dyIds5, 400, 480)
	log.Debug("method_end", "AddDynamicFans", "status", "success")
	return "100", nil
}
func addFans(dyIds []string, start int, end int) {
	if len(dyIds) > 0 {
		for _, v := range dyIds {
			rand.Seed(time.Now().Unix())
			r1 := rand.Intn(19) + start
			r2 := rand.Intn(19) + end
			for i := r1; i < r2; i++ {
				_, err := redisClient.Zadd("grpdynamic:ULD:"+fanUId+strconv.Itoa(i), getNowTime(), v)
				if err != nil {
					log.Debug("█err", err.Error())
				}
				_, err = redisClient.Zadd("grpdynamic:UL:"+v, getNowTime(), fanUId+strconv.Itoa(i))
				if err != nil {
					log.Debug("█err", err.Error())
				}
				userId, err := redisClient.Hget("grpdynamic:D:"+v, "UserId")
				if err != nil {
					log.Debug("█err", err.Error())
				}
				_, err = redisClient.Zadd("grpdynamic:USERALLLIKE:"+userId, getNowTime(), fanUId+strconv.Itoa(i)+","+v)
				if err != nil {
					log.Debug("█err", err.Error())
				}
			}
		}
	}
}

func getNowTime() int64 {
	return int64(time.Now().Unix())
}

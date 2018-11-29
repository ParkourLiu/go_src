package main

import (
	"context"
	"flag"
	"fmt"
	"mtcomm/db/redis"
	logger "mtcomm/log"
	"testing"

	"github.com/tjz101/goprop"
)

func init() {
	/* init properties */
	propFile := flag.String("prop", "prop.properties", "properties file")
	flag.Parse()

	prop = goprop.NewProp()
	prop.Read(*propFile)

	namespace = prop.Get("namespace") //kubernetes namespace

	/* init log */
	logger.SetDefaultLogLevel(logger.LevelDebug)
	logger.With("serviceName", serviceName)

	redisInfo := &redis.RedisServerInfo{
		Ctx:       context.TODO(),
		Logger:    logger.GetDefaultLogger(),
		RedisHost: prop.Get("redis_host"),
	}
	redisClient = redis.NewRedisClient(redisInfo)

}

//	ListAttentionUids(info *GrpFansInfo) (i map[string]interface{}, errr error)

func TestAttention1(t *testing.T) {
	//社群添加粉丝
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "TestFansId"}
	code, err := svc.Attention(info)
	fmt.Println(err, "err=======")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	//查看该社群的粉丝列表里是否已添加过该粉丝
	zhi, err := dao.Zscore(redisClient, prop.Get("groupFansList")+info.GroupId, info.FansId)
	//查看该粉丝的关注列表里是否已关注了社群
	zhi2, err2 := dao.Zscore(redisClient, prop.Get("attentionGroupList")+info.FansId, info.GroupId)
	if zhi == 0 || zhi2 == 0 || err != nil || err2 != nil {
		t.Error("社群添加粉丝失败")
		return
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	//测试输入参数是否合法 1
	info0 := &GrpFansInfo{GroupId: "Test   GroupId", FansId: "Test  FansId", UId: "TestUId"}
	_, err0 := svc.Attention(info0)
	if err0 != nil {
		if err0.Error() == prop.Get("101") {

		} else {
			return
		}
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info0.GroupId, info0.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info0.FansId, info0.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	//测试输入参数是否合法 2
	info00 := &GrpFansInfo{GroupId: " ", FansId: "Test  FansId", UId: " "}
	_, err00 := svc.Attention(info00)
	if err00 != nil {
		if err00.Error() == prop.Get("102") {

		} else {
			return
		}
	}

	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info00.GroupId, info00.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info00.FansId, info00.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	//测试输入参数是否合法3
	info000 := &GrpFansInfo{GroupId: "sfgs ", FansId: " "}
	_, err000 := svc.Attention(info000)
	if err000 != nil {
		if err000.Error() == prop.Get("103") {

		} else {
			return
		}
	}

	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info000.GroupId, info000.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info000.FansId, info000.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	//测试异常1
	info1 := &GrpFansInfo{GroupId: "Test   GroupId", FansId: "Test  FansId"}
	code11, err1 := svc.Attention(info1)
	if err1 != nil {
		t.Error(err1.Error())
		return
	}
	if code11 != "100" {
		t.Error("code不为100")
	}
	_, err11 := svc.Attention(info1)
	if err11 != nil {
		if err11.Error() == "社群添加粉丝失败，该用户已是社群的粉丝了" {

		} else {
			return
		}
	}

	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info1.GroupId, info1.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info1.FansId, info1.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

}
func TestAttention2(t *testing.T) {
	//活动添加粉丝ActivityId
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}
	infoAct := &GrpFansInfo{ActivityId: "TestActivityId", FansId: "TestFansId"}
	code1, errAct := svc.Attention(infoAct)
	if errAct != nil {
		t.Error(errAct.Error())
		return
	}
	if code1 != "100" {
		t.Error("code不为100")
	}
	//查看该活动的粉丝列表里是否已添加过该粉丝
	zhiAct, errAct := dao.Zscore(redisClient, prop.Get("activityfansList")+infoAct.ActivityId, infoAct.FansId)
	//查看该粉丝的关注列表里是否已关注了活动
	zhi2Act, err2Act := dao.Zscore(redisClient, prop.Get("attentionActivityList")+infoAct.FansId, infoAct.ActivityId)
	if zhiAct == 0 || zhi2Act == 0 || errAct != nil || err2Act != nil {
		t.Error("活动添加粉丝失败", zhiAct, zhi2Act, errAct, err2Act)
		return
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("activityfansList") + infoAct.ActivityId, infoAct.FansId}
		re[1] = []string{prop.Get("attentionActivityList") + infoAct.FansId, infoAct.ActivityId}
		_, errAct := dao.UnAttention(redisClient, re)
		if errAct != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//异常1
	infoAct1 := &GrpFansInfo{ActivityId: "TestActivityId异常1", FansId: "TestFansId异常1"}
	_, errAct1 := svc.Attention(infoAct1)
	if errAct1 != nil {
		t.Error(errAct1.Error())
		return
	}
	_, errAct2 := svc.Attention(infoAct1)
	if errAct2 != nil {
		if errAct2.Error() == prop.Get("105") {

		} else {
			return
		}
	}

	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("activityfansList") + infoAct1.ActivityId, infoAct1.FansId}
		re[1] = []string{prop.Get("attentionActivityList") + infoAct1.FansId, infoAct1.ActivityId}
		_, errAct := dao.UnAttention(redisClient, re)
		if errAct != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
}
func TestAttention3(t *testing.T) {
	//个人添加粉丝
	svc := &grpFansService{}
	dao := &GrpFansDao{}
	infoUId := &GrpFansInfo{UId: "TestUId个人", FansId: "TestFansId个人"}
	code, errUId := svc.Attention(infoUId)
	if errUId != nil {
		t.Error(errUId.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}

	//查看该用户的粉丝列表里是否已添加过该粉丝
	zhiUId, errUId2 := dao.Zscore(redisClient, prop.Get("uIdFansList")+infoUId.UId, infoUId.FansId)
	//查看该粉丝的关注列表里是否已关注了该用户
	zhi2UId, err2UId := dao.Zscore(redisClient, prop.Get("attentionUIdList")+infoUId.FansId, infoUId.UId)
	if zhiUId == 0 || zhi2UId == 0 || errUId2 != nil || err2UId != nil {
		t.Error("个人添加粉丝失败", zhiUId, zhi2UId, errUId2, err2UId)
		return
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + infoUId.UId, infoUId.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + infoUId.FansId, infoUId.UId}
		_, errUIdq := dao.UnAttention(redisClient, re)
		if errUIdq != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//异常1
	infoUId1 := &GrpFansInfo{UId: "Test  UId", FansId: "TestFansId"}
	code2, errUId11 := svc.Attention(infoUId1)
	if errUId11 != nil {
		t.Error(errUId11.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	_, errUId22 := svc.Attention(infoUId1)
	if errUId22 != nil {
		if errUId22.Error() == prop.Get("106") {

		} else {
			return
		}
	}

	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + infoUId1.UId, infoUId1.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + infoUId1.FansId, infoUId1.UId}
		_, errUIds := dao.UnAttention(redisClient, re)
		if errUIds != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
}

func TestUnAttention1(t *testing.T) {
	//社群删除粉丝
	svc := &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "TestFansId"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	infoUId := &GrpFansInfo{GroupId: "TestGroupId", UnFansId: "TestFansId"}
	code1, errUId := svc.UnAttention(infoUId)
	if errUId != nil {
		t.Error(errUId.Error())
		return
	}
	if code1 != "100" {
		t.Error("code不为100")
	}
	//测试输入参数是否合法 1
	info0 := &GrpFansInfo{GroupId: "Test   GroupId", UnFansId: "Test  UnFansId", UId: "TestUId"}
	_, err0 := svc.UnAttention(info0)
	if err0 != nil {
		if err0.Error() == prop.Get("101") {

		} else {
			return
		}
	}
	//测试输入参数是否合法 2
	info00 := &GrpFansInfo{GroupId: " ", UnFansId: "Test  UnFansId"}
	_, err00 := svc.UnAttention(info00)
	if err00 != nil {
		if err00.Error() == prop.Get("102") {

		} else {
			return
		}
	}

	//测试输入参数是否合法3
	info000 := &GrpFansInfo{GroupId: "sfgs ", UnFansId: " "}
	_, err000 := svc.UnAttention(info000)
	if err000 != nil {
		if err000.Error() == prop.Get("107") {

		} else {
			return
		}
	}

}
func TestUnAttention2(t *testing.T) {
	//活动删除粉丝ActivityId
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	info := &GrpFansInfo{ActivityId: "TestGroupId", FansId: "TestFansId"}
	_, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("activityfansList") + info.ActivityId, info.FansId}
		re[1] = []string{prop.Get("attentionActivityList") + info.FansId, info.ActivityId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	infoAct := &GrpFansInfo{ActivityId: "TestActivityId", UnFansId: "TestFansId"}
	_, errAct := svc.UnAttention(infoAct)
	if errAct != nil {
		t.Error(errAct.Error())
		return
	}

}
func TestUnAttention3(t *testing.T) {
	//个人删除粉丝
	svc := &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	info := &GrpFansInfo{UId: "TestUId", FansId: "TestUnFansId"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	infoUId := &GrpFansInfo{UId: "TestUId", UnFansId: "TestUnFansId"}
	code2, errUId := svc.UnAttention(infoUId)
	if errUId != nil {
		t.Error(errUId.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
}

func TestFansNum1(t *testing.T) {
	//社群粉丝数
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	infoF := &GrpFansInfo{GroupId: "TestGroupId", FansId: "TestFansId"}
	code, errF := svc.Attention(infoF)
	if errF != nil {
		t.Error(errF.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + infoF.GroupId, infoF.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + infoF.FansId, infoF.GroupId}
		_, errF := dao.UnAttention(redisClient, re)
		if errF != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "TestFansId"}
	data, code, err := svc.FansNum(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	i := data.(int64)
	if i != int64(1) {
		t.Error("粉丝数不等于1")
		return
	}

}
func TestFansNumParm(t *testing.T) {
	var svc GrpFansService
	svc = &grpFansService{}
	//测试输入参数是否合法 1
	info0 := &GrpFansInfo{GroupId: "Test   GroupId", FansId: "Test  FansId", UId: "TestUId"}
	_, _, err0 := svc.FansNum(info0)
	if err0 != nil {
		if err0.Error() == prop.Get("101") {

		} else {
			return
		}
	}
	//测试输入参数是否合法 2
	info00 := &GrpFansInfo{GroupId: " ", FansId: "Test  FansId", UId: " "}
	_, _, err00 := svc.FansNum(info00)
	if err00 != nil {
		if err00.Error() == prop.Get("102") {

		} else {
			return
		}
	}

}
func TestFansNum2(t *testing.T) {
	//活动粉丝数
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	info := &GrpFansInfo{ActivityId: "TestGroupId", FansId: "TestFansId"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("activityfansList") + info.ActivityId, info.FansId}
		re[1] = []string{prop.Get("attentionActivityList") + info.FansId, info.ActivityId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	info11 := &GrpFansInfo{ActivityId: "TestGroupId", FansId: "TestFansId"}
	data, code, err := svc.FansNum(info11)
	if err != nil {
		t.Error(err.Error())
		return
	}
	i := data.(int64)
	if i != int64(1) {
		t.Error("粉丝数不等于1")
		return
	}

}
func TestFansNum3(t *testing.T) {
	//个人粉丝数
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	info := &GrpFansInfo{UId: "TestGroupId", FansId: "TestFansId"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	info11 := &GrpFansInfo{UId: "TestGroupId", FansId: "TestFansId"}
	data, code2, err := svc.FansNum(info11)
	if err != nil {
		t.Error(err.Error())
		return
	}
	i := data.(int64)
	if i != int64(1) {
		t.Error("粉丝数不等于1")
		return
	}
	if code2 != "100" {
		t.Error("code不等于100")
	}

}

func TestAttentionNum(t *testing.T) {
	//个人关注数
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}
	//准备数据
	info := &GrpFansInfo{UId: "TestGroupId", FansId: "TestFansId"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	defer func() {
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
	}()

	info11 := &GrpFansInfo{UId: "TestFansId"}
	data, code2, err := svc.AttentionNum(info11)
	if err != nil {
		t.Error(err.Error())
		return
	}
	i := data.(int64)
	if i != int64(1) {
		t.Error("粉丝数不等于1")
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//测试异常1
	info2 := &GrpFansInfo{UId: " ", FansId: "TestFansId"}
	_, _, err2 := svc.AttentionNum(info2)
	if err2 != nil {
		if err2.Error() == prop.Get("109") {

		} else {
			t.Error(err2.Error())
		}
	}
}

func TestListMembers(t *testing.T) {
	//测试列出社群的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//社群添加粉丝两个
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//社群添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除社群和粉丝的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("groupFansList") + info2.GroupId, info2.FansId}
		re2[1] = []string{prop.Get("attentionGroupList") + info2.FansId, info2.GroupId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{GroupId: "TestGroupId"}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data========", data)

}

//lastUid 为空
func TestListMembers2(t *testing.T) {
	//测试列出社群的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//社群添加粉丝两个
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}

	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//社群添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除社群和粉丝的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("groupFansList") + info2.GroupId, info2.FansId}
		re2[1] = []string{prop.Get("attentionGroupList") + info2.FansId, info2.GroupId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{GroupId: "TestGroupId", LastUId: "testUid2"}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data====", data, "code==", code)

}

//108
func TestListMembers3(t *testing.T) {
	//测试列出社群的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//社群添加粉丝两个
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//社群添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除社群和粉丝的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("groupFansList") + info2.GroupId, info2.FansId}
		re2[1] = []string{prop.Get("attentionGroupList") + info2.FansId, info2.GroupId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{GroupId: "TestGroupId", LastUId: " sd "}
	data, code, err := svc.ListMembers(info3)

	fmt.Println("data==", data, "code==", code, "err==", err)

}

//最后一条记录
func TestListMembers4(t *testing.T) {
	//测试列出社群的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//社群添加粉丝两个
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//社群添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除社群和粉丝的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("groupFansList") + info2.GroupId, info2.FansId}
		re2[1] = []string{prop.Get("attentionGroupList") + info2.FansId, info2.GroupId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{GroupId: "TestGroupId", LastUId: "testUid"}
	data3, code3, err := svc.ListMembers(info3)
	if err != nil {
		t.Error(err.Error())
	}
	if code3 == prop.Get("108") {

	} else {
		t.Error("code不为108")
	}
	fmt.Println("data3===", data3)

}

//info.Count != 0
func TestListMembers5(t *testing.T) {
	//测试列出社群的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//社群添加粉丝两个
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{GroupId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//社群添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除社群和粉丝的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("groupFansList") + info.GroupId, info.FansId}
		re[1] = []string{prop.Get("attentionGroupList") + info.FansId, info.GroupId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("groupFansList") + info2.GroupId, info2.FansId}
		re2[1] = []string{prop.Get("attentionGroupList") + info2.FansId, info2.GroupId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{GroupId: "TestGroupId", Count: int64(3)}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data=====", data, "code==", code)

}

func TestListMembersGeRen(t *testing.T) {
	//测试列出个人的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "TestGroupId", Count: int64(3)}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data=====", data, "code==", code)
}

//lastUid 为空
func TestListMembersGeRen2(t *testing.T) {
	//测试列出个人的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "TestGroupId", LastUId: " "}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data======", data, "code===", code)
}

//108
func TestListMembersGeRen3(t *testing.T) {
	//测试列出个人的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testHmset", "testHmset2")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "TestGroupId", LastUId: " sdfsdf"}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data==", data, "code==", code)
}

//最后一条记录
func TestListMembersGeRen4(t *testing.T) {
	//测试列出个人的粉丝基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	hmset2 := []interface{}{"uid", "testUid2", "url", "tesurl2", "nickName", "testnickName2"}
	redisClient.Hmset("testUid2", hmset2)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid2"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testUid", "testUid2")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "TestGroupId", LastUId: "testUid2"}
	data, code, err := svc.ListMembers(info3)
	fmt.Println("data===", data, "code==", code)
}

func TestListAttentionUids(t *testing.T) {
	//测试列出个人关注的基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId2", FansId: "testUid"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testUid")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "testUid", LastUId: "TestGroupId2"}
	data, code, err := svc.ListAttentionUids(info3)
	fmt.Println("data=====", data, "code==", code)
}

//lastUid 为空
func TestListAttentionUids2(t *testing.T) {
	//测试列出个人关注的基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId2", FansId: "testUid"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testUid")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "testUid", LastUId: " "}
	data, code, err := svc.ListAttentionUids(info3)
	fmt.Println("data====", data, "code=", code)
}

//108
func TestListAttentionUids3(t *testing.T) {
	//测试列出个人关注的基本信息
	//准备数据
	//新建一个用户
	hmset := []interface{}{"uid", "testUid", "url", "tesurl", "nickName", "testnickName"}
	redisClient.Hmset("testUid", hmset)
	//用户点赞用户
	var svc GrpFansService
	svc = &grpFansService{}
	dao := &GrpFansDao{}

	info := &GrpFansInfo{UId: "TestGroupId", FansId: "testUid"}
	info2 := &GrpFansInfo{UId: "TestGroupId2", FansId: "testUid"}
	code, err := svc.Attention(info)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if code != "100" {
		t.Error("code不为100")
	}
	code2, err2 := svc.Attention(info2)
	if err2 != nil {
		t.Error(err2.Error())
		return
	}
	if code2 != "100" {
		t.Error("code不为100")
	}
	//个人添加禁言用户
	//删除数据
	defer func() {
		//删除用户
		redisClient.Del("testUid")
		//删除用户和用户的点赞关系
		re := make([][]string, 2)
		re[0] = []string{prop.Get("uIdFansList") + info.UId, info.FansId}
		re[1] = []string{prop.Get("attentionUIdList") + info.FansId, info.UId}
		_, err := dao.UnAttention(redisClient, re)
		if err != nil {
			t.Fatal("测试数据删除失败")
		}
		re2 := make([][]string, 2)
		re2[0] = []string{prop.Get("uIdFansList") + info2.UId, info2.FansId}
		re2[1] = []string{prop.Get("attentionUIdList") + info2.FansId, info2.UId}
		_, err2 := dao.UnAttention(redisClient, re2)
		if err2 != nil {
			t.Fatal("测试数据删除失败")
		}
	}()
	//验证数据"count":100

	info3 := &GrpFansInfo{UId: "testUid", LastUId: " 方法"}
	data, code, err := svc.ListAttentionUids(info3)
	fmt.Println("data===", data, "code==", code)
}

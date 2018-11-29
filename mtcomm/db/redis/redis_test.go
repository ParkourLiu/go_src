package redis_test

import (
	"context"
	"fmt"
	"mtcomm/db/redis"
	logger "mtcomm/log"
	"testing"
	"time"
)

var client redis.RedisClient

func init() {
	logger.SetDefaultLogLevel(logger.LevelDebug)
	info := &redis.RedisServerInfo{
		Ctx:           context.TODO(),
		Logger:        logger.GetDefaultLogger(),
		RedisHost:     "106.14.216.4:8888",
		RedisPassword: "zaq12wsx1",
	}
	client = redis.NewRedisClient(info)
}

func TestSetAndKeysAndExist(t *testing.T) {
	client.Set("lio_abc", "1")
	client.Set("lio_abd", "1")
	client.Set("lio_abe", "1")

	//	keys, err := client.Keys("lio_ab*")
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	if len(keys) != 3 {
	//		t.Error("value error")
	//	}

	exist, err1 := client.Exist("lio_abc")
	if err1 != nil {
		t.Error(err1)
	}
	if !exist {
		t.Error("value error")
	}

	client.Del("lio_abc", "lio_abd", "lio_abe")

	exist, err1 = client.Exist("lio_abc")
	if err1 != nil {
		t.Error(err1)
	}
	if exist {
		t.Error("value error")
	}
}

func TestIncr(t *testing.T) {
	client.Del("lio_incr")

	i, err := client.Incr("lio_incr")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Error("value error")
	}

	i, err = client.Incr("lio_incr")
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Error("value error")
	}

	i, err = client.Decr("lio_incr")
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Error("value error")
	}

	client.Del("lio_incr")
}

func TestExpire(t *testing.T) {
	client.Del("lio_Expire")

	_, err := client.Incr("lio_Expire")
	if err != nil {
		t.Error(err)
	}
	err = client.Expire("lio_Expire", 1)
	if err != nil {
		t.Error(err)
	}
	b, err1 := client.Exist("lio_Expire")
	if err1 != nil {
		t.Error(err1)
	}
	if !b {
		t.Error("err")
	}
	time.Sleep(1 * time.Second)
	b, err1 = client.Exist("lio_Expire")
	if err1 != nil {
		t.Error(err1)
	}
	if b {
		t.Error("err")
	}

	client.Del("lio_Expire")
}

func TestSetAndExpire(t *testing.T) {
	client.Del("lio_SetAndExpire")

	err := client.SetStringAndExpire("lio_SetAndExpire", "2", 1)
	if err != nil {
		t.Error(err)
	}
	v, err1 := client.Get("lio_SetAndExpire")
	if err1 != nil {
		t.Error(err1)
	}
	if v != "2" {
		t.Error("err")
	}
	time.Sleep(1 * time.Second)
	b, err2 := client.Exist("lio_SetAndExpire")
	if err2 != nil {
		t.Error(err2)
	}
	if b {
		t.Error("err")
	}

	client.Del("lio_SetAndExpire")
}

func TestGetStrings(t *testing.T) {
	client.Del("lio_GetStrings_1", "lio_GetStrings_2")
	client.Set("lio_GetStrings_1", "aaa")
	client.Set("lio_GetStrings_2", "bbb")

	vmap, err := client.GetStrings("lio_GetStrings_1", "lio_GetStrings_2")
	if err != nil {
		t.Error(err)
	}
	if vmap[0] != "aaa" || vmap[1] != "bbb" {
		t.Error("err")
	}

	client.Del("lio_GetStrings_1", "lio_GetStrings_2")
}

type p struct {
	Id   string
	Name string
}

func TestHmset1(t *testing.T) {
	client.Del("lio_Hmset_1")

	err := client.Hmset("lio_Hmset_1", &p{Id: "111", Name: "lio"})
	if err != nil {
		t.Error(err)
	}

	p1 := &p{}
	err = client.HgetAllStruction("lio_Hmset_1", p1)
	if err != nil {
		t.Error(err)
	}
	if p1.Id != "111" || p1.Name != "lio" {
		t.Error("err")
	}
	client.Del("lio_Hmset_1")
}

func TestHmset2(t *testing.T) {
	client.Del("lio_Hmset_2")

	err := client.Hmset("lio_Hmset_2", map[string]string{"Id": "111", "Name": "lio"})
	if err != nil {
		t.Error(err)
	}

	p1, err1 := client.HgetAllMap("lio_Hmset_2")
	if err1 != nil {
		t.Error(err1)
	}
	if p1["Id"] != "111" || p1["Name"] != "lio" {
		t.Error("err")
	}
	client.Del("lio_Hmset_2")
}

type ps struct {
	Id    string
	Name  string
	Count int64
}

func TestHsetgetincr(t *testing.T) {
	client.Del("lio_Hsetgetincr")

	client.Hmset("lio_Hsetgetincr", map[string]string{"Id": "111", "Name": "lio"})
	err := client.Hset("lio_Hsetgetincr", "Id", "222")
	if err != nil {
		t.Error(err)
	}

	//check struction 没有count字段会不会出错
	p3 := &ps{}
	err = client.HgetAllStruction("lio_Hsetgetincr", p3)
	if err != nil {
		t.Error(err)
	}
	if p3.Id != "222" || p3.Count != int64(0) {
		t.Error("err")
	}

	client.Hincr("lio_Hsetgetincr", "Count")
	if err != nil {
		t.Error(err)
	}

	//check map
	p1, err1 := client.HgetAllMap("lio_Hsetgetincr")
	if err1 != nil {
		t.Error(err1)
	}
	if p1["Id"] != "222" || p1["Count"] != "1" {
		t.Error("err")
	}
	//check struction
	p2 := &ps{}
	err = client.HgetAllStruction("lio_Hsetgetincr", p2)
	if err != nil {
		t.Error(err)
	}
	if p2.Id != "222" || p2.Count != int64(1) {
		t.Error("err")
	}

	client.Hdecr("lio_Hsetgetincr", "Count")
	if err != nil {
		t.Error(err)
	}
	c, _ := client.Hget("lio_Hsetgetincr", "Count")
	if c != "0" {
		t.Error("error")
	}

	client.Del("lio_Hsetgetincr")
}

func TestSet(t *testing.T) {
	client.Del("lio_set")

	err := client.Sadd("lio_set", "a", "b", "c", "d", "e")
	if err != nil {
		t.Error(err)
	}
	length, err1 := client.Slen("lio_set")
	if err1 != nil {
		t.Error(err1)
	}
	if length != int64(5) {
		t.Error("error")
	}

	err = client.Srem("lio_set", "b", "d", "e")
	if err != nil {
		t.Error(err)
	}

	exist, err2 := client.Sismember("lio_set", "a")
	if err2 != nil {
		t.Error(err2)
	}
	if !exist {
		t.Error("error")
	}
	exist, err2 = client.Sismember("lio_set", "b")
	if err2 != nil {
		t.Error(err2)
	}
	if exist {
		t.Error("error")
	}

	es, err3 := client.Smembers("lio_set")
	if err3 != nil {
		t.Error(err3)
	}
	if len(es) != 2 {
		t.Error("err")
	}
	if es[0] != "a" && es[0] != "c" {
		t.Error("err")
	}
	if es[1] != "a" && es[1] != "c" {
		t.Error("err")
	}

	client.Del("lio_set")
}

func TestZset(t *testing.T) {
	client.Del("lio_zset")

	err := client.ZincrScore("lio_zset", 1, "aaa")
	if err != nil {
		t.Error(err)
	}
	client.ZincrScore("lio_zset", 10, "ccc")
	client.ZincrScore("lio_zset", 2, "bbb")
	client.ZincrScore("lio_zset", 3, "bbb")

	length, err1 := client.Zlen("lio_zset")
	if err1 != nil {
		t.Error(err1)
	}
	if length != 3 {
		t.Error("err")
	}

	s, err2 := client.Zscore("lio_zset", "bbb")
	if err2 != nil {
		t.Error(err2)
	}
	if s != int64(5) {
		t.Error("err")
	}

	s, err2 = client.Zscore("lio_zset", "ccc")
	if err2 != nil {
		t.Error(err2)
	}
	if s != int64(10) {
		t.Error("err")
	}

	ss, err3 := client.ZrangeAll("lio_zset")
	if err3 != nil {
		t.Error(err3)
	}
	if len(ss) != 3 || ss[0] != "ccc" || ss[1] != "bbb" || ss[2] != "aaa" {
		t.Error("err")
	}

	ss, _ = client.Zrange("lio_zset", int64(0), int64(1))
	if len(ss) != 2 || ss[0] != "ccc" || ss[1] != "bbb" {
		t.Error("err")
	}

	err = client.Zrem("lio_zset", "aaa", "ccc")
	if err != nil {
		t.Error(err)
	}

	ss, _ = client.ZrangeAll("lio_zset")
	if len(ss) != 1 || ss[0] != "bbb" {
		t.Error("err")
	}

	client.Del("lio_zset")
}

func TestList1(t *testing.T) {
	client.Del("lio_list")
	client.LpushFromHead("lio_list", "2", "1")
	client.LpushFromTail("lio_list", "3", "4", "5")
	length, _ := client.Llen("lio_list")
	if length != int64(5) {
		t.Error("err")
	}
	ss, _ := client.LrangeAll("lio_list")
	if len(ss) != 5 || ss[0] != "1" || ss[1] != "2" || ss[2] != "3" || ss[3] != "4" || ss[4] != "5" {
		t.Error("err")
	}

	ss, _ = client.Lrange("lio_list", 0, 2)
	if len(ss) != 3 || ss[0] != "1" || ss[1] != "2" || ss[2] != "3" {
		t.Error("err")
	}

	client.Del("lio_list")
}

func TestList2(t *testing.T) {
	client.Del("lio_list2")
	client.LpushFromTail("lio_list2", "1", "2", "1", "3")
	ss, _ := client.LrangeAll("lio_list2")
	if len(ss) != 4 || ss[0] != "1" || ss[1] != "2" || ss[2] != "1" || ss[3] != "3" {
		t.Error("err")
	}
	err := client.Lrem("lio_list2", "1")
	if err != nil {
		t.Error(err)
	}
	ss, _ = client.LrangeAll("lio_list2")
	if len(ss) != 2 || ss[0] != "2" || ss[1] != "3" {
		t.Error("err")
	}
	client.Del("lio_list2")
}

func TestPipelineCommon(t *testing.T) {
	/* prepare */
	client.Set("lio_pc1", "1")
	client.Set("lio_pc2", "1")
	client.Set("lio_pc3", "1")

	c, _ := client.GetPipeline()

	/* pipeline */
	client.PipelineDel(c, "lio_pc1", "lio_pc2")
	client.PipelineExpire(c, "lio_pc3", uint32(1))
	err := client.ExecutePipeline(c)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)

	/* confirm */
	b1, _ := client.Exist("lio_pc1")
	b2, _ := client.Exist("lio_pc2")
	b3, _ := client.Exist("lio_pc3")
	if b1 || b2 || b3 {
		t.Error("error")
	}
}

func TestPipelineString(t *testing.T) {
	/* prepare */

	c, _ := client.GetPipeline()

	/* pipeline */
	client.PipelineIncr(c, "lio_ps1")
	client.PipelineIncr(c, "lio_ps1")
	client.PipelineIncr(c, "lio_ps1")
	client.PipelineSet(c, "lio_ps2", "a")
	client.PipelineSetStringAndExpire(c, "lio_ps3", "b", uint32(2))

	err := client.ExecutePipeline(c)
	if err != nil {
		t.Error(err)
	}

	/* confirm */
	b1, _ := client.Get("lio_ps1")
	b2, _ := client.Get("lio_ps2")
	b3, _ := client.Get("lio_ps3")
	if b1 != "3" || b2 != "a" || b3 != "b" {
		t.Error("error")
	}
	time.Sleep(2 * time.Second)
	b4, _ := client.Exist("lio_ps3")
	if b4 {
		t.Error("error")
	}

	/* finally */
	client.Del("lio_ps1", "lio_ps2", "lio_ps3")
}

type HashTest struct {
	UserId   string
	UserName string
}

func TestPipelineHash(t *testing.T) {
	/* prepare */
	c, _ := client.GetPipeline()

	/* pipeline */
	client.PipelineHmset(c, "angio_Hmset", &HashTest{UserId: "123", UserName: "angio_123"})
	client.PipelineHset(c, "angio_Hset", "angio_Hset_1", "angio")
	client.PipelineHincr(c, "angio_Hincr", "angio_Hincr_1")
	client.PipelineHincr(c, "angio_Hincr", "angio_Hincr_1")

	err := client.ExecutePipeline(c)
	if err != nil {
		t.Error(err)
	}

	/* confirm */
	ht := &HashTest{}
	client.HgetAllStruction("angio_Hmset", ht)
	b1, _ := client.Hget("angio_Hset", "angio_Hset_1")
	b2, _ := client.Hget("angio_Hincr", "angio_Hincr_1")
	if ht.UserId != "123" || ht.UserName != "angio_123" || b1 != "angio" || b2 != "2" {
		t.Error("error")
	}
	/* finally */
	client.Del("angio_Hmset", "angio_Hset", "angio_Hincr")
}
func TestPipelineList(t *testing.T) {
	/* prepare */
	c, _ := client.GetPipeline()

	/* pipeline */

	client.PipelineLpushFromHead(c, "angio", "a")
	client.PipelineLpushFromTail(c, "angio", "b")
	client.PipelineLpushFromHead(c, "angio", "c")
	client.PipelineLpushFromTail(c, "angio", "d")
	client.PipelineLrem(c, "angio", "a")
	err := client.ExecutePipeline(c)
	if err != nil {
		t.Error(err)
	}

	/* confirm */
	b1, _ := client.LrangeAll("angio")
	if len(b1) != 3 || b1[0] != "c" || b1[1] != "b" || b1[2] != "d" {
		t.Error("error")
	}
	/* finally */
	client.Del("angio")
}

type Person struct {
	UserId   string
	UserName string
}

func TestPipelineSet(t *testing.T) {
	/* prepare */
	c, _ := client.GetPipeline()

	/* pipeline */
	client.PipelineSadd(c, "angio", "a", "b", "c")

	err := client.ExecutePipeline(c)
	/* confirm */
	b1, _ := client.Slen("angio")
	b2, _ := client.Smembers("angio")
	if b1 != 3 || (b2[0] != "a" && b2[1] != "a" && b2[2] != "a") {
		t.Error("error")
	}

	if err != nil {
		t.Error(err)
	}
	/* pipeline */
	client.PipelineSrem(c, "angio", "a", "b", "c")

	err = client.ExecutePipeline(c)
	/* confirm */
	b3, _ := client.Slen("angio")
	if b3 != 0 {
		t.Error("error")
	}
	/* finally */
	client.Del("angio")
}
func TestPipelineZSet(t *testing.T) {
	/* prepare */
	c, _ := client.GetPipeline()

	/* pipeline */

	//	PipelineZrem(c RedisResourceConn, key string, value ...interface{})

	client.PipelineZincrScore(c, "angio_zset", 10, "aaa")
	client.PipelineZincrScore(c, "angio_zset", 5, "bbb")
	client.PipelineZincrScore(c, "angio_zset", 2, "ccc")
	client.ExecutePipeline(c)

	/* confirm */
	length, _ := client.Zlen("angio_zset")
	s2, _ := client.Zscore("angio_zset", "aaa")
	s3, _ := client.Zscore("angio_zset", "bbb")
	s4, _ := client.Zscore("angio_zset", "ccc")
	if length != 3 || s2 != 10 || s3 != 5 || s4 != 2 {
		t.Error("error")
	}

	/* pipeline */

	client.PipelineZrem(c, "angio_zset", "aaa")
	client.PipelineZrem(c, "angio_zset", "bbb")
	client.ExecutePipeline(c)

	/* confirm */
	s5, _ := client.ZrangeAll("angio_zset")

	if len(s5) != 1 || s5[0] != "ccc" {
		t.Error("error")
	}

	/* finally */
	client.Del("angio_zset")
}

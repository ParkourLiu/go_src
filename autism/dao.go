package main

import (
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
	"strconv"

	"github.com/golang/go/src/pkg/errors"
)

//根据id查询星星
func s_starById(client mysql.MysqlClient, autism *Autism) (map[string]string, error) {
	request, err := client.SearchOneRow(&mysql.Stmt{Sql: "SELECT * FROM `star` WHERE `st_id`=?", Args: []interface{}{autism.St_id}})
	return request, err
}

//根据星星id查询星星评论数量
func s_commentCount(client mysql.MysqlClient, autism *Autism) (string, error) {
	commentCountMap, err := client.SearchOneRow(&mysql.Stmt{Sql: "SELECT count(1) as count FROM `comments` WHERE st_id=? order by createTime desc", Args: []interface{}{autism.St_id}})

	return commentCountMap["count"], err
}

//根据星星id查询星星评论 新方式分页
func s_commentsByStid(client mysql.MysqlClient, autism *Autism) ([]map[string]string, error) {
	requestList, err := client.SearchMutiRows(&mysql.Stmt{Sql: "SELECT * FROM comments WHERE co_id>? AND st_id=? ORDER BY co_id ASC LIMIT 10;", Args: []interface{}{autism.LastCoid, autism.St_id}})
	return requestList, err
}

//捐款列表
func s_donation(client mysql.MysqlClient) ([]map[string]string, error) {
	requestList, err := client.SearchMutiRows(&mysql.Stmt{Sql: "SELECT * FROM `donation` ORDER BY `payMoney`", Args: []interface{}{}})
	return requestList, err
}

//添加评论
func i_comment(client mysql.MysqlClient, autism *Autism) error {
	return mysqlClient.Execute(&mysql.Stmt{Sql: "INSERT INTO `comments` (`co_id`,`userId`,`st_id`,`comment`,`createTime`)VALUES(?,?,?,?,NOW())", Args: []interface{}{autism.LastCoid, autism.ClickUserId, autism.St_id, autism.Comment}})
}

type AutismDao struct {
}

func (*AutismDao) StarList(redisClient redis.RedisClient, client mysql.MysqlClient, autism *Autism) (string, string, []map[string]string, []map[string]string, error) {
	//总点亮数
	brightCount0, err := redisClient.Get(serviceName + "brightCount0") //假数据
	if err != nil {
		return "0", "0", nil, nil, err
	}
	brightCount, err := redisClient.Zlen(serviceName + ":brightCount1") //真数据
	count0 := int(brightCount);
	if err != nil {
		return "0", "0", nil, nil, err
	}
	count1, err := strconv.Atoi(brightCount0)
	count := strconv.Itoa(count0 + count1)

	//总点赞数
	likeCount0, err := redisClient.Get(serviceName + ":likeCount0") //假数据
	if err != nil {
		return count, "0", nil, nil, err
	}
	likeCount, err := redisClient.Get(serviceName + ":likeCount1") //真数据
	if err != nil {
		return count, "0", nil, nil, err
	}
	likeCount = likeCount + likeCount0
	//首页弹幕列表
	sql := "SELECT photo,comment  FROM `rolls` ORDER BY RAND() LIMIT 100"
	rolls, err := client.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return count, likeCount, nil, nil, err
	}
	//星星
	sql = "SELECT st_id,st_name,st_head,st_type FROM `star` ORDER BY RAND() LIMIT 30"
	starList, err := client.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return count, likeCount, nil, rolls, err
	}
	return count, likeCount, rolls, starList, err
}
func (*AutismDao) Likes(redisClient redis.RedisClient, client mysql.MysqlClient, autism *Autism) error {
	//为单个人点赞用户id存储
	b, err := redisClient.Sismember(serviceName+":likeFlag:"+autism.St_id, autism.User_id)
	if b {
		return errors.New("你已经点过赞了")
	}
	err = redisClient.Sadd(serviceName+":likeFlag:"+autism.St_id, autism.User_id)
	if err != nil {
		return err
	}
	//总点赞数
	redisClient.Incr(serviceName + ":likeCount1")
	return err
}

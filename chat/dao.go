package main

import (
	"fmt"
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
)

type Dao struct {
}

//根据userid查询redis数据库里是否有token
func (*Dao) SelectRyToken(client redis.RedisClient, key string) (i string, err error) {
	i, err = client.Get(key)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return "", err
		}
	}
	return i, nil
}

//一次添加两个key
func (*Dao) Add(client redis.RedisClient, param [][]interface{}) (i interface{}, err error) {
	var c redis.RedisResourceConn
	c, err = client.GetConn()
	defer client.ReturnConn(c)
	if err != nil {
		return
	}
	c.Send("MULTI")
	for i := 0; i < len(param); i++ {
		c.Send("set", param[i][0], param[i][1])
	}
	//exec pipeline

	i, err = c.Do("EXEC")
	return
}

//批量读取用户信息
func (*Dao) ListPipelineGetHashMap_User(client redis.RedisClient, keys []string) (values []interface{}, err error) {
	var c redis.RedisResourceConn
	result := []interface{}{}
	c, err = client.GetConn()
	defer client.ReturnConn(c)
	if err != nil {
		return
	}
	for _, key := range keys {
		c.Send("HGETALL", userId+key)
	}
	c.Flush()
	//c.Receive()
	for i := 0; i < len(keys); i++ {
		values, err := c.Receive()
		if err != nil {
			return nil, err
		}
		result = append(result, values)
	}
	return result, err
}

func (*Dao) ListPipelineIntStatus(client redis.RedisClient, key string, members []string) (values []string, err error) {
	var c redis.RedisResourceConn
	var result1 []string
	c, err = client.GetConn()
	defer client.ReturnConn(c)
	if err != nil {
		return
	}
	for i := 0; i < len(members); i++ {
		c.Send("ZSCORE", key, members[i])
	}

	c.Flush()
	//c.Receive()
	for i := 0; i < len(members); i++ {
		values, err := c.Receive()
		switch v1 := values.(type) {
		case []uint8:
			result1 = append(result1, "1")
		case nil:
			result1 = append(result1, "0")
		default:
			fmt.Println("result1===", result1, "err==", err, "v1==", v1)
			return nil, err
		}
	}
	return result1, err
}
func (*Dao) ListZrange2(client redis.RedisClient, key string, startIndex int64, endIndex int64) (value []string, err error) {
	value, err = client.Zrange2(key, startIndex, endIndex)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return nil, err
		}
	}
	return value, nil
}
func (*Dao) ListZrange(client redis.RedisClient, key string, startIndex int64, endIndex int64) (value []string, err error) {
	value, err = client.Zrange(key, startIndex, endIndex)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return nil, err
		}
	}
	return value, nil
}
func (*Dao) ListZRANK(client redis.RedisClient, key, member string) (values int64, err error) {
	values, err = client.ZRANK(key, member)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return 0, err
		}
	}
	return values, nil
}

//批量读取群聊的信息
func (*Dao) ListPipelineGetHashMap_GroupChat(client redis.RedisClient, keys []string) (values []interface{}, err error) {
	var c redis.RedisResourceConn
	result := []interface{}{}
	c, err = client.GetConn()
	defer client.ReturnConn(c)
	if err != nil {
		return
	}
	for _, key := range keys {
		c.Send("HGETALL", groupChatInfo+key)
	}
	c.Flush()
	//c.Receive()
	for i := 0; i < len(keys); i++ {
		values, err := c.Receive()
		if err != nil {
			return nil, err
		}
		result = append(result, values)
	}
	return result, err
}

//社群成员列表添加一个用户，  用户参与的群聊添加一个群聊id
func (*Dao) DaoJoin(client redis.RedisClient, param [][]interface{}) (i interface{}, err error) {
	var c redis.RedisResourceConn
	c, err = client.GetConn()
	defer client.ReturnConn(c)
	if err != nil {
		return
	}
	c.Send("MULTI")
	for i := 0; i < len(param); i++ {
		c.Send("Zadd", param[i][0], param[i][1], param[i][2])
	}
	//exec pipeline

	i, err = c.Do("EXEC")
	return
}

//批量删除
func (*Dao) ListPipelineDel(client redis.RedisClient, keys []string) (values []interface{}, err error) {
	var c redis.RedisResourceConn
	result := []interface{}{}
	c, err = client.GetConn()
	defer client.ReturnConn(c)
	if err != nil {
		return
	}
	for _, key := range keys {
		c.Send("DEL", key)
	}
	c.Flush()
	//c.Receive()
	for i := 0; i < len(keys); i++ {
		values, err := c.Receive()
		if err != nil {
			return nil, err
		}
		result = append(result, values)
	}
	return result, err
}

func (*Dao) GetGroupChatId(client mysql.MysqlClient, acId string) ([]map[string]string, error) {
	return client.SearchMutiRows(&mysql.Stmt{Sql: "select groupChatId,groupChatName from groupchat where ad='a' AND acId=? order by  createTime Desc", Args: []interface{}{acId}})
}

func (*Dao) GetActivityInfo(client redis.RedisClient, acId string) (map[string]string, error) {
	value, err := client.HgetAllMap("activity:" + acId)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return nil, err
		}
	}
	return value, nil
}

func (*Dao) Zrem(client redis.RedisClient, key string, member string) error {
	return client.Zrem(key, member)
}

func GetChatInfo(chat *GroupChatInfo) (map[string]string,error) {
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql:"select * from `groupchat` where acId=? and `ad`='a'limit 1",Args:[]interface{}{chat.AcId}})
}
func GetClassId(chat *GroupChatInfo)(map[string]string,error){
	return mysqlClient.SearchOneRow(&mysql.Stmt{Sql:"select * from `groupchat` where `groupChatId`=? and `ad`='a'",Args:[]interface{}{chat.GroupChatId}})
}
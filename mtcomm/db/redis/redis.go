package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/youtube/vitess/go/pools"
	"golang.org/x/net/context"

	logger "mtcomm/log"
)

type RedisClient interface {
	/* common */
	Expire(key string, expireSecond uint32) (err error)
	Exist(key string) (exist bool, err error)
	Del(keys ...interface{}) (err error)
	Keys(key string) (values []string, err error)
	//	Keys(key string) (keys []string, err error)
	GetConn() (RedisResourceConn, error)
	ReturnConn(c RedisResourceConn)
	/* string */
	Incr(key string) (value int64, err error)
	Decr(key string) (value int64, err error)
	Get(key string) (value string, err error)
	GetStrings(keys ...string) (values []string, err error)
	Set(key string, value string) (err error)
	SetStringAndExpire(key string, value string, expireSecond uint32) (err error) //expireSecond:秒
	//设置缓存，默认1小时过期
	SetStringWithDefaultExpire(key string, value string) (err error)
	SETNX(key string, value string) (ex int64, err error) //如果没有此key则插入并返回1，存在则不插入返回0

	/* hash */
	Hmset(key string, value interface{}) (err error)
	HgetAllStruction(key string, value interface{}) (err error)
	HgetAllMap(key string) (value map[string]string, err error)
	Hset(key string, subkey string, value string) (err error)
	Hget(key string, subkey string) (value string, err error)
	Hincr(key string, subkey string) (value int64, err error) // 加1
	Hdecr(key string, subkey string) (value int64, err error)

	/* list */
	LpushFromHead(key string, value ...string) (err error)
	LpushFromTail(key string, value ...string) (err error)
	//返回值包含startIndex和endIndex的值。
	Lrange(key string, startIndex int64, endIndex int64) (value []string, err error)
	LrangeAll(key string) (value []string, err error)
	Llen(key string) (value int64, err error)
	Lrem(key string, value string) (err error)

	Lpop(key string) (value string, err error)
	Rpush(key string, value string) (err error)

	/* set */
	Sismember(key string, value string) (exist bool, err error)
	Smembers(key string) (value []string, err error)
	Srem(key string, value ...interface{}) (err error)
	Sadd(key string, value ...interface{}) (err error)
	Slen(key string) (value int64, err error)

	/* zset */
	//分值高到低排序，返回值包含startIndex和endIndex的值。
	Zrange(key string, startIndex int64, endIndex int64) (value []string, err error)
	//分值低到高排序，返回值包含startIndex和endIndex的值。
	Zrange2(key string, startIndex int64, endIndex int64) (value []string, err error)

	ZrangeWithscores(key string, startIndex int64, endIndex int64) (value []string, err error) ////分值低到高排序,并且同步返回此值得score
	ZrangeAll(key string) (value []string, err error)
	Zadd(key string, score int64, member string) (i int, err error)
	ZCARD(key string) (value int64, err error)
	ZRANK(key, members string) (i int64, err error)
	//	Zadd(key string, members map[string]int64) (err error)
	ZincrScore(key string, score int, member string) (err error)
	Zscore(key string, member string) (score int64, err error) //取得分值
	Zlen(key string) (value int64, err error)
	//	Zismember(key string, value string) (exist bool, err error)
	Zrem(key string, value ...interface{}) (err error)
	ZRANGEBYSCORE(key string, score int64) (value []string, err error)                 //根据分值区间，截取数据
	ZRANGEBYSCORE2(key string, score1 int64, score2 int64) (value []string, err error) //根据分值区间，分值开始，分值结束，截取数据
	Zremrangebyrank(key string, score1 int64, score2 int64) (count int64, err error)   //根据区间删除zset中元素，清空传0，-1
	/* *************pipeline*********** */
	GetPipeline() (c RedisResourceConn, err error)
	ExecutePipeline(c RedisResourceConn) (err error)
	/* common */
	PipelineExpire(c RedisResourceConn, key string, expireSecond uint32)
	PipelineDel(c RedisResourceConn, keys ...interface{})
	/* string */
	PipelineIncr(c RedisResourceConn, key string)
	PipelineSet(c RedisResourceConn, key string, value string)
	PipelineSetStringAndExpire(c RedisResourceConn, key string, value string, expireSecond uint32)
	PipelineSetStringWithDefaultExpire(c RedisResourceConn, key string, value string)
	/* hash */
	PipelineHmset(c RedisResourceConn, key string, value interface{})
	PipelineHset(c RedisResourceConn, key string, subkey string, value string)
	PipelineHincr(c RedisResourceConn, key string, subkey string) // 加1
	PipelineGetHashMap(keys []string) (values []interface{}, err error)
	PipelineGetHashMap2(keys []string) (values []map[string]string, err error)
	PipelineIntStatus(key string, members []string) (values []string, err error)
	PipelineIntStatus2(key []string, members string) (values []string, err error)
	PipelineIntStatusGroup(key []string, UId string) (values []string, err error)
	/* list */
	PipelineLpushFromHead(c RedisResourceConn, key string, value ...string)
	PipelineLpushFromTail(c RedisResourceConn, key string, value ...string)
	PipelineLrem(c RedisResourceConn, key string, value string)
	/* set */
	PipelineSrem(c RedisResourceConn, key string, value ...interface{})
	PipelineSadd(c RedisResourceConn, key string, value ...interface{})
	/* zset */
	PipelineZincrScore(c RedisResourceConn, key string, score int, member string)
	PipelineZrem(c RedisResourceConn, key string, value ...interface{})
	PipelineZadd(param [][]interface{}) (value interface{}, err error)
	PipelineZrem2(arrays [][]string) (value interface{}, err error)
}

type redisClient struct {
	ctx    context.Context
	pool   *pools.ResourcePool
	logger *logger.Logger
}

type RedisServerInfo struct {
	Ctx           context.Context
	Logger        *logger.Logger
	RedisHost     string
	RedisPassword string
	PoolCap       int
	PoolMaxCap    int
}

func NewRedisClient(info *RedisServerInfo) RedisClient {
	if info.PoolCap == 0 || info.PoolMaxCap == 0 {
		info.PoolCap = 100
		info.PoolMaxCap = 200
	}
	p := pools.NewResourcePool(func() (pools.Resource, error) {
		var c redis.Conn
		var err error
		if info.RedisPassword != "" {
			c, err = redis.Dial("tcp", info.RedisHost, redis.DialPassword(info.RedisPassword))
		} else {
			c, err = redis.Dial("tcp", info.RedisHost)
		}
		return RedisResourceConn{c}, err
	}, info.PoolCap, info.PoolMaxCap, 10*time.Second)
	return &redisClient{
		ctx:    info.Ctx,
		pool:   p,
		logger: info.Logger,
	}
}

// RedisResourceConn adapts a Redigo connection to a Vitess Resource.
type RedisResourceConn struct {
	redis.Conn
}

func (r RedisResourceConn) Close() {
	r.Conn.Close()
}

func (rc *redisClient) GetConn() (RedisResourceConn, error) {
	r, err := rc.pool.Get(rc.ctx)
	if err != nil {
		rc.logger.Error("status", "fail", "msg", "redis get conn fail", "detail", err.Error())
		return RedisResourceConn{}, err
	}
	c := r.(RedisResourceConn)
	return c, nil
}

func (rc *redisClient) ReturnConn(c RedisResourceConn) {
	rc.pool.Put(c)
}

func (rc *redisClient) LpushFromHead(key string, value ...string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("LPUSH", redis.Args{}.Add(key).AddFlat(value)...)
	return
}
func (rc *redisClient) LpushFromTail(key string, value ...string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("RPUSH", redis.Args{}.Add(key).AddFlat(value)...)
	return
}

func (rc *redisClient) Lrange(key string, startIndex int64, endIndex int64) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Strings(c.Do("LRANGE", key, startIndex, endIndex))
	return
}

func (rc *redisClient) LrangeAll(key string) (value []string, err error) {
	value, err = rc.Lrange(key, 0, -1)
	return
}

func (rc *redisClient) Llen(key string) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Int64(c.Do("LLEN", key))
	return
}
func (rc *redisClient) Lrem(key string, value string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("LREM", key, 0, value)
	return
}

func (rc *redisClient) Zrange(key string, startIndex int64, endIndex int64) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Strings(c.Do("ZREVRANGE", key, startIndex, endIndex))
	return
}
func (rc *redisClient) Zrange2(key string, startIndex int64, endIndex int64) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Strings(c.Do("ZRANGE", key, startIndex, endIndex))
	return
}

//分值低到高排序,并且同步返回此分值score
func (rc *redisClient) ZrangeWithscores(key string, startIndex int64, endIndex int64) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	value, err = redis.Strings(c.Do("ZRANGE", key, startIndex, endIndex, "withscores"))
	return
}

func (rc *redisClient) Zadd(key string, score int64, member string) (value int, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	value, err = redis.Int(c.Do("Zadd", key, score, member))
	return
}
func (rc *redisClient) ZCARD(key string) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Int64(c.Do("ZCARD", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}
func (rc *redisClient) ZRANK(key, members string) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	value, err = redis.Int64(c.Do("ZRANK", key, members))
	if err == redis.ErrNil {
		err = nil
	}
	return
}
func (rc *redisClient) PipelineZadd(param [][]interface{}) (value interface{}, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	c.Send("MULTI")
	for i := 0; i < len(param); i++ {
		c.Send("Zadd", param[i][0], param[i][1], param[i][2])
	}
	//exec pipeline

	value, err = c.Do("EXEC")
	return
}
func (rc *redisClient) PipelineZrem2(param [][]string) (value interface{}, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	c.Send("MULTI")
	fmt.Println(param)
	for i := 0; i < len(param); i++ {
		c.Send("Zrem", param[i][0], param[i][1])

	}
	//exec pipeline

	value, err = c.Do("EXEC")
	return
}
func (rc *redisClient) ZrangeAll(key string) (value []string, err error) {
	value, err = rc.Zrange(key, int64(0), int64(-1))
	return
}

func (rc *redisClient) ZincrScore(key string, score int, member string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("ZINCRBY", key, score, member)
	if err == redis.ErrNil {
		err = nil
	}
	return
}
func (rc *redisClient) Zscore(key string, member string) (score int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	score, err = redis.Int64(c.Do("ZSCORE", key, member))
	if err == redis.ErrNil {
		err = nil
	}
	return
}
func (rc *redisClient) Zlen(key string) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Int64(c.Do("ZCARD", key))
	return
}

func (rc *redisClient) Zrem(key string, value ...interface{}) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("ZREM", redis.Args{}.Add(key).AddFlat(value)...)
	return
}

func (rc *redisClient) ZRANGEBYSCORE(key string, score int64) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	value, err = redis.Strings(c.Do("ZRANGEBYSCORE", key, score, "+inf"))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) ZRANGEBYSCORE2(key string, score1 int64, score2 int64) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	value, err = redis.Strings(c.Do("ZRANGEBYSCORE", key, score1, score2))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) Zremrangebyrank(key string, score1 int64, score2 int64) (count int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	count, err = redis.Int64(c.Do("zremrangebyrank", key, score1, score2))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) Sismember(key string, value string) (exist bool, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	exist, err = redis.Bool(c.Do("SISMEMBER", key, value))
	return
}

func (rc *redisClient) Smembers(key string) (value []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Strings(c.Do("SMEMBERS", key))
	return
}

func (rc *redisClient) Srem(key string, value ...interface{}) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("SREM", redis.Args{}.Add(key).AddFlat(value)...)
	return
}

func (rc *redisClient) Sadd(key string, value ...interface{}) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("SADD", redis.Args{}.Add(key).AddFlat(value)...)
	return
}

func (rc *redisClient) Slen(key string) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Int64(c.Do("SCARD", key))
	return
}

func (rc *redisClient) Hset(key string, subkey string, value string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("HSET", key, subkey, value)
	return
}

func (rc *redisClient) Hget(key string, subkey string) (value string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.String(c.Do("HGET", key, subkey))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) Hdecr(key string, subkey string) (value int64, err error) {
	return rc.hincr(key, subkey, -1)
}

func (rc *redisClient) Hincr(key string, subkey string) (value int64, err error) {
	return rc.hincr(key, subkey, 1)
}

func (rc *redisClient) hincr(key string, subkey string, count int) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Int64(c.Do("HINCRBY", key, subkey, count))
	return
}

func (rc *redisClient) Hmset(key string, value interface{}) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
	return
}

func (rc *redisClient) HgetAllStruction(key string, value interface{}) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	var v []interface{}
	v, err = redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return
	}
	err = redis.ScanStruct(v, value)
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) HgetAllMap(key string) (value map[string]string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.StringMap(redis.Values(c.Do("HGETALL", key)))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) Set(key string, value string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("SET", key, value)
	return
}

func (rc *redisClient) GetStrings(keys ...string) (values []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	c.Send("MULTI")

	for _, key := range keys {
		c.Send("GET", key)
	}

	//exec pipeline
	values, err = redis.Strings(c.Do("EXEC"))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) Keys(key string) (keys []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	keys, err = redis.Strings(c.Do("KEYS", key))
	return
}

func (rc *redisClient) Del(keys ...interface{}) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("DEL", keys...)
	return
}

func (rc *redisClient) Exist(key string) (exist bool, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	exist, err = redis.Bool(c.Do("EXISTS", key))
	return
}

func (rc *redisClient) Get(key string) (value string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.String(c.Do("GET", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (rc *redisClient) SetStringAndExpire(key string, value string, expireSecond uint32) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("SETEX", key, expireSecond, value)
	return
}

func (rc *redisClient) SETNX(key string, value string) (ex int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	ex, err = redis.Int64(c.Do("SETNX", key, value))
	return
}

func (rc *redisClient) Lpop(key string) (value string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.String(c.Do("LPOP", key))
	return
}

func (rc *redisClient) Rpush(key string, value string) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("RPUSH", key, value)
	return
}

// 设置过期时间
func (rc *redisClient) Expire(key string, expireSecond uint32) (err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	_, err = c.Do("EXPIRE", key, expireSecond)
	return
}

func (rc *redisClient) Incr(key string) (value int64, err error) {
	return rc.incr(key, 1)
}

func (rc *redisClient) Decr(key string) (value int64, err error) {
	return rc.incr(key, -1)
}

func (rc *redisClient) incr(key string, count int) (value int64, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}

	value, err = redis.Int64(c.Do("INCRBY", key, count))
	return
}

func (rc *redisClient) SetStringWithDefaultExpire(key string, value string) (err error) {
	return rc.SetStringAndExpire(key, value, 3600)
}

/* pipeline */
func (rc *redisClient) GetPipeline() (c RedisResourceConn, err error) {
	c, err = rc.GetConn()
	if err != nil {
		return
	}

	c.Send("MULTI")
	return
}

func (rc *redisClient) ExecutePipeline(c RedisResourceConn) (err error) {
	defer rc.pool.Put(c)
	//exec pipeline
	_, err = c.Do("EXEC")
	return
}

/* common */
func (rc *redisClient) PipelineExpire(c RedisResourceConn, key string, expireSecond uint32) {
	c.Send("EXPIRE", key, expireSecond)
}

func (rc *redisClient) PipelineDel(c RedisResourceConn, keys ...interface{}) {
	c.Send("DEL", keys...)
}

/* string */
func (rc *redisClient) PipelineIncr(c RedisResourceConn, key string) {
	c.Send("INCR", key)
}
func (rc *redisClient) PipelineSet(c RedisResourceConn, key string, value string) {
	c.Send("SET", key, value)
}
func (rc *redisClient) PipelineSetStringAndExpire(c RedisResourceConn, key string, value string, expireSecond uint32) {
	c.Send("SETEX", key, expireSecond, value)
}
func (rc *redisClient) PipelineSetStringWithDefaultExpire(c RedisResourceConn, key string, value string) {
	c.Send("SETEX", key, 3600, value)
}

/* hash */
func (rc *redisClient) PipelineHmset(c RedisResourceConn, key string, value interface{}) {
	c.Send("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
}
func (rc *redisClient) PipelineHset(c RedisResourceConn, key string, subkey string, value string) {
	c.Send("HSET", key, subkey, value)
}
func (rc *redisClient) PipelineHincr(c RedisResourceConn, key string, subkey string) {
	c.Send("HINCRBY", key, subkey, 1)
}

func (rc *redisClient) PipelineGetHashMap(keys []string) (values []interface{}, err error) {
	var c RedisResourceConn
	result := []interface{}{}
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	for _, key := range keys {
		c.Send("HGETALL", key)
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

func (rc *redisClient) PipelineGetHashMap2(keys []string) (values []map[string]string, err error) {
	var c RedisResourceConn
	result := make([]map[string]string, 0)
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	for _, key := range keys {
		c.Send("HGETALL", key)
	}
	c.Flush()
	//c.Receive()
	for i := 0; i < len(keys); i++ {
		values, err := c.Receive()
		if err != nil {
			return nil, err
		}
		v, err1 := redis.StringMap(redis.Values(values, err))
		if err1 != nil {
			return nil, err1
		}
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result, err
}

func (rc *redisClient) PipelineIntStatus(key string, members []string) (values []string, err error) {
	var c RedisResourceConn
	var result1 []string
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
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

//个人关注个人列 表
func (rc *redisClient) PipelineIntStatus2(key []string, members string) (values []string, err error) {
	var c RedisResourceConn
	var result1 []string
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	for i := 0; i < len(key); i++ {
		c.Send("ZSCORE", key[i], members)
	}

	c.Flush()
	//c.Receive()
	for i := 0; i < len(key); i++ {
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

//个人关注社群列 表状态
func (rc *redisClient) PipelineIntStatusGroup(key []string, UId string) (values []string, err error) {
	var c RedisResourceConn
	c, err = rc.GetConn()
	defer rc.pool.Put(c)
	if err != nil {
		return
	}
	for i := 0; i < len(key); i++ {
		c.Send("HMGET", key[i], "gId")
	}

	c.Flush()
	//c.Receive()
	param2Slice := []string{}
	for i := 0; i < len(key); i++ {
		values, err := c.Receive()
		switch v := values.(type) {
		case []interface{}:
			for _, pa := range v {
				switch v1 := pa.(type) {
				case []uint8:
					strV11 := string(v1)
					if strV11 == UId { //0表示是群主
						param2Slice = append(param2Slice, "0")
					} else {
						param2Slice = append(param2Slice, "1")
					}
				case nil:
					param2Slice = append(param2Slice, "1")
					fmt.Println("values===", values, "err=====", err, fmt.Sprintf("%T", v1))
				default:
					panic("params type not supported")
				}
			}
		default:
			panic("params type not supported")
		}
	}
	return param2Slice, err
}

/* list */
func (rc *redisClient) PipelineLpushFromHead(c RedisResourceConn, key string, value ...string) {
	c.Send("LPUSH", redis.Args{}.Add(key).AddFlat(value)...)
}
func (rc *redisClient) PipelineLpushFromTail(c RedisResourceConn, key string, value ...string) {
	c.Send("RPUSH", redis.Args{}.Add(key).AddFlat(value)...)
}
func (rc *redisClient) PipelineLrem(c RedisResourceConn, key string, value string) {
	c.Send("LREM", key, 0, value)
}

/* set */
func (rc *redisClient) PipelineSrem(c RedisResourceConn, key string, value ...interface{}) {
	c.Send("SREM", redis.Args{}.Add(key).AddFlat(value)...)
}
func (rc *redisClient) PipelineSadd(c RedisResourceConn, key string, value ...interface{}) {
	c.Send("SADD", redis.Args{}.Add(key).AddFlat(value)...)
}

/* zset */
func (rc *redisClient) PipelineZincrScore(c RedisResourceConn, key string, score int, member string) {
	c.Send("ZINCRBY", key, score, member)
}
func (rc *redisClient) PipelineZrem(c RedisResourceConn, key string, value ...interface{}) {
	c.Send("ZREM", redis.Args{}.Add(key).AddFlat(value)...)
}

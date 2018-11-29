package monitor

import (
	//	"mtcomm/db/mongo"
	"encoding/json"
	"mtcomm/db/mysql"
	"mtcomm/db/redis"
	"net/http"
)

type MonitorHandler struct {
	// redis
	RedisClient redis.RedisClient
	// mysql
	MySqlTableName string
	MysqlClient    mysql.MysqlClient
	//	// mongo
	//	MongoDbName string
	//	MongoCollection string
	//	MongoClient    mongo.mongoClient
}

type Response struct {
	Err string `json:"err"`
}

func (h *MonitorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	noErrStr := "{" + "\"err\":\"\"" + "}"
	//recover panic
	defer func() {
		if x := recover(); x != nil {
			w.WriteHeader(500)
			w.Write([]byte("NG"))
			return
		}
	}()

	if h.RedisClient != nil {
		_, err = h.RedisClient.Incr("monitor_svc")
	}
	if h.MySqlTableName != "" && h.MysqlClient != nil {
		_, err = h.MysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from " + h.MySqlTableName + " limit 1", Args: []interface{}{}})
	}
	//	if h.MongoDbName != "" && h.MongoCollection != "" && h.MongoClient != nil {
	//		_, err3 := h.MongoClient.SearchById(xxx)
	//	}

	if err != nil {
		w.WriteHeader(500)
		response := Response{err.Error()}
		JsonResponse, _ := json.Marshal(response)
		w.Write(JsonResponse)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(noErrStr))
	}
}

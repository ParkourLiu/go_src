package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"
)

type Oss interface {
	GetOssTokenForWeb() (map[string]interface{}, error)
}

type oss struct{}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func (service oss) GetOssTokenForWeb() (map[string]interface{}, error) { //返回值，最后一页标识，code，err
	log.Debug("method_start", "GetOssTokenForWeb", "input", "")
	response, err := get_policy_token()
	if err != nil {
		return nil, err
	}
	tokenMap, err1 := jsonMap(response)
	if err1 != nil {
		return nil, err1
	}
	tokenMap["imageIdPrefix"] = idGenClient.GetUniqueId()
	log.Debug("████████OssToken:", fmt.Sprint(tokenMap))
	log.Debug("method_end", "GetOssTokenForWeb", "status", "success")
	return tokenMap, nil
}

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func get_policy_token() (string, error) {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = get_gmt_iso8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = upload_dir
	policyToken.Policy = string(debyte)
	response, err := json.Marshal(policyToken)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

//json转map
func jsonMap(str string) (map[string]interface{}, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(str), &dat); err == nil {
		return dat, nil
	} else {
		return make(map[string]interface{}), err
	}
}

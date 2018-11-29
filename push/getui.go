// getui
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
//个推demo测试工具配置参数
//	appid        string = "ChPdvWd91n5GCbbzP1YuP8"
//	appkey       string = "iCWqHYbMw18cYY5NluwWp2"
//	mastersecret string = "UxlvgvFh57AhhJehtHI0G7"
//	appsecret    string = "4SLMFG6GTM8IM8cldMI8e4"

//钦家测试app个推配置参数
//	appid        string = "XmsnHx8jwe6CvVKhnOzV22"
//	appkey       string = "pi1tFGmiQd77vb0rBc4DX4"
//	mastersecret string = "FbHuAZL8Iu8Krh9tw1Wvq6"
//	appsecret    string = "xHek8i1q2394ig3PYbWl3"

//钦家生产app个推配置参数
//	appid        string = "sV0NB8kGTx73YJokpdAyQA"
//	appkey       string = "V12CZrm7je6VI26xUMZyf3"
//	mastersecret string = "lpCswXNUUhAP61WpUGjuo4"
//	appsecret    string = "DTKvVhEek16KQdZc9kiez4"

)

//json转map
func jsonMap(str string) (map[string]interface{}, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(str), &dat); err == nil {
		return dat, nil
	} else {
		return make(map[string]interface{}), err
	}
}

//结构体转Json字符串
func structJson(pa interface{}) (string, error) {
	jsons, err := json.Marshal(pa) //转换成JSON返回的是byte[]
	if err != nil {
		return "", err
	}
	log.Debug(string(jsons)) //byte[]转换成string 输出
	return string(jsons), nil
}

//sig需要的sha256算法
func getSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//获取毫秒时间
func nowTime() string {
	return fmt.Sprint(time.Now().UnixNano() / 1000000)
}

type authtokenStruct struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Appkey    string `json:"appkey"`
}

//获取authtoken=====================================================================================start
func authtoken() (string, error) {
	timestamp := fmt.Sprint(nowTime())
	sign := getSha256(appkey + timestamp + mastersecret)

	url := authtoken_URL
	//json序列化
	po := authtokenStruct{
		Sign:      sign,
		Timestamp: timestamp,
		Appkey:    appkey,
	}
	//结构体转Json字符串
	post, err := structJson(po)
	if err != nil {
		return "", err
	}
	log.Debug(url, "post", post)

	var jsonStr = []byte(post)
	log.Debug("jsonStr", jsonStr)
	log.Debug("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close() // ok，即使不读取Body中的数据，即使Body是空的，也要调用close方法,一定要判断resp是否为空
	}
	if err != nil {
		return "", err
	}

	log.Debug("response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	str := string(body)
	log.Debug("response Body:", string(body))

	//把返回回来的json，提取authtoken
	authtokenMap, err := jsonMap(str)
	log.Debug("▼▼▼▼▼authtoken,authtokenMap:::::::::", fmt.Sprint(authtokenMap))
	if err != nil {
		return "", err
	}
	result := authtokenMap["result"].(string)
	log.Debug("██获取authtoken███result:", result)

	return authtokenMap["auth_token"].(string), nil
}

//获取authtoken=====================================================================================end

//推送通知结构体===============================================0
//ios特殊处理结构体========0
type push_info struct {
	Aps aps `json:"aps"`
}

type aps struct {
	Alert   alert   `json:"alert"`
	IosInfo iosInfo `json:"iosInfo"`
	Badge   int     `json:"badge"`
}

type alert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type iosInfo struct {
	Json string `json:"json"`
}

//ios特殊处理结构体=========1

type styleStruct struct {
	Type  int    `json:"type"`
	Text  string `json:"text"`
	Title string `json:"title"`
}

type notificationStruct struct {
	Style                styleStruct `json:"style"`
	Transmission_type    bool        `json:"transmission_type"`
	Transmission_content string      `json:"transmission_content"`
}

type messageStruct struct {
	Appkey              string `json:"appkey"`
	Is_offline          bool   `json:"is_offline"`
	Offline_expire_time int64  `json:"offline_expire_time"`
	Msgtype             string `json:"msgtype"`
}

type pushListStruct struct {
	Message      messageStruct      `json:"message"`
	Notification notificationStruct `json:"notification"`
	Push_info    push_info          `json:"push_info"`
	Alias        string             `json:"alias"`
	Requestid    string             `json:"requestid"`
}

type AllPush struct {
	Message      messageStruct      `json:"message"`
	Notification notificationStruct `json:"notification"`
	Push_info    push_info          `json:"push_info"`
	Requestid    string             `json:"requestid"`
}

//推送通知结构体===============================================1

//推送通知结构体封装===============================================0
func newAndroidParamJson(appkey string, text string, title string, alias string, requestid string, jsonInfo string) (string, error) {
	alert := alert{
		Title: title,
		Body:  text,
	}
	iosInfo := iosInfo{
		Json: jsonInfo, //推送通知的时候的隐藏参数
	}
	aps := aps{
		Alert:   alert,
		IosInfo: iosInfo,
		Badge:   1,
	}
	push_info := push_info{
		Aps: aps,
	}
	style := styleStruct{
		Type:  0,
		Text:  text,
		Title: title,
	}
	notification := notificationStruct{
		Style:                style,
		Transmission_type:    false,
		Transmission_content: jsonInfo, //推送通知的时候的隐藏参数
	}
	message := messageStruct{
		Appkey:              appkey,
		Is_offline:          true,
		Offline_expire_time: 10000000,
		Msgtype:             "notification",
	}
	pushList := pushListStruct{
		Message:      message,
		Notification: notification,
		Push_info:    push_info,
		Alias:        alias,
		Requestid:    requestid,
	}
	return structJson(pushList)
}

//推送通知结构体封装===============================================1

//个推，群推传入参数map{text,title,alias,terrace}=====================================================0
func PushList(paramsStruct *pushInfoList) (string, error) {
	url := List_push_URL

	authtoken, err := authtoken()
	if err != nil {
		log.Debug("authtoken:", err)
		return "", err
	}

	params := paramsStruct.PushInfos
	//遍历参数
	for i := 0; i < len(params); i++ {
		if params[i].Text != "" && params[i].Title != "" && params[i].Alias != "" {
			androidPost, err := newAndroidParamJson(appkey, params[i].Text, params[i].Title, params[i].Alias, nowTime(), params[i].JsonInfo)
			if err != nil {
				return "", err
			}
			//调用推送========================
			//json序列化
			post := androidPost
			log.Debug("████████推送通知URL:", url, "████████推送通知post", post)

			var jsonStr = []byte(post)
			log.Debug("LIST_push_jsonStr", jsonStr)
			log.Debug("LIST_push_new_str", bytes.NewBuffer(jsonStr))

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			if err != nil {
				log.Debug("http.NewRequest:", err)
				return "", err
			}
			//添加请求头文件
			req.Header.Set("Content-Type", "application/json")
			headMap := req.Header
			headMap["authtoken"] = append(headMap["authtoken"], authtoken)
			//req.Header.Add("authtoken", authtoken)
			log.Debug("authtoken:", req)
			client := &http.Client{}
			resp, err := client.Do(req)
			if resp != nil {
				defer resp.Body.Close() // ok，即使不读取Body中的数据，即使Body是空的，也要调用close方法,一定要判断resp是否为空
			}
			if err != nil {
				log.Debug("client.Do:", err)
				return "", err
			}
			//调用完成开始处理返回数据================================================================
			log.Debug("response Status:", resp.Status)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Debug("ioutil.ReadAll:", err)
				return "", err
			}
			str := string(body)
			log.Debug("response Body:", string(body))

			//把返回回来的json，提取authtoken
			authtokenMap, err := jsonMap(str)
			log.Debug("▼▼▼▼▼推送,authtokenMap:::::::::", fmt.Sprint(authtokenMap))
			if err != nil {
				return "", err
			}
			result := authtokenMap["result"].(string)
			log.Debug("██推送通知"+fmt.Sprint(i)+"███result:", result)
		}
	}

	return "", nil
}

//个推，群推传入参数map{text,title,alias,terrace}=====================================================1

//透传消息结构体===============================================0
type iosPushList struct {
	Message      messageStruct `json:"message"`
	Transmission transmission  `json:"transmission"`
	Alias        string        `json:"alias"`
	Requestid    string        `json:"requestid"`
}

type transmission struct {
	Transmission_content string `json:"transmission_content"`
}

//透传消息结构体===============================================1

//透传消息封装======================================================================0
func OSPFJson(appkey string, text string, title string, alias string, requestid string) (string, error) {

	transmission := transmission{
		Transmission_content: text,
	}
	message := messageStruct{
		Appkey:              appkey,
		Is_offline:          true,
		Offline_expire_time: 10000000,
		Msgtype:             "transmission",
	}
	iosPushList := iosPushList{
		Message:      message,
		Transmission: transmission,
		Alias:        alias,
		Requestid:    requestid,
	}
	return structJson(iosPushList)
}

//透传消息封装======================================================================1

//透传消息传入参数map{text,title,alias,terrace}
func OSPFList(paramsStruct *pushInfoList) (string, error) {
	url := List_push_URL

	authtoken, err := authtoken()
	if err != nil {
		log.Debug("authtoken:", err)
		return "", err
	}

	params := paramsStruct.PushInfos
	//遍历参数
	for i := 0; i < len(params); i++ {
		if params[i].Text != "" && params[i].Title != "" && params[i].Alias != "" {

			OSPFPost, err := OSPFJson(appkey, params[i].Text, params[i].Title, params[i].Alias, nowTime())
			if err != nil {
				return "", err
			}
			//开始调用个推===============================================0
			//json序列化
			post := OSPFPost
			log.Debug("████████透传消息URL:", url, "████████透传消息post", post)

			var jsonStr = []byte(post)
			log.Debug("LIST_push_jsonStr", jsonStr)
			log.Debug("LIST_push_new_str", bytes.NewBuffer(jsonStr))

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			if err != nil {
				return "", err
			}
			//添加请求头文件
			req.Header.Set("Content-Type", "application/json")
			headMap := req.Header
			headMap["authtoken"] = append(headMap["authtoken"], authtoken)
			//req.Header.Add("authtoken", authtoken)
			log.Debug("authtoken:", req)
			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()
			if err != nil {
				return "", err
			}
			//调用个推结束开始处理返回数据===============================================1
			log.Debug("response Status:", resp.Status)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			str := string(body)
			log.Debug("response Body:", string(body))

			//把返回回来的json，提取authtoken
			authtokenMap, err := jsonMap(str)
			log.Debug("▼▼▼▼▼透传消息,authtokenMap:::::::::", fmt.Sprint(authtokenMap))
			if err != nil {
				return "", err
			}
			result := authtokenMap["result"].(string)
			log.Debug("██透传消息"+fmt.Sprint(i)+"███result:", result)
		}
	}

	return "", nil
}

//全部用户推送通知结构体封装===============================================0
func AllPushParamJson(appkey string, text string, title string, requestid string, jsonInfo string) (string, error) {
	alert := alert{
		Title: title,
		Body:  text,
	}
	iosInfo := iosInfo{
		Json: jsonInfo, //推送通知的时候的隐藏参数
	}
	aps := aps{
		Alert:   alert,
		IosInfo: iosInfo,
		Badge:   1,
	}
	push_info := push_info{
		Aps: aps,
	}
	style := styleStruct{
		Type:  0,
		Text:  text,
		Title: title,
	}
	notification := notificationStruct{
		Style:                style,
		Transmission_type:    false,
		Transmission_content: jsonInfo, //推送通知的时候的隐藏参数
	}
	message := messageStruct{
		Appkey:              appkey,
		Is_offline:          true,
		Offline_expire_time: 10000000,
		Msgtype:             "notification",
	}
	allPush := AllPush{
		Message:      message,
		Notification: notification,
		Push_info:    push_info,
		Requestid:    requestid,
	}
	return structJson(allPush)
}

//全部用户推送通知结构体封装===============================================1

//给所有用户推送通知=====================================================0
func PushAll(paramsStruct *pushInfoList) (string, error) {
	url := Push_All_Url

	authtoken, err := authtoken()
	if err != nil {
		log.Debug("authtoken:", err)
		return "", err
	}

	params := paramsStruct.PushInfos
	//遍历参数
	for i := 0; i < 1; i++ {
		if params[i].Text != "" && params[i].Title != "" && params[i].Alias != "" {
			androidPost, err := AllPushParamJson(appkey, params[i].Text, params[i].Title, nowTime(), params[i].JsonInfo)
			if err != nil {
				return "", err
			}
			//调用推送========================
			//json序列化
			post := androidPost
			log.Debug("████████所有用户推送通知URL:", url, "████████所有用户推送通知post", post)

			var jsonStr = []byte(post)
			log.Debug("LIST_push_jsonStr", jsonStr)
			log.Debug("LIST_push_new_str", bytes.NewBuffer(jsonStr))

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			if err != nil {
				log.Debug("http.NewRequest:", err)
				return "", err
			}
			//添加请求头文件
			req.Header.Set("Content-Type", "application/json")
			headMap := req.Header
			headMap["authtoken"] = append(headMap["authtoken"], authtoken)
			//req.Header.Add("authtoken", authtoken)
			log.Debug("authtoken:", req)
			client := &http.Client{}
			resp, err := client.Do(req)
			if resp != nil {
				defer resp.Body.Close() // ok，即使不读取Body中的数据，即使Body是空的，也要调用close方法,一定要判断resp是否为空
			}
			if err != nil {
				log.Debug("client.Do:", err)
				return "", err
			}
			//调用完成开始处理返回数据================================================================
			log.Debug("response Status:", resp.Status)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Debug("ioutil.ReadAll:", err)
				return "", err
			}
			str := string(body)
			log.Debug("response Body:", string(body))

			//把返回回来的json，提取authtoken
			authtokenMap, err := jsonMap(str)
			log.Debug("▼▼▼▼▼所有用户推送,authtokenMap:::::::::", fmt.Sprint(authtokenMap))
			if err != nil {
				return "", err
			}
			result := authtokenMap["result"].(string)
			log.Debug("██所有用户推送通知"+fmt.Sprint(i)+"███result:", result)
		}
	}

	return "", nil
}

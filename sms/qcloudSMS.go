// qcloudSMS
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	ch "mtcomm/caller/http3part"
	logger "mtcomm/log"
	"net/http"
	"time"
)

var (
	nationcode                  string
	nowTime                     string
	sdkappid                    string
	random                      string
	sign                        string
	appkey                      string
	Template_More_qcloudSMS_URL string
)

func init() {
	nationcode = "86"
	/* init properties */

	sdkappid = prop.Get("sms_sdkappid")
	sign = "钦家"
	appkey = prop.Get("sms_appkey")
	Template_More_qcloudSMS_URL = "https://yun.tim.qq.com/v5/tlssmssvr/sendmultisms2"

	/* init log */
	logger.SetDefaultLogLevel(logger.LevelDebug)
	log = logger.GetDefaultLogger()

	/* init random */
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	random = fmt.Sprintf("%04v", rnd.Int31n(10000))
}

type More_Template_qcloudSMS interface {
	More_Template_qcloudSMS(inf *smsRequest) (*smsResponse, error)
}

type more_Template_qcloudSMS struct {
	caller ch.CallProxyStruction
}

func NewMore_Template_qcloudSMS(url string) More_Template_qcloudSMS {
	return &more_Template_qcloudSMS{
		caller: ch.CallProxyStruction{
			Tracer:      tracer,
			CommandName: "sms_qcloudSMS_More_Template_qcloudSMS",
			CalledUrl:   url,
		},
	}
}

//sig需要的sha256算法
func getSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
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

//模板短信群发
func (caller *more_Template_qcloudSMS) More_Template_qcloudSMS(info *smsRequest) (*smsResponse, error) {

	p := &ch.CallerParameter{
		ErrorPercentThreshold: 25,
		Timeout:               600000,
		MaxThread:             500,
		HttpMethod:            "POST",
		DecodeResponseFunc: func(_ context.Context, r *http.Response) (interface{}, error) {
			b, er := ioutil.ReadAll(r.Body)
			if er != nil {
				return nil, er
			}
			rMap, _ := jsonMap(string(b))

			return &smsResponse{Par: rMap}, nil
		},
	}

	e, err1 := caller.caller.MakeRemoteEndpoint(p)
	if err1 != nil {
		return nil, err1
	}

	resp, err2 := e(context.TODO(), info)
	if err2 != nil {
		return nil, err2
	}

	result, _ := resp.(*smsResponse)
	return result, nil
}

//func main() {
//	a := []string{"aaa"}
//	b := []string{"17671774535"}
//	SmsInfo := SmsInfo{
//		Params: a,
//		Mobile: b,
//		Tpl_id: "57999",
//	}
//	More_Template_qcloudSMS(&SmsInfo)
//}

//func main() {
//	a := []string{"aaa"}
//	b := []string{"17671774535"}
//	SmsInfo := SmsInfo{
//		Params: a,
//		Mobile: b,
//		Tpl_id: "57999",
//	}
//	More_Template_qcloudSMS(&SmsInfo)
//}

type smsRequest struct {
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    []tel    `json:"tel"`
	Time   string   `json:"time"`
	Tpl_id string   `json:"tpl_id"`
}

type tel struct {
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
}

type smsResponse struct {
	Par map[string]interface{} `json:"par"`
}

func smsRequestInit(info *SmsInfo) (*smsRequest, string) {
	tels := []tel{}

	params := info.Params
	mobile := info.Mobile
	tpl_id := info.Tpl_id
	//遍历手机号
	mobiles := ""
	for i := 0; i < len(mobile); i++ {
		mobiles += mobile[i]
		tel := tel{Mobile: mobile[i], Nationcode: nationcode}
		tels = append(tels, tel)
		if i < len(mobile)-1 {
			mobiles += ","
		}
	}
	nowTime = fmt.Sprint(time.Now().Unix())
	var sig_s bytes.Buffer
	sig_s.WriteString("appkey=")
	sig_s.WriteString(appkey)
	sig_s.WriteString("&random=")
	sig_s.WriteString(random)
	sig_s.WriteString("&time=")
	sig_s.WriteString(nowTime)
	sig_s.WriteString("&mobile=")
	sig_s.WriteString(mobiles)
	sig := getSha256(sig_s.String())

	var url bytes.Buffer
	url.WriteString(Template_More_qcloudSMS_URL)
	url.WriteString("?sdkappid=")
	url.WriteString(sdkappid)
	url.WriteString("&random=")
	url.WriteString(random)
	URL := url.String()

	smsReq := &smsRequest{
		Params: params,
		Sig:    sig,
		Sign:   sign,
		Tel:    tels,
		Time:   nowTime,
		Tpl_id: tpl_id,
	}

	return smsReq, URL
}

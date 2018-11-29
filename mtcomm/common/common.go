package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func Json2Struct(jsonStr string, result interface{}) (err error) {
	return json.Unmarshal([]byte(jsonStr), result)
}

func Struct2Json(struction interface{}) (result string, err error) {
	bs, err := json.Marshal(struction)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//map的key大小写都能映射到result的key。result为struct的指针
func Map2Struct(m map[string]interface{}, result interface{}) error {
	return mapstructure.Decode(m, result)
}

//数组变字符串
func Slice2String(s ...interface{}) string {
	if s == nil || len(s) == 0 {
		return ""
	}

	var str bytes.Buffer
	for _, v := range s {
		str.WriteString(fmt.Sprint(v))
		str.WriteString(", ")
	}
	result := str.String()
	if len(result) > 2 {
		result = result[:(len(result) - 2)]
	}
	return result
}

//数组变字符串，逗号分隔
func Slice2StringByKoma(s []string) string {
	if s == nil || len(s) == 0 {
		return ""
	}

	if len(s) == 1 {
		return s[0]
	}

	var str bytes.Buffer
	for _, v := range s {
		str.WriteString(fmt.Sprint(v))
		str.WriteString(",")
	}
	result := str.String()
	if len(result) > 2 {
		result = result[:(len(result) - 1)]
	}
	return result
}

//字符串变字符串数组，逗号分隔
func String2Slice(s string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

//二维数组变字符串，逗号及中括号分隔
func SliceTwo2String(s [][]string) string {
	if s == nil || len(s) == 0 {
		return ""
	}

	var str bytes.Buffer
	for _, v := range s {
		str.WriteString("[")
		str.WriteString(Slice2String(v))
		str.WriteString("]")
	}
	return str.String()
}

//json转map
func json2Map(str string) (map[string]interface{}, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(str), &dat); err == nil {
		return dat, nil
	} else {
		return make(map[string]interface{}), err
	}
}

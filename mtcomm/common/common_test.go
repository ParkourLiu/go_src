package common_test

import (
	"fmt"
	"mtcomm/common"
	"testing"
)

type Person struct {
	UserId   string
	UserName string
	Age      int
}

func (p *Person) String() string {
	return p.UserId + ":" + p.UserName + ":" + fmt.Sprint(p.Age)
}

func TestStruct2Json(t *testing.T) {
	j, err := common.Struct2Json(&Person{UserId: "1", UserName: "lio", Age: 30})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(j)

	s := &Person{}
	err = common.Json2Struct(j, s)
	if err != nil {
		t.Error(err)
	}
	if s.UserId != "1" || s.UserName != "lio" || s.Age != 30 {
		t.Error("error")
	}
}

func TestStruct2Map(t *testing.T) {
	m := common.Struct2Map(Person{UserId: "1", UserName: "lio", Age: 30})

	if len(m) != 3 || m["UserId"] != "1" || m["UserName"] != "lio" || m["Age"] != 30 {
		t.Error("error")
	}
}

func TestMap2Struct(t *testing.T) {
	m := map[string]interface{}{"usErId": "2", "uSerName": "yoyo", "age": 18}
	p := &Person{}
	err := common.Map2Struct(m, p)
	if err != nil {
		t.Error(err)
	}

	if p.UserId != "2" || p.UserName != "yoyo" || p.Age != 18 {
		t.Error("error")
	}
}

func TestSlice2String(t *testing.T) {
	m := []string{"usErId", "2"}
	str := common.Slice2StringByKoma(m)
	fmt.Println("str", str)
	if str != "usErId,2" {
		t.Error("error")
	}
}

func TestSlice2String_2(t *testing.T) {
	m := []string{"usErId"}
	str := common.Slice2StringByKoma(m)
	fmt.Println("str", str)
	if str != "usErId" {
		t.Error("error")
	}
}

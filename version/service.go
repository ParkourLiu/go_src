package main

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionService interface {
	VersionInfo(*Version) (map[string]string, string, error)
}

type versionService struct{}

type Version struct {
	VersionId    string `json:"versionId"`
	NewVersion   string `json:"newVersion"`
	IsForce      string `json:"isForce"`
	AppStoreLink string `json:"appStoreLink"`
	Description  string `json:"description"`
	Flag         string `json:"flag"`
	CreateTime   string `json:"createTime"`
}

//批量删除图片，传入PhotoBook{PbId,lookUserId,Users[{Uid,Pid}]}
func (service versionService) VersionInfo(version *Version) (map[string]string, string, error) {
	log.Debug("method_start", "VersionInfo", "input", fmt.Sprint(version))
	newVersion, err1 := SVersion(mysqlClient, version) //最新版本
	if err1 != nil {
		return nil, "", err1
	}
	version.IsForce = "1"
	mustUpVersion, err := SVersion(mysqlClient, version) //必须要更新的最新版本
	if err != nil {
		return nil, "", err
	}
	newVerInt := string2int(newVersion["newVersion"])
	mustUpVerInt := string2int(mustUpVersion["newVersion"])
	verInt := string2int(version.NewVersion)
	if mustUpVerInt > verInt { //有必须要更新的版本
		return mustUpVersion, flag1, nil

	} else if newVerInt > verInt { //有不是必须更新的版本
		return newVersion, flag2, nil

	}
	log.Debug("method_end", "VersionInfo", "status", "success")
	return nil, "100", nil

}

//版本号转换
func string2int(version string) int {
	a := 0
	aaa := strings.Split(version, ".")
	if len(aaa) == 3 {
		a0, _ := strconv.Atoi(aaa[0])
		a1, _ := strconv.Atoi(aaa[1])
		a2, _ := strconv.Atoi(aaa[2])
		a = a0*1000000 + a1*1000 + a2
	} else if len(aaa) == 2 {
		a0, _ := strconv.Atoi(aaa[0])
		a1, _ := strconv.Atoi(aaa[1])
		a = a0*1000000 + a1*1000
	}
	return a
}

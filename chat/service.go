package main

import (
	rcserversdk "chat/rcserversdk"
	"errors"
	"fmt"
	"mtcomm/db/mysql"
	//"mtcomm/db/redis"
	"strconv"
	"strings"
	"time"
)

type RYToKenService interface {
	GetRYToken(info *ToKenInfo) (data interface{}, code string, err error)
	GetUserInfo(info *ToKenInfo) (data interface{}, code string, err error)
	CreateGroupChat(info *GroupChatInfo) (data string, code string, err error)
	JoinGroupChat(info *GroupChatInfo) (code string, err error)
	QuitGroupChat(info *GroupChatInfo) (code string, err error)
	QueryGroupChatMemberList(info *GroupChatInfo) (data map[string]interface{}, code string, err error)
	GetArrayGroupInfo(info *GroupChatInfo) (data interface{}, code string, err error)
	GetMyGroupChat(info *GroupChatInfo) (data interface{}, code string, err error)
	UpdateGroupChat(info *GroupChatInfo) (data interface{}, code string, err error)
	Dismiss(info *GroupChatInfo) (code string, err error)

	//parkour======================================================start
	AddOfficialMSG(*Chat) (string, error)                  //添加官方消息
	LookOfficialMSG() ([]map[string]string, string, error) //查看官方消息列表
	InformChat(*Chat) (string, error)                      //举报群聊
	//parkour======================================================end

	SearchChatInfo(info *GroupChatInfo) (map[string]string, string, error) //获取群聊信息
	GetClassId(info *GroupChatInfo) (string, string, error)                //根据群聊id获取班级id

}

type rYToKenService struct{}

//parkour======================================================start

func officialMSGsKey() string { //redis官方消息id汇总key
	return "officialMSGs"
}

func officialMSGKey(officialMSGId string) string { //redis官方消息详细内容key
	return "officialMSG:" + officialMSGId
}

func informChatsKey() string { //redis举报群聊id汇总key
	return "informChats"
}

func informChatKey(informChatId string) string { //redis举报群聊详细内容key
	return "informChat:" + informChatId
}

/*官方消息存储redis
    key                            	type          value         score
officialMSGs                        Zset      	officialMSGId     时间戳
officialMSG:${officialMSGId}      	hash      	officialMSGId
			                                  	title
										       	content
												url
												createTime
*/

/*举报群聊存储redis
    key                            	type          value         score
informChats                        	Zset      	informChatId     时间戳
informChat:${informChatId}      	hash      	groupChatId
			                                  	informUserId
										       	informExplain
												informImg
												createTime
*/
//传参 Chat{title，content，url}
func (service rYToKenService) AddOfficialMSG(chat *Chat) (string, error) {
	log.Debug("method_start", "AddOfficialMSG", "input", fmt.Sprint(chat))
	/* service */
	officialMSGId := idGenClient.GetUniqueId()

	AddOfficialMSGMap := map[string]string{}
	AddOfficialMSGMap["officialMSGId"] = officialMSGId
	AddOfficialMSGMap["title"] = chat.Title
	AddOfficialMSGMap["content"] = chat.Content
	AddOfficialMSGMap["url"] = chat.Url
	AddOfficialMSGMap["createTime"] = time.Now().Format("2006-01-02 15:04:05")
	err := redisClient.Hmset(officialMSGKey(officialMSGId), AddOfficialMSGMap)
	if err != nil {
		return "", err
	}
	_, err = redisClient.Zadd(officialMSGsKey(), time.Now().Unix(), officialMSGId)
	if err != nil {
		return "", err
	}
	t := rcserversdk.TxtMessage{Content: "1", Extra: "1"}
	codeSuccess, err := message.Broadcast(UserId, t, "", "", "")
	if err != nil {
		return "", err
	}
	if codeSuccess.Code != 200 {
		return "", errors.New(codeSuccess.ErrorMessage)
	}
	log.Debug("method_end", "AddOfficialMSG", "status", "success")
	return "100", nil
}

func (service rYToKenService) LookOfficialMSG() ([]map[string]string, string, error) {
	log.Debug("method_start", "LookOfficialMSG", "input", "")
	/* service */
	returnList := []map[string]string{}
	officialMSGs, err := redisClient.Zrange2(officialMSGsKey(), int64(0), int64(-1))
	if err != nil {
		return nil, "", err
	}
	for _, v := range officialMSGs {
		OfficialMSGMap, err1 := redisClient.HgetAllMap(officialMSGKey(v))
		if err1 != nil {
			return nil, "", err1
		}
		returnList = append(returnList, OfficialMSGMap)
	}
	log.Debug("method_end", "LookOfficialMSG", "status", "success")
	return returnList, "100", nil
}

//举报群聊
func (service rYToKenService) InformChat(chat *Chat) (string, error) {
	log.Debug("method_start", "InformChat", "input", fmt.Sprint(chat))
	/* service */
	informChatId := idGenClient.GetUniqueId()

	InformChatMap := map[string]string{}
	InformChatMap["informChatId"] = informChatId
	InformChatMap["groupChatId"] = chat.GroupChatId
	InformChatMap["informUserId"] = chat.InformUserId
	InformChatMap["informExplain"] = chat.InformExplain
	InformChatMap["informImg"] = chat.InformImg
	InformChatMap["createTime"] = time.Now().Format("2006-01-02 15:04:05")
	err := redisClient.Hmset(informChatKey(informChatId), InformChatMap)
	if err != nil {
		return "", err
	}
	_, err = redisClient.Zadd(informChatsKey(), time.Now().Unix(), informChatId)
	if err != nil {
		return "", err
	}
	log.Debug("method_end", "InformChat", "status", "success")
	return "100", nil
}

//parkour======================================================end

func c(strr string) bool {
	str := strings.Replace(strr, " ", "", -1)
	if str == "" {
		return false
	}
	return true
}

//得到融云token
func (self *rYToKenService) GetRYToken(info *ToKenInfo) (data interface{}, code string, err error) {
	if !c(info.Status) || !c(info.UserId) {
		return nil, rytoken101, nil
	}
	//	if info.Status != "0" && info.Status != "1" {
	//		return nil, rytoken102, nil
	//	}
	//登陆时
	//status为0时候，根据userid查询redis数据库里是否存在token，存在就直接返回，不存在去融云上获得一个新的token，并存入redis数据库
	//if info.Status == "0" {
	re, err := dao.SelectRyToken(redisClient, RyUserId+info.UserId)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return nil, "", err
			//return nil, rytoken103, nil
		}
		//return nil, "", err
	}
	if re != "" {
		return re, "100", nil
	} else {

		//}
		//注册时
		//从redis里读取user的信息  HgetAllMap(key string) (value map[string]string, err error)
		mapR, err := redisClient.HgetAllMap(userId + info.UserId)

		if err != nil {
			if err.Error() == "redigo: nil returned" {
				return nil, rytoken104, err
			} else {
				return nil, "", err
			}

		}

		if len(mapR) == 0 {
			return nil, rytoken104, nil
		} else {
			ii := 0
			nickName := 0
			imageName := 0
			for k, _ := range mapR {
				if k == "userId" {
					ii = ii + 1
				}
				if k == "nickName" {
					nickName = nickName + 1
				}
				if k == "imageName" {
					imageName = imageName + 1
				}
			}
			if ii != 1 {
				return nil, rytoken105, nil
			}
			if nickName != 1 {
				mapR["nickName"] = defaultNickName
			}
			if imageName != 1 {
				mapR["imageName"] = defaultImageName
			}
			if !c(mapR["nickName"]) {
				mapR["nickName"] = defaultNickName
			}
			if !c(mapR["imageName"]) {
				mapR["imageName"] = defaultImageName
			}
		}
		resulu, err := user.GetToken(mapR["userId"], mapR["nickName"], mapR["imageName"])
		if err != nil {
			return nil, "", err
		}
		//RyUserId    = "chat:rytoken:keyIsUserId:"
		//UserIdRy    = "chat:userId:keyIsryToken:"
		re := make([][]interface{}, 1)
		re[0] = []interface{}{RyUserId + resulu.UserId, resulu.Token}
		//re[1] = []interface{}{UserIdRy + resulu.Token, resulu.UserId}
		_, err = dao.Add(redisClient, re)
		if err != nil {
			return nil, "", err
		}
		return resulu.Token, "100", nil
	}
	//}
	//return nil, "rytoken102", nil
}

//批量得到用户的信息
func (self *rYToKenService) GetUserInfo(info *ToKenInfo) (data interface{}, code string, err error) {

	if len(info.UserIdArray) == 0 || !c(info.UserId) {
		return nil, rytoken101, nil
	}
	//得到所有的user信息
	lis3, err3 := dao.ListPipelineGetHashMap_User(redisClient, info.UserIdArray)
	if len(lis3) == 0 {
		return nil, rytoken107, err3
	}
	if err3 != nil {
		return nil, "", err3
	}
	lis4 := GetMapArray(lis3, len(info.UserIdArray))
	lis5 := []map[string]string{}
	for _, v := range lis4 {
		if len(v) != 0 {
			lis5 = append(lis5, v)
		}
	}
	lis6 := make([]interface{}, len(lis5))
	for k, v := range lis5 {
		cc := make(map[string]string)
		for k1, v1 := range v {
			switch k1 {
			case "ad":
				cc["ad"] = v1
			case "imageName":
				cc["imageName"] = v1
			case "nickName":
				cc["nickName"] = v1
			case "platForm":
				cc["platForm"] = v1
			case "userId":
				cc["userId"] = v1
			}
		}
		lis6[k] = cc
	}
	return lis6, "100", nil

}

func GetMapArray(params []interface{}, strlen int) []map[string]string {
	paramSlice := make([]map[string]string, strlen)
	for i := 0; i < len(params); i++ {
		param2Slice := []string{}
		switch v := params[i].(type) {
		case []interface{}:
			for _, pa := range v {
				switch v1 := pa.(type) {
				case []uint8:
					strV11 := string(v1)
					param2Slice = append(param2Slice, strV11)
				default:
					panic("params type not supported")
				}
			}
		default:
			panic("params type not supported")
		}
		result := map[string]string{}
		for y := 0; y < len(param2Slice); y++ {
			if y%2 == 0 {
				key := param2Slice[y]
				values := param2Slice[y+1]
				result[key] = values
			}
		}
		paramSlice[i] = result
	}
	return paramSlice
}

//创建群聊
func (self *rYToKenService) CreateGroupChat(info *GroupChatInfo) (data string, code string, err error) {

	log.Debug("method_start", "CreateGroupChat", "input", info)
	log.Debug("chat----GroupChatInfo😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀", fmt.Sprint(info))
	//check
	if !c(info.CreateUserId) || !c(info.GroupChatName) || !c(info.GroupChatNotice) || !c(info.AcId) {
		return "", rytoken101, nil
	}

	//获取id
	id := idGenClient.GetUniqueId()
	if info.AcId[:2] == "cl" {
		id = "sc" + idGenClient.GetUniqueId()
	}
	//在融云平台创建群聊
	userId := []string{info.CreateUserId}
	codeSuccess, err := group.Create(userId, id, info.GroupChatName)
	if err != nil {
		return "", "", err
	}
	if codeSuccess.Code != 200 {
		return "", "", errors.New(codeSuccess.ErrorMessage)
	}
	//持久化到redis与mysql
	//  //持久化到redis
	mapS2 := map[string]string{}
	mapS2["groupChatId"] = id
	mapS2["createUserId"] = info.CreateUserId
	mapS2["groupChatName"] = info.GroupChatName
	mapS2["groupChatnotice"] = info.GroupChatNotice
	mapS2["groupChatUrl"] = info.GroupChatUrl
	mapS2["acId"] = info.AcId
	mapS2["gId"] = info.GId
	mapS2["createTime"] = time.Now().Format("2006-01-02 15:04:05")
	err = redisClient.Hmset(groupChatInfo+id, mapS2)
	if err != nil {
		return "", "", err
	}
	//  //持久化到mysql
	err = mysqlClient.Execute(&mysql.Stmt{Sql: "INSERT INTO groupchat(groupChatId,createUserId,groupChatName,groupChatnotice,groupChatUrl,updatetime,acId,gId)VALUES (?,?,?,?,?,NOW(),?,? ) ", Args: []interface{}{id, info.CreateUserId, info.GroupChatName, info.GroupChatNotice, info.GroupChatUrl, info.AcId, info.GId}})
	if err != nil {
		return "", "", err
	}

	//添加成员列表
	re := make([][]interface{}, 2)
	int64time, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 64)
	re[0] = []interface{}{groupChatMemberList + id, int64time, info.CreateUserId}
	re[1] = []interface{}{memberJoinGroupChatList + info.CreateUserId, int64time, id}
	_, err = dao.DaoJoin(redisClient, re)
	if err != nil {
		return "", "", err
	}
	log.Debug("method_end", "CreateGroupChat")
	return id, "100", nil
}

//加入群聊 func (self * Group)Join(userId []string, groupId string, groupName string)(*CodeSuccessReslut, error)
func (self *rYToKenService) JoinGroupChat(info *GroupChatInfo) (code string, err error) {
	log.Debug("method_start", "JoinGroupChat", "input", info)
	//check
	if len(info.UserIdArray) == 0 || !c(info.AcId) {
		return rytoken101, nil
	}

	//根据活动id判断最新的群聊id是否超过指定人数
	groupCha, err := dao.GetGroupChatId(mysqlClient, info.AcId)
	if err != nil {
		return "", err
	}
	if len(groupCha) == 0 {
		return chat113, nil
	}
	//获取群聊成员的数量  Zlen(key string) (value int64, err error

	in64, err := redisClient.Zlen(groupChatMemberList + groupCha[0]["groupChatId"])
	if err != nil {
		return "", err
	}
	if in64 == 0 {
		return chat108, nil
	}
	if in64 >= int64(groupChatNum) {
		//创建群聊
		service := &rYToKenService{}
		activityInfo, err := dao.GetActivityInfo(redisClient, info.AcId)
		if err != nil {
			return "", err
		}

		inf := &GroupChatInfo{UserId: activityInfo["uId"], CreateUserId: activityInfo["uId"], GroupChatName: activityInfo["name"] + strconv.Itoa(len(groupCha)+1), GroupChatUrl: activityInfo["cover"], GroupChatNotice: activityInfo["introduce"], GId: activityInfo["gId"], AcId: info.AcId}
		data, code, err := service.CreateGroupChat(inf)
		if err != nil {
			return "", err
		}
		if code != "100" {
			return code, nil
		}
		info.GroupChatId = data
		info.GroupChatName = activityInfo["name"] + strconv.Itoa(len(groupCha)+1)
	} else {
		info.GroupChatId = groupCha[0]["groupChatId"]
		info.GroupChatName = groupCha[0]["groupChatName"]
	}

	codeSuccess, err := group.Join(info.UserIdArray, info.GroupChatId, info.GroupChatName)
	if err != nil {
		return "", err
	}
	if codeSuccess.Code != 200 {
		return "", errors.New(codeSuccess.ErrorMessage)
	}
	//groupChatMember
	re := make([][]interface{}, len(info.UserIdArray)*2)
	int64time, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 64)
	for k, v := range info.UserIdArray {
		re[k] = []interface{}{groupChatMemberList + info.GroupChatId, int64time, v}
	}
	for k, v := range info.UserIdArray {
		re[len(info.UserIdArray)+k] = []interface{}{memberJoinGroupChatList + v, int64time, info.GroupChatId}
	}

	//加入群聊
	_, err = dao.DaoJoin(redisClient, re)
	if err != nil {
		return "", err
	}

	log.Debug("method_end", "JoinGroupChat")
	return "100", nil
}

//移除群聊 func (self * Group)Quit(userId []string, groupId string)(*CodeSuccessReslut, error)
func (self *rYToKenService) QuitGroupChat(info *GroupChatInfo) (code string, err error) {
	log.Debug("method_start", "QuitGroupChat", "input", info)
	//check
	if len(info.UserIdArray) == 0 || !c(info.GroupChatId) {
		return rytoken101, nil
	}
	//移除群聊
	resp, err := group.Quit(info.UserIdArray, info.GroupChatId)
	if err != nil {
		return "", err
	}
	if resp.Code != 200 {
		return "", errors.New(resp.ErrorMessage)
	}
	//移除群聊列表中的数据
	err = redisClient.Zrem(groupChatMemberList+info.GroupChatId, info.UserIdArray)
	if err != nil {
		return "", err
	}

	log.Debug("method_end", "QuitGroupChat")
	return "100", nil
}

//得到群聊成员基本信息，以及和群公告，群简介等信息。
//  func (self * Group)QueryUser(groupId string)(*GroupUserQueryReslut, error)
func (self *rYToKenService) QueryGroupChatMemberList(info *GroupChatInfo) (data map[string]interface{}, code string, err error) {
	log.Debug("method_start", "QueryGroupChatMemberList", "input", info)
	if !c(info.GroupChatId) || !c(info.UserId) {
		return nil, rytoken101, nil
	}

	re := map[string]interface{}{}
	var lis4 []map[string]string
	//得到群公告等信息
	//HgetAllMap(key string) (value map[string]string, err error)
	chatMap, err := redisClient.HgetAllMap(groupChatInfo + info.GroupChatId)
	if err != nil {
		return nil, "", err
	}
	re["groupName"] = chatMap["groupChatName"]
	re["groupNotice"] = chatMap["groupChatnotice"]
	re["groupChatUrl"] = chatMap["groupChatUrl"]
	re["createTime"] = chatMap["createTime"]
	re["createUserId"] = chatMap["createUserId"]
	if chatMap["createUserId"] == info.UserId {
		re["isCreate"] = "1"
	} else {
		re["isCreate"] = "0"
	}

	//列出群聊的基本信息
	if c(info.GroupChatId) {
		var start, stop int64
		//根据lastUId 判断下一条的序列号
		start1, err := dao.ListZRANK(redisClient, groupChatMemberList+info.GroupChatId, info.LastUserId)
		if err != nil {
			if err.Error() == "redigo: nil returned" {
				//return nil, rytoken107   , err
			} else {
				return nil, "", err
			}
		}
		//这里不做分页  设置1000  表示要显示全部数据
		stop1, err := strconv.ParseInt("1000", 10, 64)
		if err != nil {
			return nil, "", err
		}

		lis := []string{}
		if !c(info.LastUserId) {
			start = 0
			stop = stop1 - 1
			re["page"] = "0"
			//得到粉丝列表
			lis, err = dao.ListZrange(redisClient, groupChatMemberList+info.GroupChatId, start, stop)
			if err != nil {
				return nil, "", err
			}

		} else {
			if start1 == 0 {
				//根据条件查不到数据
				return nil, rytoken107, nil
			}
			start = start1 - stop1
			stop = start1 - 1
			if start < int64(0) {
				start = 0
				if stop < int64(0) {
					stop = 0
				}
				re["page"] = "1"
			}
			//得到粉丝列表
			lis, err = dao.ListZrange2(redisClient, groupChatMemberList+info.GroupChatId, start, stop)
			if err != nil {
				return nil, "", err
			}

		}
		//score值最小的成員值
		minMember, err := dao.ListZrange2(redisClient, groupChatMemberList+info.GroupChatId, int64(0), int64(0))
		if err != nil {
			return nil, "", err
		}
		if len(minMember) == 0 {
			return nil, "", errors.New("minMember长度为0")
		}
		if minMember[0] == lis[len(lis)-1] {
			re["page"] = "1"
		} else {
			re["page"] = "0"
		}
		//获取群聊成员的数量  Zlen(key string) (value int64, err error
		in64, err := redisClient.Zlen(groupChatMemberList + info.GroupChatId)
		if err != nil {
			return nil, "", err
		}
		str64 := strconv.FormatInt(in64, 10)
		re["userNum"] = str64
		//得到所有的user信息
		lis3, err3 := dao.ListPipelineGetHashMap_User(redisClient, lis)
		if err3 != nil {
			return nil, "", err3
		}
		if len(lis3) == 0 {
			return nil, rytoken107, err3
		}
		//拼成指定的格式
		lis4 = GetMapArray(lis3, len(lis))
		for k, _ := range lis4 {
			if lis4[k]["userId"] == chatMap["createUserId"] {
				lis4[k]["isCreate"] = "1"
			} else {
				lis4[k]["isCreate"] = "0"
			}

		}

		//反序输出
		resu := make([]map[string]string, len(lis4))
		if !c(info.LastUserId) {
			re["users"] = lis4
		} else {
			j := 0
			for i := len(lis4) - 1; i >= 0; i-- {
				resu[j] = lis4[i]
				j++
			}
			re["users"] = resu
		}

	}
	log.Debug("method_end", "QueryGroupChatMemberList")
	return re, "100", nil
}

//批量得到群聊的相关信息
func (self *rYToKenService) GetArrayGroupInfo(info *GroupChatInfo) (data interface{}, code string, err error) {
	if !c(info.UserId) || len(info.GroupChatIdArray) == 0 {
		return nil, rytoken101, nil
	}
	//得到所有的qun聊的信息
	lis3, err3 := dao.ListPipelineGetHashMap_GroupChat(redisClient, info.GroupChatIdArray)
	if len(lis3) == 0 {
		return nil, rytoken107, err3
	}
	if err3 != nil {
		return nil, "", err3
	}
	lis4 := GetMapArray(lis3, len(info.GroupChatIdArray))
	lis5 := []map[string]string{}
	for _, v := range lis4 {
		if len(v) != 0 {
			lis5 = append(lis5, v)
		}
	}
	lis6 := make([]interface{}, len(lis5))
	for k, v := range lis5 {
		cc := make(map[string]string)
		for k1, v1 := range v {
			switch k1 {
			case "groupChatId":
				cc["groupChatId"] = v1
			case "createUserId":
				cc["createUserId"] = v1
			case "groupChatName":
				cc["groupChatName"] = v1
			case "groupChatnotice":
				cc["groupChatnotice"] = v1
			case "groupChatUrl":
				cc["groupChatUrl"] = v1
			}
		}
		lis6[k] = cc
	}
	return lis6, "100", nil

}

//我的群
func (self *rYToKenService) GetMyGroupChat(info *GroupChatInfo) (data interface{}, code string, err error) {
	if !c(info.UserId) {
		return nil, rytoken101, nil
	}
	//获取我加入的群聊列表

	lis, err := dao.ListZrange2(redisClient, memberJoinGroupChatList+info.UserId, int64(0), int64(-1))
	if err != nil {
		return nil, "", err
	}
	//得到所有的qun聊的信息
	lis3, err3 := dao.ListPipelineGetHashMap_GroupChat(redisClient, lis)
	if len(lis3) == 0 {
		return nil, rytoken107, err3
	}
	if err3 != nil {
		return nil, "", err3
	}
	lis4 := GetMapArray(lis3, len(lis))
	lis5 := []map[string]string{}
	for _, v := range lis4 {
		if len(v) != 0 {
			lis5 = append(lis5, v)
		}
	}
	lis6 := make([]interface{}, len(lis5))
	for k, v := range lis5 {
		cc := make(map[string]string)
		for k1, v1 := range v {
			switch k1 {
			case "groupChatId":
				cc["groupChatId"] = v1
			case "createUserId":
				cc["createUserId"] = v1
			case "groupChatName":
				cc["groupChatName"] = v1
			case "groupChatnotice":
				cc["groupChatnotice"] = v1
			case "groupChatUrl":
				cc["groupChatUrl"] = v1
			}
		}
		lis6[k] = cc
	}
	return lis6, "100", nil

}

//修改群聊
func (self *rYToKenService) UpdateGroupChat(info *GroupChatInfo) (data interface{}, code string, err error) {
	log.Debug("method_start", "UpdateGroupChat", "input", "=="+fmt.Sprint(info))
	if !c(info.UserId) || !c(info.GroupChatId) {
		return nil, rytoken101, nil
	}
	if !c(info.GroupChatUrl) && !c(info.GroupChatName) && !c(info.GroupChatNotice) {
		return nil, chat110, nil
	}

	//判断是否是群主(value map[string]string, err error)
	mapG, err := redisClient.HgetAllMap(groupChatInfo + info.GroupChatId)
	if err != nil {
		return nil, chat111, err
	}
	if mapG["createUserId"] != info.UserId {
		return nil, chat112, nil
	}

	//跟新mysql里的数据 UPDATE groupchat SET groupChatName='d'   WHERE groupChatId='1'
	sql := " UPDATE groupchat SET "
	if c(info.GroupChatUrl) {
		sql = sql + "groupChatUrl='" + info.GroupChatUrl + "',"
	}
	if c(info.GroupChatName) {
		sql = sql + "groupChatName='" + info.GroupChatName + "',"
	}
	if c(info.GroupChatNotice) {
		sql = sql + "groupChatNotice='" + info.GroupChatNotice + "',"
	}
	sql = sql + " updateTime=NOW() WHERE groupChatId='" + info.GroupChatId + "'"
	err = mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
	if err != nil {
		return nil, "", err
	}
	//删除redis里的数据  Del  groupChatInfo
	err = redisClient.Del(groupChatInfo + info.GroupChatId)
	if err != nil {
		return nil, "", err
	}

	//再次跟新redis里的数据
	mapGroup, err := mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "SELECT * FROM groupchat WHERE groupchatId= ? ", Args: []interface{}{info.GroupChatId}})
	mapS2 := map[string]string{}
	mapS2["groupChatId"] = mapGroup["groupChatId"]
	mapS2["createUserId"] = mapGroup["createUserId"]
	mapS2["groupChatName"] = mapGroup["groupChatName"]
	mapS2["groupChatnotice"] = mapGroup["groupChatNotice"]
	mapS2["groupChatUrl"] = mapGroup["groupChatUrl"]
	mapS2["acId"] = mapGroup["acId"]
	mapS2["gId"] = mapGroup["gId"]
	mapS2["createTime"] = mapGroup["createTime"]
	err = redisClient.Hmset(groupChatInfo+info.GroupChatId, mapS2)
	if err != nil {
		return "", "", err
	}
	return info.GroupChatId, "100", nil
}

//解散群聊
func (self *rYToKenService) Dismiss(info *GroupChatInfo) (code string, err error) {
	if !c(info.UserId) || !c(info.GroupChatId) {
		return rytoken101, nil
	}
	//判断是否是群主(value map[string]string, err error)
	mapG, err := redisClient.HgetAllMap(groupChatInfo + info.GroupChatId)
	if err != nil {
		return chat111, err
	}
	if mapG["createUserId"] != info.UserId {
		return chat112, nil
	}
	//解散群聊
	codeSuccess, err := group.Dismiss(info.UserId, info.GroupChatId)
	if err != nil {
		return "", err
	}
	if codeSuccess.Code != 200 {
		return "", errors.New(codeSuccess.ErrorMessage)
	}
	//删除redis里群成员关系
	//	groupChatInfo           = "chat:groupChatInfo:"
	//	groupChatMemberList     = "chat:groupChatMemberList:"
	//	memberJoinGroupChatList = "chat:memberJoinGroupChatList:"
	lis, err := dao.ListZrange(redisClient, groupChatMemberList+info.GroupChatId, int64(0), int64(-1))
	if err != nil {
		return "", err
	}
	sql := make([]string, len(lis)+2)
	for _, v := range lis {
		err = dao.Zrem(redisClient, memberJoinGroupChatList+v, info.GroupChatId)
		if err != nil {
			return "", err
		}
	}
	sql = append(sql, groupChatInfo+info.GroupChatId)
	sql = append(sql, groupChatMemberList+info.GroupChatId)

	_, err = dao.ListPipelineDel(redisClient, sql)
	if err != nil {
		return "", err
	}
	sql2 := " UPDATE groupchat SET ad='d' where groupChatId='" + info.GroupChatId + "'"
	err = mysqlClient.Execute(&mysql.Stmt{Sql: sql2, Args: []interface{}{}})
	if err != nil {
		return "", err
	}

	return "100", nil
}

//<<<<<<< HEAD
//func (self *rYToKenService) SearchChatInfo(info *GroupChatInfo) (map[string]string, string, error) {
//	chat := map[string]string{}
//	c := &GroupChatInfo{AcId: "cl" + info.ClId}
//	m, err := GetChatInfo(c)
//	if err != nil {
//		return m, "", err
//	}
//	if len(m) == 0 {
//		return m, chat114, err
//	}
//	chat["groupChatId"] = m["groupChatId"]
//	chat["groupChatName"] = m["groupChatName"]
//	return chat, "100", err
//}
//=======
func (self *rYToKenService) SearchChatInfo(info *GroupChatInfo) (map[string]string, string, error) {
	chat := map[string]string{}
	c := &GroupChatInfo{AcId: info.ClId}
	m, err := GetChatInfo(c)
	if err != nil {
		return m, "", err
	}
	if len(m) == 0 {
		return m, chat114, err
	}
	chat["groupChatId"] = m["groupChatId"]
	chat["groupChatName"] = m["groupChatName"]
	return chat, "100", err
	//>>>>>>> 6fd973622db641c86b6aa300f3205d4369592628
}
func (self *rYToKenService) GetClassId(info *GroupChatInfo) (string, string, error) {
	m, err := GetClassId(info)
	if err != nil {
		return "", "", err
	}
	if len(m) == 0 {
		return "", chat114, err
	}
	//<<<<<<< HEAD
	//	clId := strings.Split(m["acId"], "Ac")
	//	cl := clId[1]
	//	return cl, "100", err
	//=======
	clid := m["acId"]
	return clid, "100", err
	//>>>>>>> 6fd973622db641c86b6aa300f3205d4369592628
}

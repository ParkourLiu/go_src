package main

import (
	"bytes"
	"errors"
	"mtcomm/db/mysql"
)

func userKey(userId string) string {
	userKey := "U:" + userId
	return userKey
}

func i_user(user *User) error {
	sqlBuffer := bytes.Buffer{}
	valueBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("INSERT INTO `user`(")
	valueBuffer.WriteString("VALUES(")
	if user.UserId != "" {
		sqlBuffer.WriteString("`userId`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.UserId)
		valueBuffer.WriteString("',")
	}
	if user.QrImageName != "" {
		sqlBuffer.WriteString("`qrImageName`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.QrImageName)
		valueBuffer.WriteString("',")
	}
	if user.PhoneNo != "" {
		sqlBuffer.WriteString("`phoneNo`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.PhoneNo)
		valueBuffer.WriteString("',")
	}
	if user.Password != "" {
		sqlBuffer.WriteString("`password`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.Password)
		valueBuffer.WriteString("',")
	}
	if user.MtalkNo != "" {
		sqlBuffer.WriteString("`mtalkNo`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.MtalkNo)
		valueBuffer.WriteString("',")
	}
	if user.Wechat_uid != "" {
		sqlBuffer.WriteString("`Wechat_uid`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.Wechat_uid)
		valueBuffer.WriteString("',")
	}
	if user.Wechat_name != "" {
		sqlBuffer.WriteString("`Wechat_name`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.Wechat_name)
		valueBuffer.WriteString("',")
	}
	if user.Wechat_iconurl != "" {
		sqlBuffer.WriteString("`Wechat_iconurl`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.Wechat_iconurl)
		valueBuffer.WriteString("',")
	}
	if user.Wechat_gender != "" {
		sqlBuffer.WriteString("`Wechat_gender`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.Wechat_gender)
		valueBuffer.WriteString("',")
	}
	if user.QQ_uid != "" {
		sqlBuffer.WriteString("`QQ_uid`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.QQ_uid)
		valueBuffer.WriteString("',")
	}
	if user.QQ_name != "" {
		sqlBuffer.WriteString("`QQ_name`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.QQ_name)
		valueBuffer.WriteString("',")
	}
	if user.QQ_iconurl != "" {
		sqlBuffer.WriteString("`QQ_iconurl`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.QQ_iconurl)
		valueBuffer.WriteString("',")
	}
	if user.QQ_gender != "" {
		sqlBuffer.WriteString("`QQ_gender`,")
		valueBuffer.WriteString("'")
		valueBuffer.WriteString(user.QQ_gender)
		valueBuffer.WriteString("',")
	}
	sqlBuffer.WriteString("`createTime`,`updateTime`)")
	valueBuffer.WriteString("NOW(),NOW())")

	sqlBuffer.WriteString(valueBuffer.String())
	sql := sqlBuffer.String()

	log.Debug("██i_user█SQL:", sql)
	return mysqlClient.Execute(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

func s_userById(user *User) (map[string]string, error) {
	if user.UserId == "" {
		return nil, errors.New("userId不可以为空:")
	}
	result, err := redisClient.HgetAllMap(userKey(user.UserId))

	if err != nil {
		return result, errors.New("s_userById查询redis数据error:" + err.Error())
	} else {
		if result["userId"] != "" { //如果redis有数据就直接取redis
			return result, nil
		}
	}

	result, err = mysqlClient.SearchOneRow(&mysql.Stmt{Sql: "select * from `user` where `ad`='a' and `userId`=?", Args: []interface{}{user.UserId}})
	if err != nil {
		return result, errors.New("s_userById查询数据库数据失败:" + err.Error())
	}
	if len(result) > 0 {
		//█████████████████mysql数据上传到redis████████████████████████████")
		err = redisClient.Hmset(userKey(user.UserId), result)
		if err != nil {
			return result, errors.New("s_userById插入redis数据失败:" + err.Error())
		}
	}

	return result, err
}

func s_user(user *User) ([]map[string]string, error) {
	sqlBuffer := bytes.Buffer{}
	sql := "select * from `user` where ad='a'"
	sqlBuffer.WriteString(sql)
	if user.UserId != "" {
		sqlBuffer.WriteString(" and userId='")
		sqlBuffer.WriteString(user.UserId)
		sqlBuffer.WriteString("'")
	}
	if user.PhoneNo != "" {
		sqlBuffer.WriteString(" and phoneNo='")
		sqlBuffer.WriteString(user.PhoneNo)
		sqlBuffer.WriteString("'")
	}
	if user.Sex != "" {
		sqlBuffer.WriteString(" and sex='")
		sqlBuffer.WriteString(user.Sex)
		sqlBuffer.WriteString("'")
	}
	if user.Wechat_uid != "" {
		sqlBuffer.WriteString(" and Wechat_uid='")
		sqlBuffer.WriteString(user.Wechat_uid)
		sqlBuffer.WriteString("'")
	}
	if user.QQ_uid != "" {
		sqlBuffer.WriteString(" and QQ_uid='")
		sqlBuffer.WriteString(user.QQ_uid)
		sqlBuffer.WriteString("'")
	}
	if user.IsWater != "" {
		sqlBuffer.WriteString(" and isWater='")
		sqlBuffer.WriteString(user.IsWater)
		sqlBuffer.WriteString("'")
	}
	if user.Volunteer != "" {
		sqlBuffer.WriteString(" and volunteer='")
		sqlBuffer.WriteString(user.Volunteer)
		sqlBuffer.WriteString("'")
	}
	if sqlBuffer.String() == sql {
		return nil, errors.New("必传参数至少有一个")
	}
	sqlBuffer.WriteString(" order by createTime desc")
	sql = sqlBuffer.String()
	return mysqlClient.SearchMutiRows(&mysql.Stmt{Sql: sql, Args: []interface{}{}})
}

func u_user(user *User) error {
	if user.UserId == "" {
		return errors.New("userId不可以为空:")
	}
	updatemap := make(map[string]interface{})
	if user.PhoneNo != "" {
		updatemap["phoneNo"] = user.PhoneNo
	}
	if user.Password != "" {
		updatemap["password"] = user.Password
	}
	if user.Email != "" {
		updatemap["email"] = user.Email
	}
	if user.TrueName != "" {
		updatemap["trueName"] = user.TrueName
	}
	if user.NickName != "" {
		updatemap["nickName"] = user.NickName
	}
	if user.BirthDay != "" {
		updatemap["birthDay"] = user.BirthDay
	}
	if user.ChineseZodiac != "" {
		updatemap["chineseZodiac"] = user.ChineseZodiac
	}
	if user.Sex != "" {
		updatemap["sex"] = user.Sex
	}
	if user.HomeAddress != "" {
		updatemap["homeAddress"] = user.HomeAddress
	}
	if user.ImageName != "" {
		updatemap["imageName"] = user.ImageName
	}
	if user.Hometown != "" {
		updatemap["hometown"] = user.Hometown
	}
	if user.Description != "" {
		updatemap["description"] = user.Description
	}
	if user.OpenId != "" {
		updatemap["openId"] = user.OpenId
	}
	if user.BackgroundImg != "" {
		updatemap["backgroundImg"] = user.BackgroundImg
	}
	if user.Wechat_uid != "" {
		updatemap["Wechat_uid"] = user.Wechat_uid
	}
	if user.Wechat_name != "" {
		updatemap["Wechat_name"] = user.Wechat_name
	}
	if user.Wechat_iconurl != "" {
		updatemap["Wechat_iconurl"] = user.Wechat_iconurl
	}
	if user.Wechat_gender != "" {
		updatemap["Wechat_gender"] = user.Wechat_gender
	}
	if user.QQ_uid != "" {
		updatemap["QQ_uid"] = user.QQ_uid
	}
	if user.QQ_name != "" {
		updatemap["QQ_name"] = user.QQ_name
	}
	if user.QQ_iconurl != "" {
		updatemap["QQ_iconurl"] = user.QQ_iconurl
	}
	if user.QQ_gender != "" {
		updatemap["QQ_gender"] = user.QQ_gender
	}
	if user.IsWater != "" {
		updatemap["isWater"] = user.IsWater
	}
	if user.Volunteer != "" {
		updatemap["volunteer"] = user.Volunteer
	}
	err := redisClient.Del(userKey(user.UserId))
	if err != nil {
		return errors.New("u_user删除redis数据失败:" + err.Error())
	}
	return mysqlClient.Update("user", "where userId='"+user.UserId+"'", updatemap, true)
}

package main

type User struct {
	UserId         string `json:"userId"`         //主键
	PhoneNo        string `json:"phoneNo"`        //电话
	Password       string `json:"password"`       //密码
	Email          string `json:"email"`          //邮箱
	TrueName       string `json:"trueName"`       //真名
	NickName       string `json:"nickName"`       //昵称
	BirthDay       string `json:"birthDay"`       //生日
	ChineseZodiac  string `json:"chineseZodiac"`  //生肖
	QrImageName    string `json:"qrImageName"`    //用户二维码
	Sex            string `json:"sex"`            //性别
	HomeAddress    string `json:"homeAddress"`    //家庭住址
	ImageName      string `json:"imageName"`      //头像
	ChatName       string `json:"chatName"`       //
	ChatPwd        string `json:"chatPwd"`        //
	MtalkNo        string `json:"mtalkNo"`        //钦家号
	Hometown       string `json:"hometown"`       //
	Description    string `json:"description"`    //
	PlatForm       string `json:"platForm"`       //来源（安卓，ios）
	UUID           string `json:"UUID"`           //设备识别码
	OpenId         string `json:"openId"`         //
	BackgroundImg  string `json:"backgroundImg"`  //
	Wechat_uid     string `json:"Wechat_uid"`     //微信uid
	Wechat_name    string `json:"Wechat_name"`    //微信名
	Wechat_iconurl string `json:"Wechat_iconurl"` //微信头像
	Wechat_gender  string `json:"Wechat_gender"`  //微信性别
	QQ_uid         string `json:"QQ_uid"`         //QQid
	QQ_name        string `json:"QQ_name"`        //qq名
	QQ_iconurl     string `json:"QQ_iconurl"`     //qq头像
	QQ_gender      string `json:"QQ_gender"`      //qq性别
	IsWater        string `json:"isWater"`        //是否是水军
	Volunteer      string `json:"volunteer"`      //是否是志愿者
}
type UserService interface {
	SearchUserById(user *User) (map[string]interface{}, string, error) //根据id查用户信息
	SearchUsers(user *User) (map[string]interface{}, string, error)    //根据其他参数查用户信息
	AddUser(user *User) (map[string]interface{}, string, error)        //添加用户
	UpdateUser(user *User) (map[string]interface{}, string, error)     //修改用户信息
}

type userService struct{}

func (service userService) SearchUserById(user *User) (map[string]interface{}, string, error) {
	returnMap := map[string]interface{}{} //返回值
	userMap, err := s_userById(user)
	returnMap["user"] = userMap
	return returnMap, "100", err
}

func (service userService) SearchUsers(user *User) (map[string]interface{}, string, error) {
	returnMap := map[string]interface{}{} //返回值
	userList, err := s_user(user)
	returnMap["users"] = userList
	return returnMap, "100", err
}
func (service userService) AddUser(user *User) (map[string]interface{}, string, error) {
	returnMap := map[string]interface{}{} //返回值
	id := idGenClient.GetUniqueId()       //"0" + fmt.Sprint(time.Now().Unix()) //
	user.UserId = "U" + id
	user.MtalkNo = id
	err := i_user(user)

	userMap := map[string]string{}
	userMap["userId"] = user.UserId
	returnMap["user"] = userMap
	return returnMap, "100", err
}
func (service userService) UpdateUser(user *User) (map[string]interface{}, string, error) {
	returnMap := map[string]interface{}{} //返回值
	err := u_user(user)
	return returnMap, "100", err
}

// csv
package main

import (
	//"encoding/csv"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	list := [][]string{}
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		a := strings.Split(line, "	")
		list = append(list, a)

	}
	fmt.Println("数据总长度：", len(list))
	var sig_s bytes.Buffer
	sig_s.WriteString("INSERT INTO USER(`userId`,`phoneNo`,`nickName`,`sex`,`homeAddress`,`mtalkNo`,`Wechat_uid`,`Wechat_name`,`Wechat_iconurl`,`Wechat_gender`,`createTime`,`updateTime`,`id`) VALUES")
	for i := 0; i < len(list); i++ {
		a := list[i]
		if a[8] != "" { //排除没有unionId的用户
			id := a[0]
			userId := a[1]
			nickName := a[2]
			Wechat_iconurl := a[3]
			phoneNo := a[4]
			sex_z := a[5]
			homeAddress := a[6] + " " + a[7]
			Wechat_uid := a[8]
			Wechat_uid = strings.Split(Wechat_uid, "\n")[0]
			//fmt.Println(Wechat_uid)

			sex := "男"
			if sex_z == "2" {
				sex = "女"
			}

			if len([]rune(nickName)) > 5 {
				nameRune := []rune(nickName)
				nickName = string(nameRune[:5])
			}

			sig_s.WriteString("('")
			sig_s.WriteString(userId)
			sig_s.WriteString("','")
			sig_s.WriteString(phoneNo)
			sig_s.WriteString("','")
			sig_s.WriteString(nickName)
			sig_s.WriteString("','")
			sig_s.WriteString(sex)
			sig_s.WriteString("','")
			sig_s.WriteString(homeAddress)
			sig_s.WriteString("','")
			sig_s.WriteString(userId) //mtalkNo
			sig_s.WriteString("','")
			sig_s.WriteString(Wechat_uid)
			sig_s.WriteString("','")
			sig_s.WriteString(nickName) //Wechat_name
			sig_s.WriteString("','")
			sig_s.WriteString(Wechat_iconurl)
			sig_s.WriteString("','")
			sig_s.WriteString(sex) //Wechat_gender
			sig_s.WriteString("',")
			sig_s.WriteString("NOW()")
			sig_s.WriteString(",")
			sig_s.WriteString("NOW()")
			sig_s.WriteString(",'")
			sig_s.WriteString(id)
			sig_s.WriteString("')")

			if i%5000 == 0 {
				sig_s.WriteString(";")
				sig_s.WriteString("\n")
				sig_s.WriteString("INSERT INTO USER(`userId`,`phoneNo`,`nickName`,`sex`,`homeAddress`,`mtalkNo`,`Wechat_uid`,`Wechat_name`,`Wechat_iconurl`,`Wechat_gender`,`createTime`,`updateTime`,`id`) VALUES")
			} else {
				sig_s.WriteString(",")
			}
			//fmt.Println(i)
		}
	}
	txt := sig_s.String()

	f, err2 := os.Create("userData3.sql") //创建文件
	fmt.Println("err2:", err2)
	n, err1 := io.WriteString(f, txt) //写入文件(字符串)
	fmt.Println("n:", n)
	fmt.Println("err1:", err1)
}

package main

import (
	"errors"
	"fmt"
	"testing"
	user "user/client"

	//	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	//	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

type idMock struct {
	mock.Mock
}

func (m *idMock) DeleteUser(*user.UserRequest) (*user.UserResponse, error) {
	//args := m.Called()
	return &user.UserResponse{}, errors.New("")

}
func (m *idMock) AddUser(*user.UserRequest) (*user.AddUsersResponse, error) {
	//args := m.Called()
	return &user.AddUsersResponse{UserId: "1111111"}, errors.New("")

}
func (m *idMock) SearchUser(*user.UserRequest) (*user.UserResponse, error) {
	//args := m.Called()
	return &user.UserResponse{}, errors.New("")

}
func (m *idMock) SearchUsers(*user.UserRequest) (*user.UsersResponse, error) {
	args := m.Called()
	arr := []user.UserRequest{}
	arr = append(arr, user.UserRequest{PhoneNo: "17671774535", Password: "111"})
	aaa := &user.UsersResponse{Users: arr}
	return aaa, errors.New("")

}
func (m *idMock) UpdateUser(*user.UserRequest) (*user.UserResponse, error) {
	//args := m.Called()
	return &user.UserResponse{Code: "ssssssssssssssssssssssssssssssss"}, nil

}

func TestUser_1(t *testing.T) {
	fmt.Println("测试开始")

	idm := new(idMock)
	userClient = idm

	var svc UserService
	svc = userService{}
	u := &User{PhoneNo: "17671774535", Password: "111"}
	code, request, err := svc.Login(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}
	fmt.Println(request)

	code, request, err = svc.Reg(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

	code, request, err = svc.ShortcutLogin(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

	code, request, err = svc.FindPassword(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

	code, err = svc.ChangePhoneNo(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

	code, request, err = svc.OtherLogin(u)
	fmt.Println("---------------------------------", code)
	if err == nil && "100" == code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

	u.OtherLogienType = "1"
	code, request, err = svc.OtherLogin(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

	code, err = svc.UpdateUser(u)
	if err != nil || "100" != code {
		t.Error("code:" + code + "," + "err:" + err.Error())
	}

}

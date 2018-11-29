package client

import (
	"qx_user/caller"
	"mtcomm/k8s"
	logger "mtcomm/log"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type UserClient interface {
	SearchUserById(request *caller.UserRequest) (*caller.UserResponse, error)
	SearchUsers(request *caller.UserRequest) (*caller.UsersResponse, error)
	AddUser(request *caller.UserRequest) (*caller.UserResponse, error)
	UpdateUser(request *caller.UserRequest) (*caller.UserResponse, error)
}

type userClient struct {
	Caller caller.UserCaller
	log    *logger.Logger
}

func NewUserClient(k8sClient k8s.K8sClient, tracer stdopentracing.Tracer, namespace string) UserClient {
	if !k8sClient.IsClusterEnv() {
		return nil
	}
	log := logger.GetDefaultLogger()
	c := caller.NewIdUserCaller(k8sClient, tracer, namespace, "UserClient", "qx_user", "8888")
	client := &userClient{
		Caller: c,
		log:    log,
	}
	return client
}

func (client *userClient) SearchUserById(request *caller.UserRequest) (*caller.UserResponse, error) {
	return client.Caller.SearchUserById(request)
}
func (client *userClient) SearchUsers(request *caller.UserRequest) (*caller.UsersResponse, error) {
	return client.Caller.SearchUsers(request)
}
func (client *userClient) AddUser(request *caller.UserRequest) (*caller.UserResponse, error) {
	return client.Caller.AddUser(request)
}
func (client *userClient) UpdateUser(request *caller.UserRequest) (*caller.UserResponse, error) {
	return client.Caller.UpdateUser(request)
}

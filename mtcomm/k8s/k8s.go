package k8s

import (
	"errors"
	"fmt"
	logger "mtcomm/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClient interface {
	GetHost(namespace string, serviceName string) (string, error)
	GetFirstPort(namespace string, serviceName string) (string, error)
	IsClusterEnv() bool
}

type k8sClient struct {
	Clientset    *kubernetes.Clientset
	isClusterEnv bool
}

func NewK8sClient() K8sClient {
	return newK8sClient()
}

func newK8sClient() *k8sClient {
	// creates the in-cluster config
	log := logger.GetDefaultLogger()
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Info(err.Error())
		return &k8sClient{isClusterEnv: false}
	}
	// creates the clientset
	clientset, err1 := kubernetes.NewForConfig(config)
	if err1 != nil {
		panic(err1.Error())
	}
	return &k8sClient{
		Clientset:    clientset,
		isClusterEnv: true,
	}
}

func (c *k8sClient) IsClusterEnv() bool {
	return c.isClusterEnv
}

func (c *k8sClient) GetHost(namespace string, serviceName string) (string, error) {
	if !c.IsClusterEnv() {
		return "localhost", nil
	}
	s, err := c.Clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
	if err != nil {
		//try again
		cli := newK8sClient()
		c.Clientset = cli.Clientset
		c.isClusterEnv = cli.isClusterEnv
		s, err = c.Clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
	}

	return s.Spec.ClusterIP, nil
}

func (c *k8sClient) GetFirstPort(namespace string, serviceName string) (string, error) {
	if !c.IsClusterEnv() {
		return "", errors.New("Not a k8s cluster. So can not get prot.")
	}
	s, err := c.Clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
	if err != nil {
		//try again
		cli := newK8sClient()
		c.Clientset = cli.Clientset
		c.isClusterEnv = cli.isClusterEnv
		s, err = c.Clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
	}

	port := s.Spec.Ports[0].Port

	return fmt.Sprint(port), nil
}

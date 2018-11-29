package main

import (
	"bytes"
)

type monitorRequest struct {
	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`
	Param       string `json:"param"`
	ErrorMsg    string `json:"errorMsg"`
	ErrorTime   string `json:"errorTime"`
}

func (m *monitorRequest) String() string {
	var b bytes.Buffer
	b.WriteString(" ServiceName: ")
	b.WriteString(m.ServiceName)
	b.WriteString(" MethodName: ")
	b.WriteString(m.MethodName)
	b.WriteString(" Param: ")
	b.WriteString(m.Param)
	b.WriteString(" ErrorMsg: ")
	b.WriteString(m.ErrorMsg)
	b.WriteString(" ErrorTime: ")
	b.WriteString(m.ErrorTime)
	return b.String()
}

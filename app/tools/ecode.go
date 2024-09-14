package tools

import "fmt"

var (
	ParamErr = ECode{
		Code:    10002,
		Message: "参数错误",
	}
)

type ECode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *ECode) String() string {
	return fmt.Sprintf("COde:%d,message:%s", e.Code, e.Message)
}

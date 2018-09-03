package testcase

import (
	_ "fmt"
	_ "log"
	"net/http"
)

//用例接口
type PItranslate interface {
	//转换接口
	Conver(string) string //将含有变量与表达式的模板翻译过来

}

type PIrequest interface {
	Buildrequest() *http.Request //构造请求体
}

type ParamsInterface interface {
	Encode() string
	ToJson() string
}

type CaseItemInterface interface {
	Getconditions() []Condition
}

type CasesetInterface interface {
	GetCases() []CaseItem
}

//用例解析器接口

type PIparserInsterface interface {
	GetCaseset() *Caseset
	GetUid() uint32
}

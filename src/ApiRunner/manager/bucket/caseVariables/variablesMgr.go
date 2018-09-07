// 全局管理用例变量
package caseVariables

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

//TODO 目前变量只支持字符串和number类型，后续可能需要支持列表
var varsMap sync.Map

func getKey(uid uint32, k string) string {
	sUid := fmt.Sprint(uid)
	return strings.Join([]string{sUid, k}, `-`)
}

func GetVar(uid uint32, k string) string {
	v, ok := varsMap.Load(getKey(uid, k))
	if ok {
		return v.(string)
	}
	return ""
}

func SetVar(uid uint32, k, v string) {
	varsMap.Store(getKey(uid, k), v)
}

func DelVars(uid uint32) {
	// 未真正删除，具体删除由sync.map内部管理
	sUid := fmt.Sprint(uid)
	varsMap.Range(func(k, v interface{}) bool {
		if strings.Index(k.(string), sUid) != -1 {
			varsMap.Delete(k)
			log.Println("delete key/value:", k, v)
		}
		return true
	})
}

type VarMap struct {
	Uid uint32
}

func (this VarMap) GetData() map[string]interface{} {
	return GetVarsMap(this.Uid)
}

func GetVarsMap(uid uint32) map[string]interface{} {
	m := make(map[string]interface{})
	varsMap.Range(func(k, v interface{}) bool {
		i := strings.Index(k.(string), `-`)
		_k := string([]byte(k.(string))[i+1:])
		m[_k] = v.(string)
		return true
	})
	return m
}

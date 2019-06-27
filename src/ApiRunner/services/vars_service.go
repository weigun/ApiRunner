// vars_service.go
package services

import (
	"ApiRunner/dao"
	"strings"

	//"ApiRunner/models"
	//"encoding/json"
	"errors"
	"fmt"
)

/*
用例变量管理器
存放用例中的模板变量、导出变量和全局变量
*/
type VarsManager struct {
	cache dao.Cache
}

var VarsMgr *VarsManager

const (
	DefaultVarsMgrKey = `VarsMgr`
	DefaultTimeOut    = 3600
)

func (vm *VarsManager) Add(key string, val interface{}) error {
	key = fmt.Sprintf(`%s:%s`, DefaultVarsMgrKey, key)
	err := vm.cache.Put(key, val, DefaultTimeOut) //防止长期堆积
	if err != nil {
		return errors.New(fmt.Sprintf(`can't not add var,key %s,val %v`, key, val))
	}
	return nil
}

func (vm *VarsManager) Delete(key string) error {
	key = fmt.Sprintf(`%s:%s`, DefaultVarsMgrKey, key)
	err := vm.cache.Delete(key)
	if err != nil {
		return errors.New(fmt.Sprintf(`delete var failed %s`, err.Error()))
	}
	return nil
}

func (vm *VarsManager) Get(key string) string {
	key = fmt.Sprintf(`%s:%s`, DefaultVarsMgrKey, key)
	val := vm.cache.Get(key)
	return val
}

func (vm *VarsManager) GetByGroup(groupKey string) map[string]string {
	key := fmt.Sprintf(`%s:%s`, DefaultVarsMgrKey, groupKey)
	m := make(map[string]string)
	splitChar := `:`
	for _, k := range vm.cache.Keys(fmt.Sprintf(`%s*`, key)) {
		_k := strings.Split(k, splitChar)[strings.Count(k, splitChar)]
		m[_k] = vm.cache.Get(fmt.Sprintf(`%s:%s`, key, _k))
	}
	return m
}

func init() {
	VarsMgr = &VarsManager{dao.GetCache()}
}

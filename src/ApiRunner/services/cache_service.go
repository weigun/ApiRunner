// cache_service.go
package services

import (
	"ApiRunner/dao"
	"ApiRunner/models"
	"encoding/json"
	"errors"
	"fmt"
)

/*
缓存管理器，目前只有对用例对象进行缓存
*/
type CacheManager struct {
	cache dao.Cache
}

const DefaultObjCacheKey = `ObjMgr`

func (cm *CacheManager) CacheObj(obj models.Executable) error {
	key := obj.GetName()
	val := obj.Json()
	key = fmt.Sprintf(`%s:%s`, DefaultObjCacheKey, key)
	err := cm.cache.Put(key, val, 0)
	if err != nil {
		return errors.New(`can't not add testcase to cache`)
	}
	return nil
}

func (cm *CacheManager) FromCache(key string) (obj models.Executable, err error) {
	key = fmt.Sprintf(`%s:%s`, DefaultObjCacheKey, key)
	jsonVal := cm.cache.Get(key)
	if jsonVal != `` {
		e := json.Unmarshal([]byte(jsonVal), &obj)
		if e != nil {
			err = e
			return
		}
		return
	}
	err = errors.New(`caseobj not found`)
	return
}

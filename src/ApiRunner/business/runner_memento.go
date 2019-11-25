package business

import (
	"fmt"
	"reflect"
	"sync"
)

type Memento interface {
	GetState() interface{}
}

type memento struct {
	state interface{}
}

func (m *memento) GetState() interface{} {
	return m.state
}

type mementoList []Memento

type mementoMgr struct {
	sync.RWMutex
	mementoList
}

func (m *mementoMgr) SaveMemento(IMem Memento) {
	m.Lock()
	defer m.Unlock()
	m.mementoList = append(m.mementoList, IMem)
}

//弹出最近的一个与参数相同类型的备忘对象，并返回之。该方法会从原位置中将对象删除
func (m *mementoMgr) PopMementoWith(typ interface{}) Memento {
	m.Lock()
	defer m.Unlock()
	targetType := reflect.TypeOf(typ).Elem().Name()
	log.Debug(fmt.Sprintf(`targettype:%s,%T`, targetType, reflect.TypeOf(typ).Elem()))
	// for i, tmp := range m.mementoList {
	for i := len(m.mementoList) - 1; i >= 0; i-- {
		// 模拟栈先进后出
		tmp := m.mementoList[i]
		log.Debug(fmt.Sprintf(`tmp type:%s`, reflect.TypeOf(tmp.GetState()).Elem().Name()))
		if reflect.TypeOf(tmp.GetState()).Elem().Name() == targetType {
			m.mementoList = append(m.mementoList[0:i], m.mementoList[i+1:]...)
			return tmp
		}
	}
	return nil
}

func (m *mementoMgr) Clean() int {
	m.Lock()
	defer m.Unlock()
	num := len(m.mementoList)
	m.mementoList = mementoList{}
	return num

}

func NewMementoMgr() *mementoMgr {
	return &mementoMgr{mementoList: mementoList{}}
}

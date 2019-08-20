package refs

import (
	// "fmt"

	"ApiRunner/models"
)

type Node interface {
	Name() string
	Parent() *rnode
	ValueOf(string) interface{}
	ChildAt(int) *rnode
	SetParent(*rnode)
	AddChild(*rnode)
	AddPairs(string, interface{})
}

type rnode struct {
	name     string
	vars     models.Variables
	parent   *rnode
	children []*rnode
}

func (n *rnode) Name() string {
	return n.name
}

func (n *rnode) Parent() *rnode {
	return n.parent
}

func (n *rnode) ValueOf(varName string) interface{} {
	if val, ok := n.vars[varName]; ok {
		return val
	}
	return nil
}

func (n *rnode) ChildAt(index int) *rnode {
	if index > len(n.children) {
		panic(`IndexError: list assignment index out of range`)
	}
	return n.children[index]
}

func (n *rnode) SetParent(parent *rnode) {
	n.parent = parent
}

func (n *rnode) AddChild(child *rnode) {
	n.children = append(n.children, child)
	child.SetParent(n)
}

func (n *rnode) AddPairs(key string, val interface{}) {
	n.vars[key] = val
}

func New(name string) *rnode {
	return &rnode{name: name, vars: models.Variables{}}
}
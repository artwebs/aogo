package web

import (
	"regexp"
	"strings"
)

type RouterTree struct {
	prefix    string
	child     []*RouterTree
	runObject interface{}
}

func NewRouterTree() *RouterTree {
	return &RouterTree{}
}

func (this *RouterTree) AddRouter(prefix string, obj interface{}) {
	parr := this.splitPath(prefix)
	this.addSeq(parr, obj, this)
}

func (this *RouterTree) FindRouter(url string) interface{} {

	return this.findRouter(this.splitPath(url), this)
}

func (this *RouterTree) addSeq(parr []string, obj interface{}, tree *RouterTree) {
	if len(parr) > 0 {
		var child *RouterTree
		for _, item := range tree.child {
			if parr[0] == item.prefix {
				child = item
				break
			}
		}
		if child == nil {
			child = NewRouterTree()
			child.prefix = parr[0]
			tree.child = append(tree.child, child)
			this.addSeq(parr[1:], obj, child)
		} else {
			this.addSeq(parr[1:], obj, child)
		}
	} else {
		tree.runObject = obj
	}
}

func (this *RouterTree) findRouter(parr []string, tree *RouterTree) interface{} {
	if len(parr) > 0 {
		for _, item := range tree.child {
			if parr[0] == item.prefix {
				return this.findRouter(parr[1:], item)
			}
		}
	}
	return tree.runObject
}

func (this *RouterTree) splitPath(prefix string) []string {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return []string{}
	}
	return strings.Split(prefix, "/")
}

type TreeReg struct {
	regexps   *regexp.Regexp
	keys      []string
	runObject interface{}
}

func (this *TreeReg) NewTreeInfo() *TreeReg {
	return &TreeReg{}
}

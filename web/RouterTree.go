package web

import (
	"regexp"
	"strings"

	"github.com/artwebs/aogo/logger"
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

func (this *RouterTree) AddTree(prefix string, tree *RouterTree) {
	parr := this.splitPath(prefix)
	this.addTree(parr, tree, this)
}

func (this *RouterTree) FindRouter(url string) ([]string, interface{}) {

	return this.findRouter(this.splitPath(url), this)
}
func (this *RouterTree) addTree(parr []string, sourceTree, targetTree *RouterTree) {
	if len(parr) > 0 {
		var child *RouterTree
		for _, item := range targetTree.child {
			if parr[0] == item.prefix {
				child = item
				break
			}
		}
		if child == nil {
			child = NewRouterTree()
			child.prefix = parr[0]
			targetTree.child = append(targetTree.child, child)
		}
		this.addTree(parr[1:], sourceTree, child)
	} else {
		targetTree.child = append(targetTree.child, sourceTree.child...)
		targetTree.runObject = sourceTree.runObject
	}
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
		}
		this.addSeq(parr[1:], obj, child)
	} else {
		tree.runObject = obj
	}
}

func (this *RouterTree) findRouter(parr []string, tree *RouterTree) ([]string, interface{}) {
	if len(parr) > 0 {
		for _, item := range tree.child {
			if parr[0] == item.prefix {
				return this.findRouter(parr[1:], item)
			}
		}
	}
	return parr, tree.runObject
}

func (this *RouterTree) splitPath(prefix string) []string {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return []string{}
	}
	return strings.Split(prefix, "/")
}

func (this *RouterTree) PrintTree(pprefix string) {
	logger.InfoTag(this, pprefix+this.prefix, this.runObject)
	if c := this.child; c != nil {
		for _, s := range c {
			tmp := pprefix + this.prefix + "/"
			s.PrintTree(tmp)
		}
	}
}

type TreeReg struct {
	regexps   *regexp.Regexp
	keys      []string
	runObject interface{}
}

func (this *TreeReg) NewTreeInfo() *TreeReg {
	return &TreeReg{}
}

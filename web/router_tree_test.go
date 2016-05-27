package web

import (
	"log"
	"strings"
	"testing"
)

func init() {

}

func TestAddRouter(t *testing.T) {
	child := NewRouterTree()
	child.AddRouter("/", "/manage")
	child.AddRouter("/user/add", "/manage/user/add")
	child.AddRouter("/user/del", "/manage/user/del")

	child1 := NewRouterTree()
	child1.AddRouter("/user/add", "/admin/dev/user/add")
	child1.AddRouter("/user/del", "/admin/dev/user/del")
	child1.AddRouter("/user/edit", "/admin/dev/user/edit")

	root := NewRouterTree()
	root.AddRouter("/", "/")
	root.AddRouter("/index", "/index")
	root.AddRouter("/index/add", "/index/add")
	root.AddRouter("/index/del", "/index/del")
	root.AddRouter("/admin/add", "/admin/add")
	root.AddRouter("/admin/del", "/admin/del")
	root.AddRouter("/admin/", "/admin")
	root.AddTree("/manage", child)
	root.AddTree("/admin/dev", child1)
	printTree("/", root, t)

	findTree("/admin/del", root, t)
	findTree("/admin/edit", root, t)
}

func printTree(before string, tree *RouterTree, t *testing.T) {
	if tree.runObject != nil {
		if strings.Trim(before, "/") != strings.Trim((tree.runObject).(string), "/") {
			t.Errorf("%s!=%s", before, tree.runObject)
		}
	}
	log.Print("prefix=", before, ",runObject=", tree.runObject)
	for _, item := range tree.child {
		url := before + item.prefix + "/"
		printTree(url, item, t)
	}
}

func findTree(url string, tree *RouterTree, t *testing.T) {
	keys, obj := tree.FindRouter(url)
	log.Println("keys=", keys, ",obj=", obj, ",url=", url)
}

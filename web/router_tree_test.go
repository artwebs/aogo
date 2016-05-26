package web

import (
	"log"
	"strings"
	"testing"
)

func init() {

}

func TestAddRouter(t *testing.T) {
	root := NewRouterTree()
	root.AddRouter("/", "/")
	root.AddRouter("/index", "/index")
	root.AddRouter("/index/add", "/index/add")
	root.AddRouter("/index/del", "/index/del")
	root.AddRouter("/admin/add", "/admin/add")
	root.AddRouter("/admin/del", "/admin/del")
	root.AddRouter("/admin/", "/admin")
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
	log.Println(tree.FindRouter(url), "==", url)
}

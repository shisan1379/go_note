package msgo

import (
	"fmt"
	"testing"
)

func Test_treeNode_Put(t1 *testing.T) {
	node := treeNode{name: "/", children: make([]*treeNode, 0)}
	node.Put("/user/**")
	node.Put("/user/add/:id")
	node.Put("/user/getMap/delete/:id")
	fmt.Printf("%v", node)

	get := node.Get("/user/**")
	fmt.Printf("%v", get)
}

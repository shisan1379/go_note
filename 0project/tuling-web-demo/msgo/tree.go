package msgo

import (
	"strings"
)

type treeNode struct {
	name       string
	children   []*treeNode
	routerName string
	isEnd      bool
}

func newTreeNode() *treeNode {
	return &treeNode{
		name:     "/",
		children: make([]*treeNode, 0),
	}
}

// put path: /user/get/:id
func (t *treeNode) Put(path string) {
	//记录 根节点 指针
	root := t
	strs := strings.Split(path, "/")
	for i, name := range strs {
		if i == 0 {
			continue
		}

		// 对子节点进行匹配，如果存在则不添加，并进入下一个节点
		children := t.children
		notExist := true
		for _, node := range children {
			if node.name == name {
				t = node //当匹配到当前节点后，步入下个节点，在下次循环中再次进行适配
				notExist = false
				break
			}
		}

		//创建节点并添加到子节点中，进入下个子节点
		if notExist {

			isEnd := false
			if i == len(strs)-1 {
				isEnd = true
			}

			node := &treeNode{name: name, children: make([]*treeNode, 0), isEnd: isEnd}
			children = append(children, node)
			t.children = children
			t = node
		}
	}

	// 将根节点给到外面
	t = root
}

// get path: /user/get/1
func (t *treeNode) Get(path string) *treeNode {
	strs := strings.Split(path, "/")
	routerName := ""
	for index, name := range strs {
		if index == 0 {
			continue
		}
		children := t.children
		isMatch := false
		for _, node := range children {
			// 比较节点名称
			// 节点是够为 *
			// 节点是够包含 :
			if node.name == name ||
				node.name == "*" ||
				strings.Contains(node.name, ":") {
				isMatch = true
				routerName += "/" + node.name
				node.routerName = routerName
				t = node
				if index == len(strs)-1 {
					return node
				}
				break
			}
		}
		if !isMatch {
			for _, node := range children {
				// /user/**
				// /user/get/userInfo
				// /user/aa/bb
				if node.name == "**" {
					routerName += "/" + node.name
					node.routerName = routerName
					return node
				}
			}

		}
	}
	return nil
}

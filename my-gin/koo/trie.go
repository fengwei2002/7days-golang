package koo

import (
	"fmt"
	"strings"
)

// trie 树中的基础 node 节点
// pattern: 待匹配路由，例如 /p/:lang
// /p/:lang/do c只有在第三层节点，即 doc 节点，pattern 才会设置为 /p/:lang/doc。p 和 :lang 节点的 pattern 属性皆为空
// part: 路由中的一部分，例如 :lang
// children: 子节点，例如 [doc, tutorial, intro]
// isWild: 是否精确匹配，part 含有 : 或 * 时为true
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

// node.String() 方法，将这个 node 的信息进行输出
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}


// insert 函数在 trie 树中插入一个 parts 的对应一系列 node
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	} // 如果到达了最后一层，则更新 n 节点的 pattern 属性

	part := parts[height]       // 取出这一层需要匹配的路由部分 a part of route
	child := n.matchChild(part) // 取出匹配 this part 的第一个 children in n
	if child == nil {
		// 如果不存在匹配 part 的 child
		// 创建一个新的 node，包含这个 part 的 isWild 属性， 和这个 part 的具体 value
		child = &node{
			part:     part,
			isWild:   part[0] == ':' || part[0] == '*',
		}
		// 将创建的新的 node 添加到 n.children 的后面
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1) // 递归执行，直到所有 part in parts 插入完毕
	// 也不会覆盖已经存在的路由，只是开辟新的道路，有之前的路就尽量走之前的路了
}

// search 寻找一个 parts 数组是否可以匹配一个最终的 node
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		// 如果到达 parts 的末尾端点，或者 n.part 是 * 开头的时候， 如果 n 存在结尾 pattern 则返回
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height + 1)
		if result != nil {
			return result
		}
	} // 递归搜索下一层，直到完全匹配返回对应节点，否则返回

	return nil
}

// travel 函数将 n 节点的一整棵 trie 树都放到 list 中
func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
} // 将 node 中的所有的 children 作为一个 list 返回

// matchChild 遍历 n 的所有 children，找到一个模糊匹配，或者精确匹配的 node 则返回（返回的是第一个匹配的节点）
// 否则返回空节点，代表 n 的后面不存在可以匹配 part 的路由
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 找出 n 节点中所有可以匹配 part 的路由，拼接为一个 slice 返回
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
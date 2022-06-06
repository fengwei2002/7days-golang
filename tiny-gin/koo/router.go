package koo

import (
	"net/http"
	"strings"
)

// 将 router 部分进行独立，方便在 router 中添加功能
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots    key eg, roots['GET']                  roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

// newRouter 是 router 的构造函数，返回一个 router
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 传入一个完整的 link，转换为一个 parts 数组
func parsePattern(pattern string) []string {
	pp := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range pp {
		if item != "" { // 剔除是空格的，遇到 * 之后截断
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 提供接口，method， pattern 和 handler 参数
// 将信息打印到 log 中，同时，将路由信息保存到 router 的 trie node 对应的 tree 中
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	if _, ok := r.roots[method]; ok == false {
		r.roots[method] = &node{}
	} // 如果 method 没有根 node，先创建根 node

	r.roots[method].insert(pattern, parts, 0)
	// 将这个 parts 插入到 method 对应的 trie 中
	r.handlers[key] = handler
	// 在 r 中存储具体的 key 对应的 handlerFunc
}

// getRoute 根据路由的方法，以及具体的 routePath 得到对应的 node 以及对应的 map 解析结果
// /:lang, /go -> {lang: go}
// /static/css/background.css 匹配到 static/*filepath
// 返回的对应 map key value 是: {filepath: "css/background.css"}
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if ok == false {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	// 使用 search 方法，搜索匹配 parts 的 node

	if n != nil { // 如果存在对应的 node
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
				// /:lang, /go -> {lang: go}
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				// /static/css/background.css 匹配到 static/*filepath
				// 返回的对应 map key value 是: {filepath: "css/background.css"}
			}
		}
		return n, params
	}
	return nil, nil
}

// getRoutes 方法，返回 method 对应的所有的 node
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// handle 传入上下文，使用 r 内存储的 handlers 进行处理
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" +  n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key]) // 添加到 c.Handlers 后面
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
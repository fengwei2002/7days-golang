package koo

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc 是路由处理函数类型的一个缩写
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 实现了 ServerHTTP 方法，可以作为 ListenAndServer 的第二个参数使用
// - Engine.router 将 string 类型的路由映射到一个路由处理函数
type Engine struct {
	router map[string]HandlerFunc
}

// New 是 koo.Engine 的构造函数
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc, 0),
	}
}

// addRoute 将一个路由的信息存储到 engine 的 router 中
// 第一个参数是使用的方法，第二个参数是具体的路由，第三个参数是路由处理函数
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET 使用 GET 方法调用 addRoute 不用在参数列表中声明 GET 而是使用 GET 方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}


// POST 使用 POST 方法调用 addRoute 不用在参数列表中声明 POST 而是使用 POST 方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 使用 run 接口，将传入的 addr 使用  http ListenAndServe 运行，第二个参数的 engine 已经实现 ServeHTTP 方法
// 所有的路由请求交给 engine 处理
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok == true {
		handler(w, req)
		// 使用对应的 handler function 处理这个 w 和 req
	} else {
		log.Fatal(fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL))
	}
}
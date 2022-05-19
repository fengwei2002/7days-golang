package koo

import (
	"log"
	"time"
)

// 添加中间件 logger 记录请求到响应中间花费的时间

// 中间件的定义和路由映射函数的 handler 一样，处理的输入是 Context 对象

// 中间件是应用在RouterGroup上的，应用在最顶层的 Group，相当于作用于全局，所有的请求都会被中间件处理。
// 那为什么不作用在每一条路由规则上呢？作用在某条路由规则，那还不如用户直接在 Handler 中调用直观。
// 只作用在某条路由规则的功能通用性太差，不适合定义为中间件

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
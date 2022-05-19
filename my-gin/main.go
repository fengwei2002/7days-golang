package main

import (
	"koo"
	"net/http"
)

func main()  {
	r := koo.Default()
	// 在 koo.go 中，使用 Default 在 New() 上额外挂载 logger.go 和 recovery.go 两个中间件

	r.GET("/", func(c *koo.Context) {
		c.String(http.StatusOK, "hello fengwei\n")
	})

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *koo.Context) {
		names := []string{"fengwei"}
		c.String(http.StatusOK, names[100])
	})

	r.Run("localhost:8080")
}
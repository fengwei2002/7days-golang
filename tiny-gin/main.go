package main

import (
	"fmt"
	"html/template"
	"koo" // use go work to test (!!! after 1.18)
	"log"
	"net/http"
	"time"
)

func onlyForV2() koo.HandlerFunc {
	return func(c *koo.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

type student struct {
	Name string
	Age  int8
}

func main() {
	r := koo.Default()
	// 在 koo.go 中，使用 Default 在 New() 上额外挂载 logger.go 和 recovery.go 两个中间件

	r.GET("/", func(c *koo.Context) {
		c.String(http.StatusOK, "hello fengwei\n")
	})
	// string 返回

	r.POST("/login", func(c *koo.Context) {
		c.JSON(http.StatusOK, koo.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	// JSON 返回

	r.GET("/hello/:name", func(c *koo.Context) {
		// expect /hello/fengwei
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	// trie 树提供的路由匹配

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *koo.Context) {
			c.String(http.StatusOK, "hello fengwei\n")
		})

		v1.GET("/hello", func(c *koo.Context) {
			// expect /hello?name=fengwei
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	// group 路由控制

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *koo.Context) {
			// expect /hello/fengwei
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	// 在 context 中添加的中间件功能

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *koo.Context) {
		names := []string{"fengwei"}
		c.String(http.StatusOK, names[100])
	})

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "fengwei", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *koo.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *koo.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", koo.H{
			"title":  "koo",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *koo.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", koo.H{
			"title": "koo",
			"now":   time.Date(2022, 7, 9, 0, 0, 0, 0, time.UTC),
		})
	})

	// 使用 template 功能

	r.Run("localhost:8080")
}

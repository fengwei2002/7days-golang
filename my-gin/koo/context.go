package koo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 是将 string 映射到任意类型的一个简写
// 在构建 JSON 类型的数据的时候，显得更加简洁
type H map[string]any

type Context struct {
	// context 是一个百宝箱 对外提供一个接口，功能在 context 上进行扩展
	// 源参数
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求的信息
	Path   string
	Method string
	Params map[string]string // 将路由解析后的参数存储到 Params 中
	// 返回的信息
	StatusCode int
	// 自己添加的中间件
	handlers []HandlerFunc // 每个 Context 一组 handlerFunc
	index    int           // 代表当前执行到了哪一个 handlerFunc

	// template
	engine *Engine // Engine Pointer
}

// newContext 是 context 的构造函数，返回一个 context 对象
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next 方法，对于一个 context，处理从 index 开始之后所有的 handlerFunc
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail 方法将 c 的 index 跳转到最后一个元素的下一个，然后，将错误以 JSON 格式返回
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	fmt.Println("------------------------------------------------------")
	c.JSON(code, H{"message": err})
}

// Param 返回 map[key] 对应的 value
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 接收一个 string 返回 http.Request.FormValue(string) 的结果
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 接收一个 string 返回 http.Request.URL.Query().Get(string) 的结果
// ?name=fengwei Query(name) 返回 fengwei
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 传入一个 code，然后将 StatusCode 记录到 c 中
// 同时在 writer.WriteHeader 中写入这个 code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置 c 中的 writer.Header 中的 key 对应的具体 value
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

/*使用上面提供的接口实现更加集中的 API 接口 */

// 使用 String 传入一部分信息，然后将 code 的信息记录在 c 中
// 并且将 values... 中的内容，转换为 byte 数组，写入到 c.Writer 中
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 方法，传入 code 和 具体的 json 内容
// 将 json 的基本信息写入到 c 中，然后将 json 的内容写入也写入到 writer 中
// 解析 json 出现错误的话 返回 500 error
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 接口，传入一个 []byte 类型的 data ，直接写入 c.Writer 中
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 接口，根据模板文件名选择模板进行渲染。
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

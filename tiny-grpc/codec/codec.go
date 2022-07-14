package codec

import (
	"io"
)

/*
一个典型的 rpc 调用：
err := client.Call("ServiceName.MethodName", args, &reply)
服务端的响应包含 error 和返回值 reply 两个，

将请求和响应中的参数和返回值抽象为 body 剩下的信息放在 header 中
*/

type Header struct {
	ServiceMethod string // format "Service.Method" 服务名和方法名，
	Seq           uint64 // sequence number chosen by client 某个请求的 id 用来区分不同的请求
	Error         string // error message 客户端置为空，如果服务端发生错误，就将信息放到 error 中
}

type Codec interface { // 对消息体进行编码和解码，抽象为接口就可以实现不同的 codec 实例
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

// 使用一个 map 来存储每一种类型和每一种编解码方式

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}

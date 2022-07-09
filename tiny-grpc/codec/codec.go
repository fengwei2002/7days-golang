package codec

/*

一个典型的 rpc 调用如下
err = client.Call("servername.methodName", args, &reply)

客户端发送的请求包括 服务名, 方法名，参数

服务端的响应内容包括错误 error 返回值 reply 两个

将请求和响应中的参数和返回值抽象为 body 剩余的信息放在 header 中

*/
import (
	"io"
)

type Header struct {
	ServiceMethod string // 服务名和方法名，通常和 go 中的结构体和方法进行映射
	Seq           uint64 // 请求的序号，可以认为是某一个请求的 ID 用于区分不同的请求
	Error         string // 客户端置为 空，服务端如果发生错误，将错误信息写在 error 中
}

// Codec 对消息体进行解码的接口 Codec
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(closer io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"  // 主要实现
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}

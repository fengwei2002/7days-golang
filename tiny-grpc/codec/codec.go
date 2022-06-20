package codec

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
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
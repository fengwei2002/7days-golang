// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tinygrpc

import "tinygrpc/codec"

/*

客户端和服务端的通信需要协商一些内容，例如 http 报文
分为 header 和 body 两个部分
body 的格式和长度用过 header 中的 content-type 和 content-length 指定
服务端通过解析 header 就可以知道如何从 body 中读取需要的信息

对于 rpc 来说这部分的功能是自己实现的，为了提升性能，一般在报文的开始会来协商相关的信息
目前需要协商的唯一内容是消息的编码和解码方式
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

// Option 中规定需要具体协商的内容
// 服务端首先使用 json 解码 option 中的内容，然后通过 option 的 codeType 解码剩余的内容
// | Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
// | <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
//
// 在一次连接中，option 在报文的最开始的地方，header 和 body 可以由很多个
// | Option | Header1 | Body1 | Header2 | Body2 | ...
type Option struct {
	MagicNumber int        // MagicNumber marks this's a geerpc request
	CodecType   codec.Type // client may choose different Codec to encode body
}

// DefaultOption 定义默认使用的方式
var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server represents an RPC Server.
type Server struct{}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// 首先使用 json.NewDecoder 反序列化得到 option 实例
// 检查 MagicNumber 和 CodeType 的值是否正确
// 然后根据 CodeType 得到对应的 消息编码和解码器
// 接下来的处理交给 serverCodec
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn))
}

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

// serverCodec
// 主要三个阶段
// 读取请求 readRequest
// 处理请求 handleRequest
// 恢复请求 sendResponse
func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {                      // 一次连接中，允许接收多个请求，requestHeader 和 requestBody 所以这里使用 for 来无限制的接收请求的到来，直到发生错误
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg) // 使用协程并发执行请求
	}
	wg.Wait()
	// 处理请求是并发的，但是回复请求的报文必须是逐个发送的，并发容易导致多个恢复报文交织在一次
	// 客户端无法解析， 使用 锁来保证
	_ = cc.Close()
}

// request stores all information of a call
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	// TODO: now we don't know the type of request argv
	// day 1, just suppose it's string
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyv
	// day 1, just print argv and send a hello message
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("tinyrpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
// net.Listener 作为参数，for 循环等待 socket 的连接建立，并且开启子协程处理
// 处理过程交给了 serverConn 方法
// 如果想要启动服务：
// lis, _ := net.Listen("tcp", ":9999")
// tinygrpc.Accept(lis)
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

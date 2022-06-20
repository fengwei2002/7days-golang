package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// GobCodec gob 类型的解码器
type GobCodec struct {
	conn io.ReadWriteCloser // 通常是通过 tcp 或者 unix 建立 socket 时候得到的连接实例
	buf  *bufio.Writer      // 为了防止阻塞而创建的 带有缓冲的 writer
	dec  *gob.Decoder       // gob 中的 decoder
	enc  *gob.Encoder       // gob 中的 encoder
}

// 以下函数提供 Codec 所需要的接口

func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body any) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body any) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}

	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

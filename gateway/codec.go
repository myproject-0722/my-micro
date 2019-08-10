package gateway

import (
	"encoding/binary"
	"errors"
	"net"
	"time"

	packet "github.com/myproject-0722/my-micro/proto/packet"
)

type Codec struct {
	Conn     net.Conn
	ReadBuf  buffer // 读缓冲
	WriteBuf []byte // 写缓冲
}

var ErrOutOfSize = errors.New("package content out of size") // package的content字节数组过大

// newCodec 创建一个解码器
func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:     conn,
		ReadBuf:  newBuffer(conn, packet.BufLen),
		WriteBuf: make([]byte, packet.BufLen),
	}
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader()
}

// Decode 解码数据
func (c *Codec) Decode() (*packet.Package, bool) {
	var err error
	// 读取数据类型
	typeBuf, err := c.ReadBuf.seek(0, packet.TypeLen)
	if err != nil {
		return nil, false
	}

	// 读取数据长度
	lenBuf, err := c.ReadBuf.seek(packet.TypeLen, packet.HeadLen)
	if err != nil {
		return nil, false
	}

	// 读取数据内容
	valueType := int32(binary.BigEndian.Uint32(typeBuf))
	valueLen := int(binary.BigEndian.Uint32(lenBuf))

	valueBuf, err := c.ReadBuf.read(packet.HeadLen, valueLen)
	if err != nil {
		return nil, false
	}
	message := packet.Package{Code: valueType, Content: valueBuf}
	return &message, true
}

// Eecode 编码数据
func (c *Codec) Eecode(pack packet.Package, duration time.Duration) error {
	contentLen := len(pack.Content)
	if contentLen > packet.ContentMaxLen {
		return ErrOutOfSize
	}

	binary.BigEndian.PutUint32(c.WriteBuf[0:packet.TypeLen], uint32(pack.Code))
	binary.BigEndian.PutUint32(c.WriteBuf[packet.TypeLen:packet.HeadLen], uint32(len(pack.Content)))
	copy(c.WriteBuf[packet.HeadLen:], pack.Content[:contentLen])

	c.Conn.SetWriteDeadline(time.Now().Add(duration))
	_, err := c.Conn.Write(c.WriteBuf[:packet.HeadLen+contentLen])
	if err != nil {
		return err
	}
	return nil
}

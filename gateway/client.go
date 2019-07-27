package gateway

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro"
	user "github.com/myproject-0722/my-micro/proto/user"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

const (
	ReadDeadline  = 10 * time.Minute
	WriteDeadline = 10 * time.Second
)

// 消息协议
const (
	CodeSignIn         = 1 // 设备登录
	CodeSignInACK      = 2 // 设备登录回执
	CodeSyncTrigger    = 3 // 消息同步触发
	CodeHeadbeat       = 4 // 心跳
	CodeHeadbeatACK    = 5 // 心跳回执
	CodeMessageSend    = 6 // 消息发送
	CodeMessageSendACK = 7 // 消息发送回执
	CodeMessage        = 8 // 消息投递
	CodeMessageACK     = 9 // 消息投递回执
)

type Client struct {
	Codec      *Codec // 编解码器
	ReadBuffer buffer
	WriteBuf   []byte
	UserId     int64
	DeviceId   int64
	IsSignIn   bool // 是否登录
}

// Package 消息包
type Package struct {
	Code    int    // 消息类型
	Path    string //
	Method  string //
	Content []byte // 消息体
}

func NewClient(conn net.Conn) *Client {
	codec := NewCodec(conn)
	return &Client{
		Codec:      codec,
		ReadBuffer: newBuffer(conn, BufLen),
		WriteBuf:   make([]byte, BufLen),
	}
}

func (c *Client) DoConn() {
	defer RecoverPanic()

	c.HandleConnect()

	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		_, err = c.Read()
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		for {
			message, ok := c.Decode()
			if ok {
				c.HandlePackage(message)
				continue
			}
			break
		}
	}
}

// HandlePackage 处理消息包
func (c *Client) HandlePackage(pack *Package) {
	// 未登录拦截
	if pack.Code != CodeSignIn && c.IsSignIn == false {
		c.Release()
		return
	}

	switch pack.Code {
	case CodeSignIn:
		c.HandlePackageSignIn(pack)
		break
	default:
		break
	}
	/*

	  switch pack.Code {
	  case CodeSignIn:
	      c.HandlePackageSignIn(pack)
	  case CodeSyncTrigger:
	      c.HandlePackageSyncTrigger(pack)
	  case CodeHeadbeat:
	      c.HandlePackageHeadbeat()
	  case CodeMessageSend:
	      c.HandlePackageMessageSend(pack)
	  case CodeMessageACK:
	      c.HandlePackageMessageACK(pack)
	  }*/
	return
}

// HandlePackageSignIn 处理登录消息包
func (c *Client) HandlePackageSignIn(pack *Package) {
	var sign user.SignInRequest
	err := proto.Unmarshal(pack.Content, &sign)
	if err != nil {
		log.Fatal(err)
		c.Release()
		return
	}

	log.Print(sign.DeviceId, " ", sign.UserId, " ", sign.Token)

	// 推荐使用etcd集群 做为服务发现,为测试暂用consul
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// create a new service
	service := micro.NewService(micro.Registry(reg))
	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := user.NewUserService("go.mymicro.srv.user", service.Client())

	// Make request
	rsp, err := cl.SignIn(context.Background(), &sign)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.ResCode)
	/*var sign pb.SignIn
	  err := proto.Unmarshal(pack.Content, &sign)
	  if err != nil {
	      log.Fatal(err)
	      c.Release()
	      return
	  }*/
}

// Decode 解码数据
func (c *Client) Decode() (*Package, bool) {
	var err error
	// 读取数据类型
	typeBuf, err := c.ReadBuffer.seek(0, TypeLen)
	if err != nil {
		return nil, false
	}

	// 读取数据长度
	lenBuf, err := c.ReadBuffer.seek(TypeLen, HeadLen)
	if err != nil {
		return nil, false
	}

	// 读取数据内容
	valueType := int(binary.BigEndian.Uint16(typeBuf))
	valueLen := int(binary.BigEndian.Uint16(lenBuf))

	valueBuf, err := c.ReadBuffer.read(HeadLen, valueLen)
	if err != nil {
		return nil, false
	}
	message := Package{Code: valueType, Content: valueBuf}
	return &message, true
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Client) Read() (int, error) {
	return c.ReadBuffer.readFromReader()
}

// HandleConnect 建立连接
func (c *Client) HandleConnect() {
	log.Print("tcp connect")
}

// HandleReadErr 读取conn错误
func (c *Client) HandleReadErr(err error) {
	log.Println("连接读取异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}
	c.Release()

	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
	log.Print("连接读取未知异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
}

// Release 释放TCP连接
func (c *Client) Release() {
	delete(c.DeviceId)
	err := c.Codec.Conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	/*
	   publishOffLine(transfer.OffLine{
	       DeviceId: c.DeviceId,
	       UserId:   c.UserId,
	   })*/
}

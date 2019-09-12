package gateway

import (
	"context"
	"io"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro"
	libmq "github.com/myproject-0722/my-micro/lib/mq"
	"github.com/myproject-0722/my-micro/lib/register"
	mq "github.com/myproject-0722/my-micro/proto/mq"
	packet "github.com/myproject-0722/my-micro/proto/packet"
	user "github.com/myproject-0722/my-micro/proto/user"
	"github.com/nsqio/go-nsq"
)

const (
	ReadDeadline  = 10 * time.Minute
	WriteDeadline = 10 * time.Second
)

type Client struct {
	Codec *Codec // 编解码器
	//ReadBuffer buffer
	//WriteBuf []byte
	UserId   int64
	DeviceId int64
	IsSignIn bool // 是否登录
}

func NewClient(conn net.Conn) *Client {
	codec := NewCodec(conn)
	return &Client{
		Codec: codec,
		//ReadBuffer: newBuffer(conn, packet.BufLen),
		//WriteBuf: make([]byte, packet.BufLen),
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

		//_, err = c.Read()
		_, err = c.Codec.Read()
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		for {
			message, ok := c.Decode()
			if ok {
				c.HandlePackage(message)
				continue
			} /*else {
				c.HandleDecodeErr()
				return
			}*/
			break
		}
	}
}

// HandlePackage 处理消息包
func (c *Client) HandlePackage(pack *packet.Package) {
	// 未登录拦截
	if pack.CodeType != packet.CodeSignIn && c.IsSignIn == false {
		c.Release()
		return
	}

	switch pack.CodeType {
	case packet.CodeSignIn:
		c.HandlePackageSignIn(pack)
		break
	case packet.CodeHeadbeat:
		c.HandlePackageHeadbeat()
		break
	default:
		c.HandlePackageOther(pack)
		break
	}

	return
}

// HandlePackageHeadbeat 处理心跳包
func (c *Client) HandlePackageHeadbeat() {
	err := c.Codec.Eecode(packet.Package{CodeType: packet.CodeHeadbeatACK, Content: []byte{}}, WriteDeadline)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("心跳：", "device_id:", c.DeviceId, "user_id:", c.UserId)
}

func (c *Client) HandlePackageOther(pack *packet.Package) {
	message := mq.MQMessage{
		UserId:    c.UserId,
		DeviceId:  c.DeviceId,
		CodeType:  pack.CodeType,
		PbMessage: pack.Content,
	}

	libmq.PublishMessage("gateway", message)
}

func (c *Client) HandleMQ2ClientMessage(msg *nsq.Message) error {
	var message mq.MQMessage
	err := proto.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Error(err)
	}

	err = c.Codec.Eecode(packet.Package{CodeType: message.CodeType, Content: message.PbMessage}, WriteDeadline)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// HandlePackageSignIn 处理登录消息包
func (c *Client) HandlePackageSignIn(pack *packet.Package) {
	var sign user.SignInRequest
	err := proto.Unmarshal(pack.Content, &sign)
	if err != nil {
		log.Error(err)
		c.Release()
		return
	}

	log.Debug("Recv signin:", sign.DeviceId, " ", sign.UserId, " ", sign.Token)

	reg := register.NewRegistry()

	// create a new service
	service := micro.NewService(micro.Registry(reg))
	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := user.NewUserService("go.mymicro.srv.user", service.Client())

	// Make request
	rsp, err := cl.SignIn(context.Background(), &sign)
	if err != nil {
		log.Debug(err)
		return
	}

	content, err := proto.Marshal(rsp)
	if err != nil {
		log.Error(err)
		return
	}

	err = c.Codec.Eecode(packet.Package{CodeType: packet.CodeSignInACK, Content: content}, WriteDeadline)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug("SignIn rescode:", rsp.ResCode)
	if rsp.ResCode == 200 {
		c.IsSignIn = true
		c.DeviceId = sign.DeviceId
		c.UserId = sign.UserId
		store(c.DeviceId, c)

		//var clientMsgTopic = "clientmsg_" + strconv.FormatInt(c.DeviceId, 10)
		//libmq.NsqConsumer(clientMsgTopic, "1", HandleMQ2ClientMessage, 20)
	}
}

// Decode 解码数据
func (c *Client) Decode() (*packet.Package, bool) {
	return c.Codec.Decode()
	/*var err error
	// 读取数据类型
	typeBuf, err := c.ReadBuffer.seek(0, packet.TypeLen)
	if err != nil {
		return nil, false
	}

	// 读取数据长度
	lenBuf, err := c.ReadBuffer.seek(packet.TypeLen, packet.HeadLen)
	if err != nil {
		return nil, false
	}

	// 读取数据内容
	valueType := int32(binary.BigEndian.Uint32(typeBuf))
	valueLen := int(binary.BigEndian.Uint32(lenBuf))

	valueBuf, err := c.ReadBuffer.read(packet.HeadLen, valueLen)
	if err != nil {
		return nil, false
	}
	message := packet.Package{CodeType: valueType, Content: valueBuf}
	return &message, true*/
	/*var err error

	//读取magic
	magicNumber, err := c.ReadBuffer.read(0, packet.MagicLen)
	if err != nil {
		return nil, false
	}

	if bytes.Compare(magicNumber, packet.MagicNumber) != 0 {
		log.Error("magicNumber error: ", magicNumber, packet.MagicNumber)
		return nil, false
	}

	// 读取数据类型
	typeBuf, err := c.ReadBuffer.seek(packet.MagicLen, packet.MagicLen+packet.TypeLen)
	if err != nil {
		return nil, false
	}

	//读取checksum
	checkSum, err := c.ReadBuffer.read(packet.MagicLen+packet.TypeLen, packet.MagicLen+packet.TypeLen+packet.CheckSumLen)
	if err != nil {
		return nil, false
	}

	// 读取数据长度
	lenBuf, err := c.ReadBuffer.seek(packet.MagicLen+packet.TypeLen+packet.CheckSumLen, packet.HeadLen)
	if err != nil {
		return nil, false
	}

	// 读取数据内容
	valueType := int32(binary.BigEndian.Uint32(typeBuf))
	valueLen := int(binary.BigEndian.Uint32(lenBuf))

	valueBuf, err := c.ReadBuffer.read(packet.HeadLen, valueLen)
	if err != nil {
		return nil, false
	}

	sum := sha256.Sum256(valueBuf)
	if sum[0] != checkSum[0] || sum[1] != checkSum[1] || sum[2] != checkSum[2] || sum[3] != checkSum[3] {
		return nil, false
	}

	message := packet.Package{CodeType: valueType, Content: valueBuf}
	return &message, true*/
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Client) Read() (int, error) {
	return c.Codec.ReadBuf.readFromReader()
}

// HandleConnect 建立连接
func (c *Client) HandleConnect() {
	log.Debug("HandleConnect")
}

// HandleReadErr 读取conn错误
func (c *Client) HandleReadErr(err error) {
	log.Debug("连接读取异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
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
	log.Debug("连接读取未知异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
}

// HandleReadErr 读取conn错误
func (c *Client) HandleDecodeErr() {
	log.Debug("连接读取Decode异常：", "device_id", c.DeviceId, "user_id", c.UserId)
	c.Release()
}

// Release 释放TCP连接
func (c *Client) Release() {
	delete(c.DeviceId)
	err := c.Codec.Conn.Close()
	if err != nil {
		log.Error(err)
	}
}

package client

import (
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/myproject-0722/my-micro/gateway"
	user "github.com/myproject-0722/my-micro/proto/user"
)

type TcpClient struct {
	DeviceId     string
	UserId       int64
	Token        string
	SendSequence int64
	SyncSequence int64
	codec        *gateway.Codec
}

func (c *TcpClient) Start() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
		return
	}

	codec := gateway.NewCodec(conn)
	c.codec = codec
}

func (c *TcpClient) SignIn() {
	signIn := user.SignInRequest{
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
		Token:    c.Token,
	}

    log.Println(signIn.DeviceId, signIn.UserId, signIn.Token)
	signInBytes, err := proto.Marshal(&signIn)
	if err != nil {
		log.Fatal(err)
		return
	}

	pack := gateway.Package{Code: gateway.CodeSignIn, Content: signInBytes}
	c.codec.Eecode(pack, 10*time.Second)
}

func (c *TcpClient) SyncTrigger() {
	/*
	   bytes, err := proto.Marshal(&pb.SyncTrigger{SyncSequence: c.SyncSequence})
	   if err != nil {
	       fmt.Println(err)
	       return
	   }
	   err = c.codec.Eecode(connect.Package{Code: connect.CodeSyncTrigger, Content: bytes}, 10*time.Second)
	   if err != nil {
	       fmt.Println(err)
	   }*/
}

func (c *TcpClient) HeadBeat() {
	ticker := time.NewTicker(time.Second * 1)
	for _ = range ticker.C {
		err := c.codec.Eecode(gateway.Package{Code: gateway.CodeHeadbeat, Content: []byte{}}, 10*time.Second)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *TcpClient) Receive() {
	for {
		_, err := c.codec.Read()
		if err != nil {
			log.Fatal(err)
			return
		}

		for {
			pack, ok := c.codec.Decode()
			if ok {
				c.HandlePackage(*pack)
				continue
			}
			break
		}
	}
}

func (c *TcpClient) HandlePackage(pack gateway.Package) error {
	switch pack.Code {
	/*case gateway.CodeSignInACK:
	  ack := pb.SignInACK{}
	  err := proto.Unmarshal(pack.Content, &ack)
	  if err != nil {
	      log.Fatal(err)
	      return err
	  }
	  log.Println("设备登录回执：%#v\n", ack)*/
	case gateway.CodeHeadbeatACK:
		log.Println("心跳回执")
		/*case connect.CodeMessageSendACK:
		        ack := pb.MessageSendACK{}
		        err := proto.Unmarshal(pack.Content, &ack)
		        if err != nil {
		            logger.Sugar.Error(err)
		            return err
		        }
				logger.Sugar.Info("消息发送回执：%#v\n", ack)*/
	/*case connect.CodeMessage:
	        message := pb.Message{}
	        err := proto.Unmarshal(pack.Content, &message)
	        if err != nil {
	            logger.Sugar.Error(err)
	            return err
	        }

	        for _, v := range message.Messages {
	            logger.Sugar.Info(message)
	        }

	        if len(message.Messages) == 0 {
	            return nil
	        }

	        ack := pb.MessageACK{SyncSequence: message.Messages[len(message.Messages)-1].SyncSequence}
	        ackBytes, err := proto.Marshal(&ack)
	        if err != nil {
	            logger.Sugar.Error(err)
	            return err
	        }

	        c.SyncSequence = ack.SyncSequence

	        err = c.codec.Eecode(connect.Package{Code: connect.CodeMessageACK, Content: ackBytes}, 10*time.Second)
	        if err != nil {
	            fmt.Println(err)
	            return err
			}*/
	default:
		log.Println("switch other")
	}
	return nil
}

func (c *TcpClient) SendMessage() {
	/*
			send := pb.MessageSend{}
		    fmt.Println("input ReceiverType,ReceiverId,Content")
		    fmt.Scanf("%d %d %s", &send.ReceiverType, &send.ReceiverId, &send.Content)
		    send.Type = 1
		    c.SendSequence++
		    send.SendSequence = c.SendSequence
		    bytes, err := proto.Marshal(&send)
		    if err != nil {
		        fmt.Println(err)
		        return
		    }
		    err = c.codec.Eecode(connect.Package{Code: connect.CodeMessageSend, Content: bytes}, 10*time.Second)
		    if err != nil {
		        fmt.Println(err)
		    }*/
}
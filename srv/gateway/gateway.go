package main

import (
	"github.com/gogo/protobuf/proto"
	"github.com/myproject-0722/my-micro/conf"
	"github.com/myproject-0722/my-micro/gateway"
	liblog "github.com/myproject-0722/my-micro/lib/log"
	libmq "github.com/myproject-0722/my-micro/lib/mq"
	mq "github.com/myproject-0722/my-micro/proto/mq"
	"github.com/myproject-0722/my-micro/proto/packet"
	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

func HandleMQ2ClientMessage(msg *nsq.Message) error {
	var message mq.MQMessage
	err := proto.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Error(err)
	}

	c := gateway.Load(message.DeviceId)
	if c != nil {
		err = c.Codec.Eecode(packet.Package{CodeType: message.CodeType, Content: message.PbMessage}, gateway.WriteDeadline)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

/*
type srv struct {
	*mux.Router
}*/

func main() {
	log.Debug("gateway started!")
	liblog.InitLog("/var/log/my-micro/", "gateway.log")
	gateway.InitRedis()
	libmq.Init()

	var clientMsgTopic = "clientmsg"
	go libmq.NsqConsumer(clientMsgTopic, "1", HandleMQ2ClientMessage, 20)

	//长连接服务
	conf := gateway.Conf{
		Address:    conf.GatewayListenAddress,
		MaxConnNum: conf.GatewayMaxConn,
		AcceptNum:  conf.AcceptNum,
	}
	server := gateway.NewTcpServer(conf)
	server.Start()
	/*
		var h http.Handler
		r := mux.NewRouter()
		s := &srv{r}
		h = s*/
	//init tcpserver

	//init webserver
}

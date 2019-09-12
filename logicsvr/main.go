package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	conf "github.com/myproject-0722/my-micro/conf"
	liblog "github.com/myproject-0722/my-micro/lib/log"
	libmq "github.com/myproject-0722/my-micro/lib/mq"
	clientmsg "github.com/myproject-0722/my-micro/proto/message"
	mq "github.com/myproject-0722/my-micro/proto/mq"
	packet "github.com/myproject-0722/my-micro/proto/packet"
	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr: conf.RedisIP,
			DB:   2,
		},
	)

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Error("redis err ")
		panic(err)
	}
}

type LogicHandler interface {
	HandleMessage(message *nsq.Message) error
}

type logicFunc func(*mq.MQMessage) error

var funcMap = map[int32]logicFunc{}

var op func(*nsq.Message) error

func HandleSingleMessage(message *mq.MQMessage) error {
	//log.Debug("im start!")

	var SingleMessage clientmsg.SingleMessage
	err2 := proto.Unmarshal(message.PbMessage, &SingleMessage)
	if err2 != nil {
		log.Fatal("msgqueue2 unmarshal err")
	}
	/*
		{
			//ack
			var topic = "clientmsg_" + strconv.FormatInt(message.DeviceId, 10)
			SingleMessageAck := clientmsg.SingleMessageAck{
				Seq: SingleMessage.Seq,
			}

			ackdata, err := proto.Marshal(&SingleMessageAck)
			if err != nil {
				log.Fatal("proto ack marshal err")
			}

			ackmq := mq.MQMessage{
				UserId:    message.UserId,
				DeviceId:  message.DeviceId,
				CodeType:  packet.CodeMessageSendACK,
				PbMessage: ackdata,
			}
			libmq.PublishMessage(topic, ackmq)
		}*/

	var to int64 = SingleMessage.To
	var key string = "userdevice:" + strconv.FormatInt(to, 10)
	sMembers := redisClient.SMembers(key)
	for _, v := range sMembers.Val() {
		deviceid := v
		log.Debug(deviceid)
		//var topic = "clientmsg_" + deviceid
		var topic = "clientmsg"
		libmq.PublishMessage(topic, *message)
	}
	return nil
}

func RegisterLogicFunc(CodeType int32, handle func(message *mq.MQMessage) error) {
	funcMap[CodeType] = handle
	//op = handle
}

func HandleLogic(CodeType int32, message *mq.MQMessage) {
	//op(message)
	if v, ok := funcMap[CodeType]; ok {
		v(message)
	} else {

	}
}

func HandleGatewayMessage(msg *nsq.Message) error {
	var message mq.MQMessage
	err := proto.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Fatal("msgqueue2 unmarshal err")
	}
	fmt.Println(message.DeviceId, message.UserId)

	HandleLogic(message.CodeType, &message)
	return nil
}

func StartNsqConsumer() {
	libmq.NsqConsumer("gateway", "1", HandleGatewayMessage, 20)
}

func StartLogic() {
	RegisterLogicFunc(packet.CodeMessage, HandleSingleMessage)
}

func main() {
	libmq.Init()

	liblog.InitLog("/var/log/my-micro/", "logicsvr.log")

	InitRedis()

	StartLogic()

	StartNsqConsumer()

	log.Print("logic server start!")
	for {
		time.Sleep(time.Second * 1)
	}
}

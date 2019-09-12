package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	mq "github.com/myproject-0722/my-micro/proto/mq"
	user "github.com/myproject-0722/my-micro/proto/user"
)

func main() {
	//fmt.Print("test")

	user1 := user.SignInRequest{
		DeviceId: 123456,
		UserId:   1,
		Token:    "123456789",
	}

	userdata, err := proto.Marshal(&user1)
	if err != nil {
		log.Fatal("proto user marshal err")
	}

	msgqueue := mq.MQMessage{
		UserId:    1,
		PbMessage: userdata,
	}

	mqdata, err := proto.Marshal(&msgqueue)
	if err != nil {
		log.Fatal("proto mq marshal err")
	}

	var msgqueue2 mq.MQMessage
	err = proto.Unmarshal(mqdata, &msgqueue2)
	if err != nil {
		log.Fatal("msgqueue2 unmarshal err")
	}

	var user2 user.SignInRequest
	err = proto.Unmarshal(msgqueue2.PbMessage, &user2)
	if err != nil {
		log.Fatal("user2 unmarshal err")
	}

	fmt.Print(user2.DeviceId, ",", user2.UserId, ",", user2.Token)
}

package main

import (
	"time"

	"github.com/myproject-0722/my-micro/conf"
	liblog "github.com/myproject-0722/my-micro/lib/log"
	"github.com/myproject-0722/my-micro/robot/client"
	//log "github.com/sirupsen/logrus"
)

func main() {
	liblog.InitLog("/var/log/my-micro/", "robot.log")

	var clientArray [conf.RobotNum]client.TcpClient

	for i := 0; i < conf.RobotNum; i++ {
		userId := int64(i + 1)
		clientArray[i].UserId = userId
		clientArray[i].DeviceId = userId
		clientArray[i].Token = "123456"
		clientArray[i].SendSequence = 1
		clientArray[i].SyncSequence = 1

		clientArray[i].Start()
		clientArray[i].SignIn()

		time.Sleep(100 * time.Millisecond)
		/*
			client[i] = client.TcpClient{
				UserId:       userId,
				DeviceId:     123456,
				Token:        "999999",
				SendSequence: 1,
				SyncSequence: 1,
			}*/
	}

	//client.Start()
	//client.SignIn()
	for {
		for j := 0; j < conf.RobotNum; j++ {
			clientArray[j].SendMessage()
		}
		//client.SendMessage()
		time.Sleep(100 * time.Millisecond)
	}
}

package main

import (
	"time"

	liblog "github.com/myproject-0722/my-micro/lib/log"
	"github.com/myproject-0722/my-micro/robot/client"
	//log "github.com/sirupsen/logrus"
)

func main() {
	liblog.InitLog("/var/log/my-micro/", "robot.log")

	client := client.TcpClient{
		UserId:       1,
		DeviceId:     123456,
		Token:        "999999",
		SendSequence: 1,
		SyncSequence: 1,
	}
	client.Start()
	client.SignIn()
	for {
		client.SendMessage()
		time.Sleep(10 * time.Millisecond)
	}
}

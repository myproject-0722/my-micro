package main

import (
	"fmt"

	"github.com/myproject-0722/my-micro/robot/client"
)

func main() {
	fmt.Print("rebot start!")

	client := client.TcpClient{}
    client.UserId = 1;
    client.DeviceId = "12234";
    client.Token = "9887743";
    client.SendSequence = 1;
    client.SyncSequence = 1;
	fmt.Println("input UserId,DeviceId,Token,SendSequence,SyncSequence")
	//fmt.Scanf("%d %d %s %d %d", &client.UserId, &client.DeviceId, &client.Token, &client.SendSequence, &client.SyncSequence)
	client.Start()
	client.SignIn()
	for {
		client.SendMessage()
	}
}

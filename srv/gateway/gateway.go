package main

import (
	"github.com/myproject-0722/my-micro/conf"
	"github.com/myproject-0722/my-micro/gateway"
	liblog "github.com/myproject-0722/my-micro/lib/log"
	libmq "github.com/myproject-0722/my-micro/lib/mq"
	log "github.com/sirupsen/logrus"
)

/*
type srv struct {
	*mux.Router
}*/

func main() {
	log.Debug("gateway started!")
	liblog.InitLog("/var/log/my-micro/", "gateway.log")
	gateway.InitRedis()
	libmq.Init()

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

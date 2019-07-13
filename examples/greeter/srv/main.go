package main

import (
	"log"
	"time"

    "github.com/micro/go-micro"
    "github.com/micro/go-micro/registry"
    "github.com/micro/go-micro/registry/consul"
	hello "github.com/micro/examples/greeter/srv/proto/hello"

	"context"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
    // 推荐使用etcd集群 做为服务发现,为测试暂用consul
    reg := consul.NewRegistry(func(op *registry.Options) {
        op.Addrs = []string{
            "127.0.0.1:8500",
        }
    })
    /* reg := etcdv3.NewRegistry(func(op *registry.Options) {
        op.Addrs = []string{
        "http://192.168.3.34:2379", "http://192.168.3.18:2379", "http://192.168.3.110:2379",
       }
    })*/
    service := micro.NewService(
        micro.Registry(reg),
        micro.Name("go.micro.srv.greeter"),
	    micro.RegisterTTL(time.Second*30),
	    micro.RegisterInterval(time.Second*10),
    )

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

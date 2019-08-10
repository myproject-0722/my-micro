package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	hello "github.com/myproject-0722/my-micro/examples/greeter/srv/proto/hello"
)

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

	// create a new service
	service := micro.NewService(micro.Registry(reg))

	// parse command line flags
	service.Init()

	//client := micro.NewClient()
	c := service.Client()
	req := c.NewRequest("go.mymicro.srv.greeter", "Say.Hello", &hello.Request{
		Name: "John",
	})

	rsp := &hello.Response{}

	if err := c.Call(context.Background(), req, rsp); err != nil {
		fmt.Println("call err: ", err, rsp)
		return
	}

	fmt.Println("rsp:", rsp.Msg)

	// Use the generated client stub
	/*cl := hello.NewSayService("go.mymicro.srv.greeter", service.Client())

	// Make request
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)

	time.Sleep(time.Duration(300) * time.Second)

	// Make request
	rsp1, err1 := cl.Hello(context.Background(), &hello.Request{
		Name: "lixu",
	})
	if err != nil {
		fmt.Println(err1)
		return
	}

	fmt.Println(rsp1.Msg)*/
}

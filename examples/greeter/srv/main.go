package main

import (
	"log"
	"time"

	"github.com/myproject-0722/my-micro/lib/tracer"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	hello "github.com/myproject-0722/my-micro/examples/greeter/srv/proto/hello"

	"context"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	t, io, err := tracer.NewTracer("go.micro.srv.greeter", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

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
	//TraceBoot()

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.mymicro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
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

package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/micro/go-micro"
	api "github.com/micro/go-micro/api/proto"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	hello "github.com/myproject-0722/my-micro/examples/greeter/srv/proto/hello"

	"context"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/myproject-0722/my-micro/lib/tracer"
	opentracing "github.com/opentracing/opentracing-go"
)

type Say struct {
	Client hello.SayService
}

func (s *Say) Hello(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received Say.Hello API request")

	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.api.greeter", "Name cannot be blank")
	}

	response, err := s.Client.Hello(ctx, &hello.Request{
		Name: strings.Join(name.Values, " "),
	})
	if err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})
	rsp.Body = string(b)

	return nil
}

func main() {
	t, io, err := tracer.NewTracer("go.mymicro.api.greeter", "localhost:6831")
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

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.mymicro.api.greeter"),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// parse command line flags
	service.Init()

	service.Server().Handle(
		service.Server().NewHandler(
			&Say{Client: hello.NewSayService("go.mymicro.srv.greeter", service.Client())},
		),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

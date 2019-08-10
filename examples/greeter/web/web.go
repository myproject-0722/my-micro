package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/web"
	hello "github.com/myproject-0722/my-micro/examples/greeter/srv/proto/hello"

	"context"
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
	webservice := micro.NewService(micro.Registry(reg))

	service := web.NewService(
		web.Name("go.mymicro.web.greeter"),
	)

	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			name := r.Form.Get("name")
			if len(name) == 0 {
				name = "World"
			}

			cl := hello.NewSayService("go.micro.srv.greeter", webservice.Client())
			rsp, err := cl.Hello(context.Background(), &hello.Request{
				Name: name,
			})

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))
			return
		}

		fmt.Fprint(w, `<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
	})

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

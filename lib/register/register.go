package register

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

var registerService micro.Service
var reg registry.Registry

func NewRegistry() registry.Registry {
	if reg != nil {
		log.Print("reg is not null")
		return reg
	}

	log.Print("reg is null")
	// 推荐使用etcd集群 做为服务发现,为测试暂用consul
	reg = consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	return reg
	/* reg := etcdv3.NewRegistry(func(op *registry.Options) {
	    op.Addrs = []string{
	    "http://192.168.3.34:2379", "http://192.168.3.18:2379", "http://192.168.3.110:2379",
	   }
	})*/
}

func NewRegistryService() micro.Service {
	if registerService != nil {
		log.Print("service is not null")
		return registerService
	}

	log.Print("service is null")
	reg := NewRegistry()

	// create a new service
	registerService = micro.NewService(micro.Registry(reg))
	// parse command line flags
	registerService.Init()

	//registerService = service

	if registerService == nil {
		log.Print("service is null 2")
	}
	return registerService
}

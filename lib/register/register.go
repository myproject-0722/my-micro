package register

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

func NewRegistry() registry.Registry {
	// 推荐使用etcd集群 做为服务发现,为测试暂用consul
	reg := consul.NewRegistry(func(op *registry.Options) {
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

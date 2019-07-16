module github.com/myproject-0722/my-micro

go 1.12

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

require (
	github.com/99designs/gqlgen v0.7.1
	github.com/emicklei/go-restful v2.8.1+incompatible
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/micro/examples v0.1.0
	github.com/micro/go-micro v1.7.0
	github.com/micro/go-plugins v0.22.0
	github.com/micro/micro v0.22.0
	github.com/nats-io/nats-server/v2 v2.0.2 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/vektah/gqlparser v1.1.0
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	google.golang.org/grpc v1.21.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

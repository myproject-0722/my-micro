# my-micro
my-micro基于[micro](https://github.com/micro/micro)，编写微服务框架所使用到场景的例子。

# 1.环境安装：
## 1.protobuf
[下载](https://github.com/protocolbuffers/protobuf/archive/v3.6.0.1.zip)或wget https://github.com/protocolbuffers/protobuf/archive/v3.6.0.1.zip
./autogen.sh && ./configure && make && make check  
sudo make install    
sudo ldconfig  
## 2.protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go
## 3.protoc-gen-micro
go get github.com/micro/protoc-gen-micro
## 4.Consul安装运行
consul agent -dev

# 2.范例说明
## 1. 运行micro api
micro --registry=consul --registry_address=127.0.0.1:8500 --server_advertise=192.168.0.194:8080 api --handler=api --address=0.0.0.0:8080 --namespace=go.mymicro.api
## 2.启动srv及api
cd examples/greeter
go run srv/main.go
go run api/api.go
## 3.测试
curl http://localhost:8080/greeter/say/hello?name=John
## 4.如使用web
go run web/web.go
micro --registry=consul --registry_address=127.0.0.1:8500 web
打开网址:http://localhost:8082/greeter
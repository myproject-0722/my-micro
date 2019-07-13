# gateway
作为客户端与服务端的网关，将后台服务与外网隔离
# 启动
micro --registry=consul --registry_address=127.0.0.1:8500 --server_advertise=192.168.0.180:8080 api --handler=api --address=0.0.0.0:8080
8080可修改为其它地址
发现服务为方便测试，这里暂用consul
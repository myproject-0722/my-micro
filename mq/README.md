# NSQ

## 拉取NSQ镜像
docker pull nsqio/nsq

## 启动NSQ服务
docker run -d --name lookupd -p 4160:4160 -p 4161:4161 nsqio/nsq:latest /nsqlookupd

## 启动nsqadmin管理系统
docker run -d --name nsqadmin -p 4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=192.168.0.140:4161

## 开启一个nsqd节点服务
docker run -d --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq:latest /nsqd --broadcast-address=192.168.0.140 --lookupd-tcp-address=192.168.0.140:4160
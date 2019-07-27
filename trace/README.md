# 链路追踪
这里暂以Jaeger作为链路追踪系统(zipkin也可)
## 运行
启动es
启动jaeger
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp \
  -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest
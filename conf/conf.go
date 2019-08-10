package conf

//consul
var (
	ConsulAddresses string = "127.0.0.1:8500" //consul addresses(such as 127.0.0.1:8500;127.0.0.1:8600)
)

//gateway
var (
	GatewayListenAddress string = "0.0.0.0:9999"   //listen ip and port
	GatewayMaxConn       int    = 10000            //max conn num
	AcceptNum            int    = 5                //
	NSQIP                string = "127.0.0.1:4150" //NSQ
	RedisIP              string = "127.0.0.1:6379" //redis
)

#dashd -testnet -daemon
#curl -X PUT -d '' 'http://localhost:8500/v1/kv/micro/config/jwt-key'
#./micro --registry=consul --registry_address=127.0.0.1:8500 --server_advertise=127.0.0.1:8000 api --handler=api --address=0.0.0.0:8000 --namespace=go.mnhosted.api &
#./micro --registry=consul --registry_address=161.189.42.122:8500 --server_advertise=127.0.0.1:8000 api --handler=api --address=0.0.0.0:8000 --namespace=go.mnhosted.api &
./user &
./logicsvr &
./gateway &

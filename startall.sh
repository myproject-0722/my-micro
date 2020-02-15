docker-compose up -d
sleep 5
docker-compose restart consul
docker-compose restart elasticsearch
sleep 2
docker-compose restart nsqlookupd
docker-compose restart nsqd
docker-compose restart nsqadmin
docker-compose restart collector
docker-compose restart agent
docker-compose restart query
docker-compose restart hotrod
killall user
killall logicsrv
killall gateway
./bin/user &
./bin/logicsrv &
./bin/gateway &

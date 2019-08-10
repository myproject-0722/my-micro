protoc --proto_path=./user --micro_out=./user --go_out=./user ./user/*.proto
protoc --proto_path=./mq --micro_out=./mq --go_out=./mq ./mq/*.proto
protoc --proto_path=./message --micro_out=./message --go_out=./message ./message/*.proto
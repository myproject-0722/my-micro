go build -race ../rpcsrv/user/
go build -race ../logicsrv/
go build -race ../srv/gateway/
#go build -race -gcflags "-N -l" ../notifysrv/alipay/alipay-notify.go 

package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/micro/go-micro"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/myproject-0722/my-micro/lib/register"
	"github.com/myproject-0722/my-micro/lib/tracer"
	user "github.com/myproject-0722/my-micro/proto/user"
)

type User struct{}

func (s *User) SignIn(ctx context.Context, req *user.SignInRequest, rsp *user.SignInResponse) error {
	log.Print("Received SignInRequest DeviceId: ", req.DeviceId, " UserId: ", req.UserId, " Token: ", req.Token)
	rsp.ResMsg = strconv.FormatInt(req.UserId, 10) + " SignIn OK!"
	return nil
}

func main() {
	tracer.InitTracer("go.micro.srv.user")

	reg := register.NewRegistry()

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.mymicro.srv.user"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	user.RegisterUserHandler(service.Server(), new(User))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

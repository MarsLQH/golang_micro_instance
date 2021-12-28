package main

import (
	"fmt"

	"MarsLuo/config"
	"MarsLuo/services/handler"
	pb "MarsLuo/services/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// 初始化配置
	config.Init()

	srv := service.New(
		service.Name("srv.micro.instance"),
		service.BeforeStop(func() error {
			//todo
			return nil
		}),
		service.AfterStop(func() error {
			//todo
			return nil
		}),
	)

	pb.RegisterbookHandler(srv.Server(), new(handler.Videos))
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	fmt.Println("the server is under service...")
}

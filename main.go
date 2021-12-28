package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"

	"git.myarena7.com/arena/hicourt/config"
	"git.myarena7.com/arena/hicourt/services/handler"
	pb "git.myarena7.com/arena/hicourt/services/proto"
	vSync "git.myarena7.com/arena/hicourt/sync"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// 初始化配置
	config.Init()
	fmt.Println("hello Arena Hicourt")
	//=======get video from other platform
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Err(err.(error)).Msgf("goroutine recovered with err=,%v", err)
			}
		}()
		vSync.GetVideo()
	}()

	srv := service.New(
		service.Name("hicourt"),
		service.BeforeStop(func() error {
			return nil
		}),
		service.AfterStop(func() error {
			return nil
		}),
	)

	pb.RegisterHicourtHandler(srv.Server(), new(handler.Videos))
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	fmt.Println("the server is under service...")
	wg.Wait()
	/*
		tracer, err := core.GetTracer("hicourt", "http://127.0.0.1:9411")
		if err != nil {
			log.Err(err).Msgf("unable to create local endpoint: %+v\n", tracer)
		}

	*/
}

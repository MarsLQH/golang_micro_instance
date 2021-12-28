package main

import (
	"github.com/micro/micro/plugin"
	//"github.com/micro/micro/plugin"
	//"github.com/urfave/cli/v2"
)

func init() {
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("etcd"),
		//plugin.WithFlag(&cli.StringFlag{Name:   "etcd",Usage:  "This is an example plugin flag",EnvVars: []string{"EXAMPLE_FLAG"},		Value: "avalue"		}),
		//plugin.WithInit(func(ctx *cli.Context) error {
		//	log.Println("Got value for example_flag", ctx.String("example_flag"))
		//	return nil
		//}),
	))
}

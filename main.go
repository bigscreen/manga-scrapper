package main

import (
	"os"

	"github.com/bigscreen/manga-scrapper/config"
	"github.com/bigscreen/manga-scrapper/logger"
	"github.com/bigscreen/manga-scrapper/server"
	"github.com/urfave/cli"
)

func main() {
	defer handleInitError()

	config.Load()
	logger.SetupLogger()

	clientApp := cli.NewApp()
	clientApp.Name = "mangajack"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "server",
			Description: "Start HTTP api server",
			Action: func(c *cli.Context) error {
				server.StartAPIServer()
				return nil
			},
		},
	}

	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func handleInitError() {
	if e := recover(); e != nil {
		logger.FatalF("Failed to load the app due to error : %s", e)
	}
}

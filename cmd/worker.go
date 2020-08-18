package cmd

import (
	"log"
	"os"
	"os/signal"

	cli "gopkg.in/urfave/cli.v2"
)

func newWorkerCmd() *cli.Command {
	return &cli.Command{
		Name:   "worker",
		Usage:  "worker",
		Action: handleWorkerCmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   defaultConfigFileFolder,
				Usage:   "path of config file",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Usage: "log level",
			},
		},
	}
}

func handleWorkerCmd(c *cli.Context) error {
	// cfg, err := loadConfig(c)
	// if err != nil {
	// 	return err
	// }

	// logger := newLogger(c.String("log-level"))
	// wk := worker.NewWorker(cfg, message.NewPubSub(), logger)
	// if err := wk.Setup(); err != nil {
	// 	return err
	// }

	// if err := wk.Run(c.Context); err != nil {
	// 	return err
	// }

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Service ...")
	return nil
}

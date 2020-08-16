package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/koolay/sqlboss/pkg/agent"
	"github.com/koolay/sqlboss/pkg/conf"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/koolay/sqlboss/pkg/worker"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

func newAgentCmd() *cli.Command {
	serveCmd := &cli.Command{
		Name:  "agent",
		Usage: "sqlboss agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   defaultConfigFileFolder,
				Usage:   "path of config file",
			},
		},
		Action: serveAction,
	}

	return serveCmd
}

func serveAction(c *cli.Context) error {
	cfg, err := loadConfig(c)
	if err != nil {
		return err
	}

	logger := newLogger(c.String("log-level"))

	log.Println(cfg)

	log.Println("start worker")

	go startWorker(context.Background(), cfg, logger)

	go func() {
		pub, err := agent.NewPublisher(cfg, message.NewPubSub())
		if err != nil {
			log.Fatal(err)
		}

		count := 0
		for {
			log.Println("publish message")
			count++
			if perr := pub.Publish(cfg.Stream.Topic, []byte(fmt.Sprintf("hello, %d", count))); perr != nil {
				logger.WithError(perr).Error("failed to pushlish")
			}

			time.Sleep(1 * time.Second)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Service ...")
	return nil
}

func startWorker(ctx context.Context, cfg *conf.Config, logger *logrus.Logger) {
	wk := worker.NewWorker(cfg, message.NewPubSub(), logger)
	if err := wk.Setup(); err != nil {
		log.Fatal(err)
	}

	if err := wk.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

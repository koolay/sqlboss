package cmd

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/koolay/sqlboss/pkg/agent"
	"github.com/koolay/sqlboss/pkg/lineage"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/koolay/sqlboss/pkg/proto"
	"github.com/koolay/sqlboss/pkg/store"
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
			&cli.StringFlag{
				Name:  "log-level",
				Value: "INFO",
				Usage: "log level",
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
	ctx := c.Context

	cqrsMarshaler := cqrs.JSONMarshaler{}
	cqrsServer, err := message.NewCQRSServer(logger.WithContext(ctx), cqrsMarshaler)
	if err != nil {
		return err
	}

	if err = cqrsServer.Setup([]message.CommandHandlerGenerator{
		func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.CommandHandler {
			return agent.NewSQLCommandHandler(eb)
		},
	},

		[]message.EventHandlerGenerator{
			func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.EventHandler {
				return store.StoreOnSQLEventHandler{}
			},
			func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.EventHandler {
				return lineage.LineageOnSQLEventHandler{}
			},
		}); err != nil {
		return err
	}

	commandBus := cqrsServer.GetCommandBus()

	go func() {
		if serr := cqrsServer.Start(); serr != nil {
			logger.Fatal(serr)
		}
	}()

	go func() {
		count := 0
		time.Sleep(1 * time.Second)
		for {
			log.Println("publish message")
			count++

			data := &proto.SqlCommand{}
			if perr := commandBus.Send(ctx, data); perr != nil {
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

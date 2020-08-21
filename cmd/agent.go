package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/koolay/sqlboss/pkg/lineage"
	"github.com/koolay/sqlboss/pkg/logging"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/koolay/sqlboss/pkg/obs"
	"github.com/koolay/sqlboss/pkg/proxy"
	"github.com/koolay/sqlboss/pkg/store"
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	mysqlVersion    = "5.6"
	defaultAddr     = "127.0.0.1:3309"
	defaultUser     = "root"
	defaultPassword = "123"
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

func initCQRSServer(ctx context.Context) (*message.CQRSServer, error) {
	logger := logging.LoggerFromContext(ctx)
	obser := obs.TelemetryFromContext(ctx)

	cqrsMarshaler := cqrs.JSONMarshaler{}
	cqrsServer, err := message.NewCQRSServer(logger, cqrsMarshaler)
	if err != nil {
		return nil, err
	}

	commandHandlerGenerators := []message.CommandHandlerGenerator{
		func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.CommandHandler {
			return store.NewStoreCommandHandler(logger)
		},

		func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.CommandHandler {
			return lineage.NewLineageCommandHandler(logger)
		},
	}
	eventHandlerGenerators := []message.EventHandlerGenerator{
		func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.EventHandler {
			return proxy.NewParseOnSQLEventHandler(logger, cb, obser.Meter)
		},
	}

	if err = cqrsServer.Setup(commandHandlerGenerators, eventHandlerGenerators); err != nil {
		return nil, err
	}

	return cqrsServer, nil
}

func serveAction(c *cli.Context) error {
	cfg, err := loadConfig(c)
	if err != nil {
		return err
	}

	logger := logging.NewLogger(c.String("log-level"))
	logEntry := logger.WithContext(context.Background())

	ctx := logging.WithLogger(c.Context, logEntry)

	obser, err := obs.NewTelemetry(obs.Config{
		ExporterPath: "/metrics",
		ExportPort:   2222,
	})

	ctx = obs.WithTelemetry(ctx, obser)

	if err != nil {
		return errors.Wrap(err, "failed to new telemetry")
	}

	cqrsServer, err := initCQRSServer(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to setup cqrs server")
	}

	eventBus := cqrsServer.GetEventBus()

	mysqlServer, err := proxy.NewProxy(cfg, logger, &proxy.MysqlServerConfig{
		Version:  mysqlVersion,
		Addr:     defaultAddr,
		User:     defaultUser,
		Password: defaultPassword,
		TargetConnection: proxy.Connection{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			Database: cfg.DB.Database,
		},
	}, eventBus)

	if err != nil {
		return errors.Wrap(err, "failed to new proxy")
	}

	go func() {
		if err := obser.Serve(); err != nil {
			log.Fatalf("failed to serve metrics, %+v", err)
		}
	}()

	go func() {
		log.Println("start mysql proxy server")
		if serr := mysqlServer.Start(); serr != nil {
			log.Fatal(serr)
		}
	}()

	go func() {
		if serr := cqrsServer.Start(); serr != nil {
			logger.Fatal(serr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Service ...")
	return nil
}

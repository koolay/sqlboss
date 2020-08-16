package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/koolay/sqlboss/pkg/conf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	defaultConfigFileFolder = "."
	defaultConfigFile       = "app.config"
)

func newLogger(levelName string) *logrus.Logger {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		level = logrus.ErrorLevel
	}

	logFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, fmt.Sprintf("%s:%v", filename, f.Line)
		},
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(logFormatter)
	logger.SetReportCaller(true)

	return logger
}

func loadConfig(c *cli.Context) (*conf.Config, error) {
	configFolder := c.String("config")
	if configFolder == "" {
		configFolder = defaultConfigFileFolder
	}

	configFolder, err := filepath.Abs(configFolder)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(configFolder); os.IsNotExist(err) {
		return nil, fmt.Errorf("configFolder '%s' not exist", configFolder)
	}

	configFilePath := filepath.Join(configFolder, defaultConfigFile)
	log.Printf("Load config file: %s\n", configFilePath)
	cfg, err := conf.LoadConfig(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	return cfg, nil
}

// NewApp new root command.
func NewApp() *cli.App {
	return &cli.App{
		Name:        "sqlboss",
		Usage:       "sqlboss",
		Description: "boss of sql",
		Version:     "1.0.0",
		Commands: []*cli.Command{
			newInitCmd(),
			newAgentCmd(),
			newWorkerCmd(),
		},
		Before: func(ctx *cli.Context) error {
			return nil
		},

		Action: func(c *cli.Context) error {
			return nil
		}}
}

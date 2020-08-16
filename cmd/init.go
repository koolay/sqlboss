package cmd

import (
	"path/filepath"

	"github.com/koolay/sqlboss/pkg/conf"
	cli "gopkg.in/urfave/cli.v2"
)

func newInitCmd() *cli.Command {
	initCmd := &cli.Command{
		Name:  "init",
		Usage: "init config file",
		Action: func(c *cli.Context) error {
			defaultConfigFile := filepath.Join(".", defaultConfigFile)
			return conf.InitDeafultCfgFile(defaultConfigFile)
		},
	}
	return initCmd
}

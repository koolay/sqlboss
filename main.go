package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/koolay/sqlboss/cmd"
)

var (
	version = "0.0.0"
	build   = "---"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	v := fmt.Sprintf("version=%s date=%s", version, build)

	rootCmd := cmd.NewApp()
	rootCmd.Version = v
	log.SetOutput(os.Stdout)

	if err := rootCmd.Run(os.Args); err != nil {
		panic(err)
	}
}

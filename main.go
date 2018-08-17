package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

const (
	DEFAULT_NODEURI = "http://localhost:8545"
)

func main() {

	// set up config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("$HOME/.wanutil")
	viper.AddConfigPath(".")

	viper.SetDefault("nodeuri", DEFAULT_NODEURI)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file:\n\n%s\n", err))
	}

	// set up cli
	app := cli.NewApp()

	app.Usage = "Wanchain CLI utility tool"
	app.UsageText = "wanutil <command> [options]"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{}
	app.Before = beforeApp
	app.Commands = commands

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func beforeApp(c *cli.Context) error {
	return nil
}

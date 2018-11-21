package main

import (
	"github.com/stevenxie/begone/cmd"
	gencmd "github.com/stevenxie/begone/strgen/cmd"
	"github.com/urfave/cli"
)

func main() {
	// Configure name, usage, metadata, etc.
	app := cli.NewApp()
	app.Name = cmd.Name
	app.Usage = cmd.Desc
	app.UsageText = cmd.Usage
	app.Version = cmd.Version
	app.HideVersion = true
	app.Flags = cmd.GlobalFlags

	// Configure app commands.
	app.Commands = []cli.Command{cmd.LoginCmd}
	app.Commands = append(app.Commands, gencmd.Commands...)

	// Run app.
	app.RunAndExitOnError()
}

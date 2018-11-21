package cmd

import (
	"github.com/stevenxie/begone/config"
	"github.com/stevenxie/begone/mbot"
	"github.com/urfave/cli"
)

// botConfig creates an mbot.Config using 'cfg' and cli flags from 'ctx'.
//
// Uses cli flags "cycles" and "delay".
func botConfig(ctx *cli.Context, cfg *config.Config) *mbot.Config {
	mcfg := mbot.NewConfig()
	mcfg.Username = cfg.Username
	mcfg.Password = cfg.Password

	if ctx.GlobalIsSet("cycles") {
		mcfg.Cycles = ctx.Int("cycles")
	}
	if ctx.GlobalIsSet("delay") {
		mcfg.Delay = ctx.Int("delay")
	}
	return mcfg
}

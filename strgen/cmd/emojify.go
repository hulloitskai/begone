package cmd

import (
	"fmt"
	"os"

	"github.com/stevenxie/begone/cmd"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
	"github.com/urfave/cli"
)

// EmojifyCmdName is the name of command emojify.
const EmojifyCmdName = "emojify"

var (
	// EmojifyCmd spams a target using an emoji generator.
	EmojifyCmd = cli.Command{
		Name:      EmojifyCmdName,
		Usage:     "spam a target using an emoji generator",
		UsageText: "emojify [-m <mode>] <converation ID>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "mode, m",
				Usage: "emoji generation method ('single', 'staircase', 'unique')",
				Value: "single",
			},
		},
		Action: emojify,
	}
)

func emojify(ctx *cli.Context) error {
	// Ensure arguments are valid.
	if len(ctx.Args()) > 1 {
		fmt.Fprint(os.Stderr, "Too many arguments.\n\n")
		cli.ShowCommandHelpAndExit(ctx, EmojifyCmdName, 1)
	}

	var (
		mode     = ctx.String("mode")
		gen, err = strgen.NewEmojifier(mode)
	)
	if err != nil {
		return ess.AddCtx("strgen/cmd: creating emojifier", err)
	}

	convoID := ctx.Args().First()
	return cmd.Begone(ctx, gen, convoID)
}

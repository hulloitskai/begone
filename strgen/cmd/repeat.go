package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/stevenxie/begone/cmd"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
	"github.com/urfave/cli"
)

const (
	// RepeatCmdName is the name of command repeat.
	RepeatCmdName = "repeat"

	usageText = `repeat <message> <conversation ID>

   Remember to surround your message with quotes if it has spaces!`
)

var (
	// RepeatCmd spams a target by.
	RepeatCmd = cli.Command{
		Name:      RepeatCmdName,
		Usage:     "repeatedly spam a target with a message",
		UsageText: usageText,
		Action:    repeat,
	}
)

func repeat(ctx *cli.Context) error {
	// Ensure arguments are valid.
	args := ctx.Args()
	if len(args) > 2 {
		fmt.Fprint(os.Stderr, "Too many arguments.\n\n")
		cli.ShowCommandHelpAndExit(ctx, RepeatCmdName, 1)
	}

	var msg, convoID = args.Get(0), args.Get(1)
	if msg == "" {
		var err error
		if msg, err = queryRepeatMsg(); err != nil {
			return ess.AddCtx("cmd: requesting repeat message", err)
		}
	}

	gen := strgen.NewRepeater(msg)
	return cmd.Begone(ctx, gen, convoID)
}

// queryRepeatMsg requests the message to repeat from stdin.
func queryRepeatMsg() (string, error) {
	var (
		msg string
		err error
	)

	for msg == "" {
		fmt.Println("Enter the message to repeat: ")
		r := bufio.NewReader(os.Stdin)
		if msg, err = r.ReadString('\n'); err != nil {
			return "", err
		}

		if msg == "" {
			fmt.Fprintln(os.Stderr, "Message must be non-empty!")
		}
	}

	return msg, nil
}

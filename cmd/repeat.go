package cmd

import (
	"bufio"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
)

func registerRepeatCmd(app *kingpin.Application) {
	repeatCmd = app.Command("repeat", "Spam a target with a message.")

	// Args:
	repeatOpts.Msg = repeatCmd.Arg("message", "The message to send repeatedly.").
		String()

	repeatOpts.ConvoID = repeatCmd.Arg(
		"conversation ID",
		"The target conversation ID (last portion of a www.messenger.com link).",
	).Default("").String()
}

var (
	repeatCmd *kingpin.CmdClause

	repeatOpts struct {
		Msg, ConvoID *string
	}
)

func repeat() error {
	// Parse arguments.
	p := interact.NewPrompter()

	// Validate arguments.
	if repeatOpts.Msg == nil {
		var err error
		if *repeatOpts.Msg, err = queryRepeatMessage(p); err != nil {
			return ess.AddCtx("querying for message", err)
		}
	}

	// Create generator.
	gen := strgen.NewRepeater(*repeatOpts.Msg)

	// Derive convoURL.
	convoURL, err := deriveConvoURL(*repeatOpts.ConvoID, p)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	var (
		bcfg   = deriveBotConfig()
		runner = deriveBotRunner(p)
	)
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleText(convoURL, gen)
	})
}

// queryRepeatMessage requests the message to repeat from os.Stdin.
func queryRepeatMessage(p *interact.Prompter) (string, error) {
	var (
		msg string
		err error
	)

	for msg == "" {
		p.Println("Enter the message to repeat:")
		r := bufio.NewReader(os.Stdin)
		if msg, err = r.ReadString('\n'); err != nil {
			return "", err
		}

		if msg == "" {
			p.Errln("Message must be non-empty!")
		}
	}

	return msg, nil
}

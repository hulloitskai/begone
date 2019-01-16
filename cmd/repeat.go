package cmd

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/pkg/mbot"
	"github.com/stevenxie/begone/pkg/strgen"
	ess "github.com/unixpickle/essentials"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func registerRepeatCmd(app *kingpin.Application) {
	repeatCmd = app.Command("repeat", "Spam a target with a message.")

	// Args:
	repeatCmd.Arg("message", "The message to send repeatedly.").
		StringVar(&repeatOpts.Msg)

	// Flags:
	repeatCmd.Flag(
		"stdin",
		"Read message from stdin (first arg becomes convo ID). Requires "+
			"you to be logged-in.",
	).BoolVar(&repeatOpts.Stdin)

	// Common options:
	registerCommonOpts(repeatCmd)
}

var (
	repeatCmd  *kingpin.CmdClause
	repeatOpts struct {
		Msg   string
		Stdin bool
	}
)

// stdinReadLimit is the maximum number of bytes to be read from os.Stdin.
const stdinReadLimit = 10000

func repeat() error {
	// Parse arguments.
	var (
		p      = interact.NewPrompter()
		bcfg   = deriveBotConfig()
		runner = deriveBotRunner(p)
		msg    string
	)

	// Validate arguments.
	if repeatOpts.Stdin { // configure runner first to validate login
		runner.Interactive = false // cannot ask for login details on os.Stdin
		if err := runner.Configure(bcfg); err != nil {
			return ess.AddCtx("configuring BotRunner", err)
		}

		data, err := ioutil.ReadAll(io.LimitReader(os.Stdin, stdinReadLimit))
		if err != nil {
			return ess.AddCtx("reading from os.Stdin", err)
		}

		msg = string(data)
		copts.ConvoID = repeatOpts.Msg // first arg is now the convo ID
	} else {
		msg = repeatOpts.Msg
		if msg == "" {
			var err error
			if msg, err = queryRepeatMessage(p); err != nil {
				return ess.AddCtx("querying for message", err)
			}
		}
	}

	// Create generator.
	gen := strgen.NewRepeater(msg)

	// Derive convoURL.
	convoURL, err := deriveConvoURL(copts.ConvoID, p)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	if !repeatOpts.Stdin {
		if err := runner.Configure(bcfg); err != nil {
			return ess.AddCtx("configuring BotRunner", err)
		}
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

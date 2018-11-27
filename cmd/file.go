package cmd

import (
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
	"gopkg.in/alecthomas/kingpin.v2"
)

func registerFileCmd(app *kingpin.Application) {
	fileCmd = app.Command(
		"file", "Send a text file line-by-line to a target.",
	)

	// Args:
	fileOpts.Path = fileCmd.Arg(
		"filepath",
		"The path to the text file to be read.",
	).Required().String()

	fileOpts.ConvoID = fileCmd.Arg(
		"conversation ID",
		"The target conversation ID (last portion of a www.messenger.com link).",
	).Default("").String()
}

var (
	fileCmd *kingpin.CmdClause

	fileOpts struct {
		Path, ConvoID *string
	}
)

func file() error {
	// Create generator.
	gen, err := strgen.NewFileReader(*fileOpts.Path)
	if err != nil {
		return ess.AddCtx("creating FileReader", err)
	}

	// Derive convoURL.
	runner := deriveBotRunner(nil)
	convoURL, err := deriveConvoURL(*fileOpts.ConvoID, runner.Prompter)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	bcfg := deriveBotConfig()
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleText(convoURL, gen)
	})
}

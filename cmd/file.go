package cmd

import (
	"github.com/stevenxie/begone/pkg/mbot"
	"github.com/stevenxie/begone/pkg/strgen"
	ess "github.com/unixpickle/essentials"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func registerFileCmd(app *kingpin.Application) {
	fileCmd = app.Command(
		"file", "Send a text file line-by-line to a target.",
	)

	// Args:
	fileCmd.Arg("filepath", "The path to the text file to be read.").Required().
		HintOptions(" ").StringVar(&filePath)

	// Common options:
	registerCommonOpts(fileCmd)
}

var (
	fileCmd  *kingpin.CmdClause
	filePath string
)

func file() error {
	// Create generator.
	gen, err := strgen.NewFileReader(filePath)
	if err != nil {
		return ess.AddCtx("creating FileReader", err)
	}

	// Derive convoURL.
	runner := deriveBotRunner(nil)
	convoURL, err := deriveConvoURL(copts.ConvoID, runner.Prompter)
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

package cmd

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/stevenxie/begone/pkg/mbot"
	"github.com/stevenxie/begone/pkg/strgen"
	ess "github.com/unixpickle/essentials"
)

func registerEmojifyCmd(app *kingpin.Application) {
	emojifyCmd = app.Command(
		"emojify",
		"Spam a target using an emoji generator.",
	).Alias("emoji")

	// Flags:
	emojifyCmd.Flag(
		"mode",
		"Emoji generation method (single, staircase, or unique).",
	).Short('m').Default("single").
		EnumVar(&emojifyMode, "single", "staircase", "unique")

	// Common options:
	registerCommonOpts(emojifyCmd)
}

var (
	emojifyCmd  *kingpin.CmdClause
	emojifyMode string
)

func emojify() error {
	// Create generator.
	gen, err := strgen.NewEmojifier(emojifyMode)
	if err != nil {
		return ess.AddCtx("creating emojifier", err)
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

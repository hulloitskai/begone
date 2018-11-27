package cmd

import (
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
	"gopkg.in/alecthomas/kingpin.v2"
)

func registerEmojifyCmd(app *kingpin.Application) {
	emojifyCmd = app.Command(
		"emojify",
		"Spam a target using an emoji generator.",
	).Alias("emoji")

	// Args:
	emojifyOpts.ConvoID = emojifyCmd.Arg(
		"conversation ID",
		"The target conversation ID (last portion of a www.messenger.com link).",
	).Default("").String()

	// Flags:
	emojifyOpts.Mode = emojifyCmd.Flag(
		"mode",
		"Emoji generation method (single, staircase, or unique).").Short('m').
		Default("single").Enum("single", "staircase", "unique")
}

var (
	emojifyCmd *kingpin.CmdClause

	emojifyOpts struct {
		ConvoID, Mode *string
	}
)

func emojify() error {
	// Create generator.
	gen, err := strgen.NewEmojifier(*emojifyOpts.Mode)
	if err != nil {
		return ess.AddCtx("creating emojifier", err)
	}

	// Derive convoURL.
	runner := deriveBotRunner(nil)
	convoURL, err := deriveConvoURL(*emojifyOpts.ConvoID, runner.Prompter)
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

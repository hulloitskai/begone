package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
)

var fileCmd = &cobra.Command{
	Use:   "file [flags]\n  begone file [flags] FILEPATH [CONVERSATION ID]",
	Short: "Send a text file line-by-line to a target",
	Long:  "File sends text from a file, line-by-line, to a target.",
	Args:  withUsage(cobra.RangeArgs(1, 2)),
	RunE:  fileReader,
}

func fileReader(cmd *cobra.Command, args []string) error {
	// Parse arguments.
	var (
		fpath, convoID string
		p              = interact.NewPrompter()
	)

	fpath = args[0]
	if len(args) > 1 {
		convoID = args[1]
	}

	gen, err := strgen.NewFileReader(fpath)
	if err != nil {
		return ess.AddCtx("creating FileReader", err)
	}

	// Derive Bot config.
	bcfg, err := deriveBotConfig(cmd.Flags())
	if err != nil {
		return ess.AddCtx("deriving bot config", err)
	}

	// Derive convoURL.
	convoURL, err := deriveConvoURL(convoID, p)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	runner := interact.NewBotRunnerWith(p)
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleText(convoURL, gen)
	})
}

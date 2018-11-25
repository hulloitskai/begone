package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
)

func initEmojifyCmd() {
	emojifyCmd.LocalFlags().StringP("mode", "m", "single",
		"emoji generation method ('single', 'staircase', 'unique')")
}

var emojifyCmd = &cobra.Command{
	Use:   "emojify [flags]\n  begone emojify [flags] [CONVERSATION ID]",
	Short: "Spam a target using an emoji generator",
	Long:  "Emojify spams a target using an emoji generator.",
	Args:  withUsage(cobra.MaximumNArgs(1)),
	RunE:  emojify,
}

func emojify(cmd *cobra.Command, args []string) error {
	mode, err := cmd.LocalFlags().GetString("mode")
	if err != nil {
		return err
	}

	// Create generator.
	gen, err := strgen.NewEmojifier(mode)
	if err != nil {
		return ess.AddCtx("creating emojifier", err)
	}

	// Derive Bot config.
	bcfg, err := deriveBotConfig(cmd.Flags())
	if err != nil {
		return ess.AddCtx("deriving bot config", err)
	}

	// Derive convoURL.
	var (
		convoID string
		runner  = interact.NewBotRunner()
	)
	if len(args) > 0 {
		convoID = args[0]
	}
	convoURL, err := deriveConvoURL(convoID, runner.Prompter)
	if err != nil {
		return err
	}

	// Configure and run BotRunner.
	if err = runner.Configure(bcfg); err != nil {
		return ess.AddCtx("configuring BotRunner", err)
	}
	return runner.Run(func(bot *mbot.Bot) error {
		return bot.CycleText(convoURL, gen)
	})
}

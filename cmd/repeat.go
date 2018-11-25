package cmd

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"
	"github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/begone/strgen"
	ess "github.com/unixpickle/essentials"
)

var repeatCmd = &cobra.Command{
	Use:   "repeat [flags]\n  begone repeat [flags] [MESSAGE] [CONVERSATION ID]",
	Short: "Spam a target with a message",
	Long:  "Repeat spams a target with a message.",
	Args:  withUsage(cobra.MaximumNArgs(2)),
	RunE:  repeat,
}

func repeat(cmd *cobra.Command, args []string) error {
	// Parse arguments.
	var (
		msg, convoID string
		p            = interact.NewPrompter()
	)
	switch len(args) {
	case 2:
		convoID = args[1]
		fallthrough
	case 1:
		msg = args[0]
	}

	// Validate arguments.
	if msg == "" {
		var err error
		if msg, err = queryRepeatMessage(p); err != nil {
			return ess.AddCtx("querying for message", err)
		}
	}

	// Create generator.
	gen := strgen.NewRepeater(msg)

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

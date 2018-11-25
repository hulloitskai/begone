package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	inter "github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// deriveBotConfig creates an mbot.Config from flags. It does not fill out the
// Username or Password fields in the Config.
func deriveBotConfig(flags *pflag.FlagSet) (*mbot.Config, error) {
	var (
		mcfg = mbot.NewConfig()
		err  error
	)

	if mcfg.Cycles, err = flags.GetInt("cycles"); err != nil {
		return nil, err
	}
	if mcfg.Delay, err = flags.GetInt("delay"); err != nil {
		return nil, err
	}

	return mcfg, nil
}

// deriveConvoURL derives a convoURL using id and p.
func deriveConvoURL(id string, p *inter.Prompter) (convoURL string, err error) {
	if p == nil {
		p = inter.NewPrompter()
	}

	if id != "" {
		convoURL = fmt.Sprintf("%s/t/%s", fbmsgr.BaseURL, id)
	} else if convoURL, err = p.QueryConvoURL(); err != nil {
		return "", ess.AddCtx("querying for conversation URL", err)
	}
	return convoURL, nil
}

func withUsage(pa cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := pa(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			if err := cmd.Usage(); err != nil {
				ess.Die(err)
			}
			os.Exit(1)
		}
		return nil
	}
}

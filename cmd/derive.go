package cmd

import (
	"fmt"
	"strings"

	"github.com/stevenxie/begone/internal/interact"
	inter "github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/pkg/mbot"
	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// deriveBotConfig creates an mbot.Config from global app opts. It does not
// fill out the Username or Password fields in the Config.
func deriveBotConfig() *mbot.Config {
	mcfg := mbot.NewConfig()
	mcfg.Cycles = copts.Cycles
	mcfg.Delay = copts.Delay
	mcfg.AssumeUser = copts.AssumeUser
	return mcfg
}

// deriveBotRunner creates an interact.BotRunner from global app opts.
func deriveBotRunner(p *interact.Prompter) *interact.BotRunner {
	runner := interact.NewBotRunnerWith(p)
	runner.Debug = copts.Debug

	if copts.NoFancy {
		runner.Mode = inter.Reduced
	}
	return runner
}

// deriveConvoURL derives a convoURL using id and p.
func deriveConvoURL(id string, p *inter.Prompter) (convoURL string, err error) {
	if p == nil {
		p = inter.NewPrompter()
	}

	if id != "" {
		if strings.HasPrefix(id, "https://www.messenger.com") {
			return "", fmt.Errorf("expected convo ID, not convo URL (got '%s')", id)
		}
		return fmt.Sprintf("%s/t/%s", fbmsgr.BaseURL, id), nil
	}

	if convoURL, err = p.QueryConvoURL(); err != nil {
		return "", ess.AddCtx("querying for conversation URL", err)
	}
	return convoURL, nil
}

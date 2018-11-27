package cmd

import (
	"fmt"

	"github.com/stevenxie/begone/internal/interact"
	inter "github.com/stevenxie/begone/internal/interact"
	"github.com/stevenxie/begone/mbot"
	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// deriveBotConfig creates an mbot.Config from global app opts. It does not
// fill out the Username or Password fields in the Config.
func deriveBotConfig() *mbot.Config {
	mcfg := mbot.NewConfig()
	mcfg.Cycles = copts.Cycles
	mcfg.Delay = copts.Delay
	return mcfg
}

// deriveBotRunner creates an interact.BotRunner from global app opts.
func deriveBotRunner(p *interact.Prompter) *interact.BotRunner {
	runner := interact.NewBotRunnerWith(p)
	runner.Debug = copts.Debug
	return runner
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

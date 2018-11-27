package interact

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/stevenxie/begone/internal/config"
	"github.com/stevenxie/begone/mbot"
	ess "github.com/unixpickle/essentials"
)

// A BotRunner knows how to configure and run a Bot.
type BotRunner struct {
	*Prompter
	Debug bool

	rng *rand.Rand
	Bot *mbot.Bot
}

// NewBotRunner returns a new BotRunner.
func NewBotRunner() *BotRunner {
	return NewBotRunnerWith(NewPrompter())
}

// NewBotRunnerWith returns a new BotRunner that interacts with the user using
// Prompter p.
//
// If p is nil, it will be set to a Prompter with the default configuration.
func NewBotRunnerWith(p *Prompter) *BotRunner {
	if p == nil {
		p = NewPrompter()
	}

	src := rand.NewSource(time.Now().Unix())
	return &BotRunner{
		Prompter: p,
		rng:      rand.New(src),
	}
}

// Configure configures the Bot using bcfg. It attempts to fill out the
// Username and Password fields of bcfg using values from a config.Config.
//
// If the resulting bcfg is not valid, an error will be returned.
func (br *BotRunner) Configure(bcfg *mbot.Config) error {
	// Ensure bcfg is non-nil.
	if bcfg == nil {
		bcfg = mbot.NewConfig()
	}

	// Query for any missing values.
	cfg, err := config.Load()
	if err != nil {
		return ess.AddCtx("interact: loading config file", err)
	}

	// Prevent overriding bcfg values.
	if bcfg.Username != "" {
		cfg.Username = bcfg.Username
	}
	if bcfg.Password != "" {
		cfg.Password = bcfg.Password
	}

	if err = br.QueryMissing(cfg, false); err != nil {
		return ess.AddCtx("interact: querying for missing values", err)
	}

	// Amend bcfg with new credentials.
	bcfg.Username = cfg.Username
	bcfg.Password = cfg.Password

	// Build and configure logger if applicable.
	br.Bot, err = bcfg.Build()
	if br.Debug {
		br.Bot.Logger = log.New(os.Stderr, "", 0)
	}
	return err
}

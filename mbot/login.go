package mbot

import (
	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// Login logs the Bot into the FB Messenger account specified in b.Config.
func (b *Bot) Login() error {
	if err := b.Cfg.Validate(); err != nil {
		return ess.AddCtx("mbot: invalid config", err)
	}

	sess, err := fbmsgr.Auth(b.Cfg.Username, b.Cfg.Password)
	if err != nil {
		return ess.AddCtx("mbot: authenticating user", err)
	}
	b.Session = sess

	return nil
}

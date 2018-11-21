package mbot

import (
	"time"

	"github.com/stevenxie/begone/strgen"
	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// A Bot is capable of messaging people on Facebook Messenger.
type Bot struct {
	// Used for generating messages.
	strgen.Generator

	// A Facebook Messenger client session.
	*fbmsgr.Session

	// A live-updating count of the number of messages sent using the Begone
	// method.
	Counter chan int

	Cfg        *Config
	killswitch bool // stops Begone loop if active
}

// NewBot returns a new Bot with a nil Session.
func NewBot(cfg *Config, gen strgen.Generator) *Bot {
	return &Bot{Cfg: cfg, Generator: gen}
}

// Login logs the Bot into the FB Messenger account specified in b.Config.
func (b *Bot) Login() error {
	if err := b.Cfg.CheckValidity(); err != nil {
		return ess.AddCtx("mbot: invalid config", err)
	}

	sess, err := fbmsgr.Auth(b.Cfg.Username, b.Cfg.Password)
	if err != nil {
		return ess.AddCtx("mbot: authenticating user", err)
	}
	b.Session = sess

	return nil
}

const (
	// definesToken indicates the section of HTML that defines FB IDs in terms of
	// convoIDs.
	definesToken = "require(\"ServerJSDefine\").handleDefines"
)

// Begone is a Bot's primary action; it repeatedly spams the FB Messenger
// thread identified by convoID until it has reached its target cycle count
// (b.Cfg.Cycles), or until its Generator has no more content to produce.
func (b *Bot) Begone(convoID string) error {
	if b.Session == nil {
		if err := b.Login(); err != nil {
			return err
		}
	}

	// Parse convoID into user / group fbid.
	var (
		fbid    string
		isGroup bool
	)
	if fchar := convoID[0]; ('0' <= fchar) && (fchar <= '9') {
		fbid = convoID
		isGroup = true
	} else {
		var err error
		if fbid, err = b.parseUserConvoID(convoID); err != nil {
			return ess.AddCtx("mbot: looking up corresponding fbid", err)
		}
	}

	// The Begone™ Spam Loop.
	var count int
	for b.HasMore() && (count != b.Cfg.Cycles) && (!b.killswitch) {
		// Message generation cycle.
		msg, err := b.Generate()
		if err != nil {
			return ess.AddCtx("mbot: generating message", err)
		}
		if isGroup {
			_, err = b.SendGroupText(fbid, msg)
		} else {
			_, err = b.SendText(fbid, msg)
		}
		if err != nil {
			return ess.AddCtx("mbot: sending message", err)
		}

		// Post-cycle operations.
		if b.Cfg.Delay > 0 {
			time.Sleep(time.Duration(b.Cfg.Delay) * time.Millisecond)
		}
		count++
		if b.Counter != nil { // update counter if applicable
			b.Counter <- count
		}
	}

	if b.Counter != nil { // send close signal if applicable
		close(b.Counter)
	}
	return nil
}

// Kill activates the Bot's killswitch, stopping a active Begone loop.
func (b *Bot) Kill() {
	b.killswitch = true
}

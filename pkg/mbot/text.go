package mbot

import (
	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// A StringGenerator can continuously generate strings.
type StringGenerator interface {
	Generate() (string, error)
	HasMore() bool
}

type stringCycler struct {
	StringGenerator
	FBID    string
	IsGroup bool
}

// Cycle implements Cycler's Cycle method.
func (sc *stringCycler) Cycle(session *fbmsgr.Session) error {
	msg, err := sc.Generate()
	if err != nil {
		return ess.AddCtx("mbot: generating message", err)
	}

	if sc.IsGroup {
		_, err = session.SendGroupText(sc.FBID, msg)
	} else {
		_, err = session.SendText(sc.FBID, msg)
	}
	return err
}

// Finished implements Cycler's Continue method.
func (sc *stringCycler) Finished() bool {
	return !sc.StringGenerator.HasMore()
}

// CycleText repeatedly sends text to a conversation identified by convoURL
// using gen.
func (b *Bot) CycleText(convoURL string, gen StringGenerator) error {
	entity, err := b.extractFBEntity(convoURL)
	if err != nil {
		return ess.AddCtx("mbot: looking up fbid", err)
	}

	return b.CycleUsing(&stringCycler{
		StringGenerator: gen,
		FBID:            entity.ID,
		IsGroup:         entity.IsGroup,
	})
}

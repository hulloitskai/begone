package mbot

import (
	"fmt"
	"strings"
	"time"

	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// A Cycler is capable of performing some action using an fbmsgr.Session.
//
// It may also indicate that there are no cycles left when its Finished method
// returns true.
type Cycler interface {
	// Cycle performs some action using an fbmsgr.Session.
	//
	// If an error returns, it is considered a send failure if it contains the
	// text "fbmsgr:".
	Cycle(session *fbmsgr.Session) error
	Finished() bool
}

// CycleUsing repeated performs c.Cycle until c.Finished returns true, or until
// the target cycle count (b.Cfg.Cycles) has been reached (whichever happens
// first).
//
// If b.Counter is set, it will be nil after Do has finished.
func (b *Bot) CycleUsing(c Cycler) error {
	if b.Session == nil {
		if err := b.Login(); err != nil {
			return err
		}
	}

	// Close b.Counter (and reset it) if applicable upon returning.
	defer func() {
		if b.Counter != nil {
			close(b.Counter)
			b.Counter = nil
		}
	}()

	// The Begone™ Spam Loop.
	var count, sendFails int
	for !c.Finished() && (count != b.Cfg.Cycles) {
		// Perform cycle action.
		if err := c.Cycle(b.Session); err != nil {
			if strings.Contains(err.Error(), "fbmsgr:") { // is a send fail
				// Fail early if this is the first cycle.
				if count == 0 {
					return fmt.Errorf("mbot: encountered send failure on first cycle: %v",
						err)
				}

				// If there have been too many consecutive send fails, crash.
				sendFails++
				if sendFails > b.Cfg.MaxSendFails {
					return fmt.Errorf("mbot: exceed maximum consecutive send failures "+
						"(%d): %v", b.Cfg.MaxSendFails, err)
				}

				if b.Logger != nil {
					b.Logger.Printf("Error during send (%d of %d maximum): %v", sendFails,
						b.Cfg.MaxSendFails, err)
				}

				time.Sleep(time.Duration(b.Cfg.SendFailDelay) * time.Millisecond)
				continue
			}
			return ess.AddCtx("mbot: cycle error", err)
		}

		// No send fail occurred; reset our consecutive send failures tracker.
		sendFails = 0

		// Post-action operations:
		if b.Cfg.Delay > 0 {
			time.Sleep(time.Duration(b.Cfg.Delay) * time.Millisecond)
		}
		count++
		if b.Counter != nil { // update counter if applicable
			b.Counter <- count
		}
	}
	return nil
}

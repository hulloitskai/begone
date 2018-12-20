package mbot

import (
	"log"

	"github.com/stevenxie/fbmsgr"
)

// A Bot is capable of continuously messaging people on Facebook Messenger.
type Bot struct {
	// A Facebook Messenger client session.
	*fbmsgr.Session

	// A live-updating count of the number of messages sent during a spam loop.
	Counter chan int

	Cfg    *Config
	Logger *log.Logger
}

// NewBot returns a new Bot with a nil Session.
func NewBot(cfg *Config, logger *log.Logger) *Bot {
	return &Bot{Cfg: cfg, Logger: logger}
}

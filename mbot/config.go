package mbot

import (
	"errors"
	"log"
)

// Config contains options that control the behaviour of a Bot.
type Config struct {
	// Cycles is the number of action cycles the Bot should perform before
	// stopping.
	//
	// The Bot may stop before this number if its Cycler reports that it has
	// finished.
	//
	// A value of -1 indicates that a bot should run indefinitely, generating
	// as many messages as Generator can produce.
	Cycles int

	// Delay is the delay (in ms) between each message (useful for avoiding spam
	// detection).
	Delay int

	// MaxSendFails is the maximum number of consecutive send failures the Bot
	// will tolerate before erroring out.
	MaxSendFails int

	// SendFailDelay is a delay (in ms) that occurs after a send failure as a
	// sort of 'cool-down period'.
	SendFailDelay int

	Username, Password string
}

// NewConfig returns a Config with default values.
func NewConfig() *Config {
	return &Config{Cycles: -1, Delay: 250, MaxSendFails: 2, SendFailDelay: 1250}
}

// Validate checks if the Config is valid, i.e. if a Bot can run properly if
// provided with the Config.
//
// If the Config is not valid, an error will be returned.
func (cfg *Config) Validate() error {
	if cfg.Username == "" {
		return errors.New("mbot: Username must be non-empty")
	}
	if cfg.Password == "" {
		return errors.New("mbot: Password must be non-empty")
	}
	if cfg.Delay < 0 {
		return errors.New("mbot: Delay must be non-negative")
	}
	if cfg.Cycles < -1 {
		return errors.New("mbot: Cycles must be -1 or greater")
	}
	if cfg.MaxSendFails < 0 {
		return errors.New("mbot: MaxSendFails must be non-negative")
	}
	if cfg.SendFailDelay < 0 {
		return errors.New("mbot: SendFailDelay must be non-negative")
	}
	return nil
}

// Build builds a Bot configured using this Config.
//
// An error will be returned if the Config is not valid.
func (cfg *Config) Build() (*Bot, error) {
	return cfg.BuildWith(nil)
}

// BuildWith builds a Bot configured using this Config and logger for
// reporting send failures and other errors.
//
// An error will be returned if the Config is not valid.
func (cfg *Config) BuildWith(logger *log.Logger) (*Bot, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return NewBot(cfg, logger), nil
}

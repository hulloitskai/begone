package mbot

import (
	"errors"
)

// Config contains options that control the behaviour of a Bot.
type Config struct {
	// The maximum number of message generation cycles the Bot should
	// perform before stopping. The Bot may stop before this number if Generator
	// reports that it has no more output to produce.
	//
	// A value of -1 indicates that a bot should run indefinitely, generating
	// as many messages as Generator can produce.
	Cycles int

	// The delay between each message (to avoid spam detection).
	Delay int

	Username, Password string
}

// NewConfig returns a Config with default values (Cycles = -1, Delay = 10).
func NewConfig() *Config {
	return &Config{Cycles: -1, Delay: 10}
}

// CheckValidity checks if the Config is valid, i.e. if a Bot can run properly
// if provided with the Config.
//
// If the Config is not valid, an error will be returned.
func (c *Config) CheckValidity() error {
	if c.Username == "" {
		return errors.New("mbot: Username must be non-empty")
	}
	if c.Password == "" {
		return errors.New("mbot: Password must be non-empty")
	}
	if c.Delay < 0 {
		return errors.New("mbot: Delay must be non-negative")
	}
	if c.Cycles < -1 {
		return errors.New("mbot: Cycles must be either -1, or a natural number")
	}
	return nil
}

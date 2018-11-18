package messenger

import (
	"fmt"

	"github.com/stevenxie/begone/loader"
	"github.com/tebeka/selenium"
)

// Engine is capable of executing operations on Facebook Messenger using
// Selenium.
type Engine struct {
	*Config

	wd  selenium.WebDriver
	svc *selenium.Service
	l   *loader.Loader
}

// NewEngine makes a new Engine.
func NewEngine(cfg *Config, l *loader.Loader) *Engine {
	return &Engine{Config: cfg, l: l}
}

// StartSession starts the Engine session, and navigates to the target URL.
//
// This is a no-op if the Engine session has already started.
func (e *Engine) StartSession() error {
	if e.wd != nil { // Engine already in-session
		return nil
	}

	if err := e.l.Ensure(); err != nil {
		return fmt.Errorf("messenger: failed to ensure valid wd "+
			"installation: %v", err)
	}

	var (
		caps = selenium.Capabilities{"browserName": "chrome"}
		err  error
	)

	if e.svc, err =
		selenium.NewChromeDriverService(e.l.DriverPath, 4444); err != nil {
		return fmt.Errorf("messenger: error making Selenium svc: %v", err)
	}
	if e.wd, err = selenium.NewRemote(caps, ""); err != nil {
		return fmt.Errorf("messenger: error making Selenium remote: %v", err)
	}

	return nil
}

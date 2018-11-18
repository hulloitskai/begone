package messenger

import (
	"fmt"

	"github.com/tebeka/selenium"
)

const (
	efID = "email" // email field ID
	pfID = "pass"  // password field ID
)

// Login logs the user into FB Messenger on messenger.com.
// Creates a new Selenium session, if none currently exists.
func (e *Engine) Login() error {
	if err := e.CheckLogin(); err != nil {
		return err
	}

	// Start a new session, if necessary.
	if e.wd == nil {
		if err := e.StartSession(); err != nil {
			return err
		}
	}

	// Navigate to login page.
	e.wd.Get(e.MessengerURL())

	// Enter username.
	if err := e.wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		elem, err := wd.FindElement(selenium.ByID, efID)
		if err != nil {
			return false, err
		}

		displayed, err := elem.IsDisplayed()
		if err != nil {
			return false, err
		}

		return displayed, nil
	}); err != nil {
		return fmt.Errorf("messenger: error while waiting for email field: %v", err)
	}

	emailField, err := e.wd.FindElement(selenium.ByID, efID)
	if err != nil {
		return fmt.Errorf("messenger: error finding email field: %v", err)
	}
	if err = emailField.SendKeys(e.Config.User); err != nil {
		return fmt.Errorf("messenger: error sending keys on email field: %v", err)
	}

	// Enter password.
	passField, err := e.wd.FindElement(selenium.ByID, pfID)
	if err != nil {
		return fmt.Errorf("messenger: error finding password field: %v", err)
	}
	if err = passField.SendKeys(e.Config.Pass); err != nil {
		return fmt.Errorf("messenger: error sending keys on password field: %v",
			err)
	}

	// Submit login.
	passField.SendKeys(selenium.EnterKey)
	return nil
}

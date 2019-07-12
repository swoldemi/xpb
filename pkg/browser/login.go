package browser

import (
	"time"

	"github.com/tebeka/selenium"
)

// TypeLoginEmail types an account's login email.
func (c *ChromeBrowser) TypeLoginEmail(email string) error {
	emailField, err := c.WebDriver.FindElement(selenium.ByName, LoginEmailSelector)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(email)
	if err != nil {
		return err
	}
	return nil
}

// SubmitEmail submits the email field to progress to the login field.
func (c *ChromeBrowser) SubmitEmail() error {
	nextBtn, err := c.WebDriver.FindElement(selenium.ByCSSSelector, EmailSubmitSelector)
	if err != nil {
		return err
	}

	err = nextBtn.Click()
	time.Sleep(time.Second * 1)
	if err != nil {
		return err
	}
	return nil
}

// TypeLoginPassword types an account's login password.
func (c *ChromeBrowser) TypeLoginPassword(password string) error {
	passField, err := c.WebDriver.FindElement(selenium.ByName, HostPasswordSelector)
	if err != nil {
		return err
	}

	err = passField.SendKeys(password)
	if err != nil {
		return err
	}
	return nil
}

// SubmitPassword submits the login form.
func (c *ChromeBrowser) SubmitPassword() error {
	nextBtn, err := c.WebDriver.FindElement(selenium.ByCSSSelector, LoginSubmitSelector)
	if err != nil {
		return err
	}

	err = nextBtn.Click()
	time.Sleep(time.Second * 1)
	if err != nil {
		return err
	}
	return nil
}

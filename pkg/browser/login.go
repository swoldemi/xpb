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
	return emailField.SendKeys(email)
}

// SubmitEmail submits the email field to progress to the login field.
func (c *ChromeBrowser) SubmitEmail() error {
	nextBtn, err := c.WebDriver.FindElement(selenium.ByCSSSelector, EmailSubmitSelector)
	if err != nil {
		return err
	}

	if err := nextBtn.Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	return nil
}

// TypeLoginPassword types an account's login password.
func (c *ChromeBrowser) TypeLoginPassword(password string) error {
	passField, err := c.WebDriver.FindElement(selenium.ByName, LoginPasswordSelector)
	if err != nil {
		return err
	}
	return passField.SendKeys(password)
}

// SubmitPassword submits the login form.
func (c *ChromeBrowser) SubmitPassword() error {
	nextBtn, err := c.WebDriver.FindElement(selenium.ByCSSSelector, LoginSubmitSelector)
	if err != nil {
		return err
	}

	if err := nextBtn.Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	return nil
}

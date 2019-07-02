package browser

import (
	"time"
	"github.com/tebeka/selenium"
)

func (c *ChromeBrowser) typeHostEmail() error {
	emailField, err := c.WebDriver.FindElement(selenium.ByName, HostEmailSelector)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(c.Config.NamedHostEmail)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) submitEmail() error {
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

func (c *ChromeBrowser) typeHostPassword() error {
	passField, err := c.WebDriver.FindElement(selenium.ByName, HostPasswordSelector)
	if err != nil {
		return err
	}

	err = passField.SendKeys(c.Config.HostPass)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) submitPassword() error {
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

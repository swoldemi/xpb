package browser

import (
	"time"

	"github.com/tebeka/selenium"
)

func (g *GCPBrowser) typeHostEmail() error {
	emailField, err := g.WebDriver.FindElement(selenium.ByName, HostEmailSelector)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(g.Config.NamedHostEmail)
	if err != nil {
		return err
	}
	return nil
}

func (g *GCPBrowser) submitEmail() error {
	nextBtn, err := g.WebDriver.FindElement(selenium.ByCSSSelector, EmailSubmitSelector)
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

func (g *GCPBrowser) typeHostPassword() error {
	passField, err := g.WebDriver.FindElement(selenium.ByName, HostPasswordSelector)
	if err != nil {
		return err
	}

	err = passField.SendKeys(g.Config.HostPass)
	if err != nil {
		return err
	}
	return nil
}

func (g *GCPBrowser) submitPassword() error {
	nextBtn, err := g.WebDriver.FindElement(selenium.ByCSSSelector, LoginSubmitSelector)
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

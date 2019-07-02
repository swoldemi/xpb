package browser

import (
	"time"

	"github.com/tebeka/selenium"
)

func (g *GCPBrowser) clickAdd() error {
	addBtn, err := g.WebDriver.FindElement(selenium.ByXPATH, AddUserBtnXPATH)
	if err != nil {
		return err
	}

	err = addBtn.Click()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 8) // Need to wait for the drawer to open
	return nil
}

func (g *GCPBrowser) typeGuestEmail() error {
	emailField, err := g.WebDriver.FindElement(selenium.ByCSSSelector, GuestEmailSelector)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(g.Config.NamedGuestEmail)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(selenium.EnterKey)
	if err != nil {
		return err
	}
	return nil
}

func (g *GCPBrowser) clickHeader() error {
	headerEl, err := g.WebDriver.FindElement(selenium.ByCSSSelector, HeaderSelector)
	if err != nil {
		return err
	}

	err = headerEl.Click()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 1) // Wait for the option to be gone
	return nil
}

func (g *GCPBrowser) addFirstRole() error {
	roleField, err := g.WebDriver.FindElement(selenium.ByCSSSelector, FirstRoleSelector)
	if err != nil {
		return err
	}
	err = roleField.Click()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1) // Wait for the menu to be visible

	menuInput, err := g.WebDriver.FindElement(selenium.ByXPATH, RolesMenuXPATH)
	if err != nil {
		return err
	}

	err = menuInput.SendKeys(OwnerRole)
	if err != nil {
		return err
	}

	// The filter has completed, click the role
	ownerRoleEl, err := g.WebDriver.FindElement(selenium.ByCSSSelector, OwnerRoleSelector)
	if err != nil {
		return err
	}
	err = ownerRoleEl.Click()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 1) // Wait for the menu to close
	return nil
}

func (g *GCPBrowser) clickAddAnother() error {
	addRoleBtn, err := g.WebDriver.FindElement(selenium.ByXPATH, AddRoleBtnXPATH)
	if err != nil {
		return err
	}

	err = addRoleBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

func (g *GCPBrowser) addSecondRole() error {
	roleFields, err := g.WebDriver.FindElements(selenium.ByXPATH, RoleInputsXPATH)
	if err != nil {
		return err
	}

	// Click the second element that matches the selector
	err = roleFields[1].Click()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1) // Wait for the menu to be visible

	menuInput, err := g.WebDriver.FindElement(selenium.ByXPATH, RolesMenuXPATH)
	if err != nil {
		return err
	}

	// Because this is the second time we're opening the menu, clear it
	err = menuInput.Clear()
	if err != nil {
		return err
	}

	err = menuInput.SendKeys(BillingRole)
	if err != nil {
		return err
	}

	// The menu has been filtered to show ONLY the Project Billing Manager Role,
	// but because the Project Billing Manager is assigned a, seemingly,
	// random option ID on refresh, we can't give it a static selector.
	// Instead, do a search on the document tree. Stop when the
	// inner text of the element matches the name of the role
	optionElements, err := g.WebDriver.FindElements(selenium.ByXPATH, MatOptionsXPATH)
	if err != nil {
		return err
	}

	// Find the only one that is visible, the Project Billing Manager button and click it
	for _, element := range optionElements {
		displayed, err := element.IsDisplayed()
		if err != nil {
			return err
		}

		if displayed {
			err = element.Click()
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (g *GCPBrowser) submitGuestInvite() error {
	submitBtn, err := g.WebDriver.FindElement(selenium.ByXPATH, RolesSubmitXPATH)
	if err != nil {
		return err
	}

	err = submitBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

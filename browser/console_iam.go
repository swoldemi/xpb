package browser

import (
	"time"

	"github.com/tebeka/selenium"
)

func (c *ChromeBrowser) clickAdd() error {
	addBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, AddUserBtnXPATH)
	if err != nil {
		return err
	}

	err = addBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) typeGuestEmail() error {
	emailField, err := c.WebDriver.FindElement(selenium.ByCSSSelector, GuestEmailSelector)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(c.Config.NamedGuestEmail)
	if err != nil {
		return err
	}

	err = emailField.SendKeys(selenium.EnterKey)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) clickHeader() error {
	headerEl, err := c.WebDriver.FindElement(selenium.ByCSSSelector, HeaderSelector)
	if err != nil {
		return err
	}

	err = headerEl.Click()
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) addFirstRole() error {
	roleField, err := c.WebDriver.FindElement(selenium.ByCSSSelector, FirstRoleSelector)
	if err != nil {
		return err
	}
	err = roleField.Click()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1) // Wait for the menu to be visible

	menuInput, err := c.WebDriver.FindElement(selenium.ByXPATH, RolesMenuXPATH)
	if err != nil {
		return err
	}

	err = menuInput.SendKeys(OwnerRole)
	if err != nil {
		return err
	}

	// The filter has completed, click the role
	ownerRoleEl, err := c.WebDriver.FindElement(selenium.ByCSSSelector, OwnerRoleSelector)
	if err != nil {
		return err
	}
	err = ownerRoleEl.Click()
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) clickAddAnother() error {
	addRoleBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, AddRoleBtnXPATH)
	if err != nil {
		return err
	}

	err = addRoleBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

func (c *ChromeBrowser) addSecondRole() error {
	roleFields, err := c.WebDriver.FindElements(selenium.ByXPATH, RoleInputsXPATH)
	if err != nil {
		return err
	}

	// Click the second element that matches the selector
	err = roleFields[1].Click()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1) // Wait for the menu to be visible

	menuInput, err := c.WebDriver.FindElement(selenium.ByXPATH, RolesMenuXPATH)
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
	optionElements, err := c.WebDriver.FindElements(selenium.ByXPATH, MatOptionsXPATH)
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

func (c *ChromeBrowser) submitGuestInvite() error {
	submitBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, RolesSubmitXPATH)
	if err != nil {
		return err
	}

	err = submitBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

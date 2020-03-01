package browser

import (
	"time"

	"github.com/tebeka/selenium"
)

// ClickAdd clicks the add button on the IAM tab of the GCP console.
func (c *ChromeBrowser) ClickAdd() error {
	addBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, AddUserBtnXPATH)
	if err != nil {
		return err
	}
	return addBtn.Click()
}

// TypeGuestEmail types the guest's email in the 'Add members to "Project"' drawer.
func (c *ChromeBrowser) TypeGuestEmail() error {
	emailField, err := c.WebDriver.FindElement(selenium.ByCSSSelector, GuestEmailSelector)
	if err != nil {
		return err
	}
	if err := emailField.SendKeys(c.Config.GuestEmail); err != nil {
		return err
	}
	return emailField.SendKeys(selenium.EnterKey)
}

// ClickHeader clicks the 'Add members, roles to "Project" project' header to hide
// the Material options drop down for the AddFirstRole method.
func (c *ChromeBrowser) ClickHeader() error {
	headerEl, err := c.WebDriver.FindElement(selenium.ByCSSSelector, HeaderSelector)
	if err != nil {
		return err
	}
	return headerEl.Click()
}

// AddFirstRole adds the Owner role for the guest being added.
func (c *ChromeBrowser) AddFirstRole() error {
	roleField, err := c.WebDriver.FindElement(selenium.ByCSSSelector, FirstRoleSelector)
	if err != nil {
		return err
	}
	if err := roleField.Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1) // Wait for the menu to be visible

	menuInput, err := c.WebDriver.FindElement(selenium.ByXPATH, RolesMenuXPATH)
	if err != nil {
		return err
	}
	if err := menuInput.SendKeys(OwnerRole); err != nil {
		return err
	}

	// The filter has completed, click the role
	ownerRoleEl, err := c.WebDriver.FindElement(selenium.ByCSSSelector, OwnerRoleSelector)
	if err != nil {
		return err
	}
	return ownerRoleEl.Click()
}

// ClickAddAnother adds another role field for the guest.
func (c *ChromeBrowser) ClickAddAnother() error {
	addRoleBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, AddRoleBtnXPATH)
	if err != nil {
		return err
	}
	return addRoleBtn.Click()
}

// AddSecondRole adds the Project Billing Manager role for the guest being added.
func (c *ChromeBrowser) AddSecondRole() error {
	roleFields, err := c.WebDriver.FindElements(selenium.ByXPATH, RoleInputsXPATH)
	if err != nil {
		return err
	}

	// Click the second element that matches the selector
	if err := roleFields[1].Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1) // Wait for the menu to be visible

	menuInput, err := c.WebDriver.FindElement(selenium.ByXPATH, RolesMenuXPATH)
	if err != nil {
		return err
	}

	// Because this is the second time we're opening the menu, clear it
	if err := menuInput.Clear(); err != nil {
		return err
	}

	if err := menuInput.SendKeys(BillingRole); err != nil {
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

	// Do the search for the visible role, the Project Billing Manager button, and click it
	for _, element := range optionElements {
		displayed, err := element.IsDisplayed()
		if err != nil {
			return err
		}

		if displayed {
			if err := element.Click(); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// SubmitGuestInvite clicks the save button at the bottom of the
// 'Add members to "Project"' drawer. If the guest
// is already bound as a poilcy on the current project,
// this may overwrite existing roles. If no new changes
// are made to the user that is already a member of the project,
// no changes will be made.
func (c *ChromeBrowser) SubmitGuestInvite() error {
	submitBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, RolesSubmitXPATH)
	if err != nil {
		return err
	}
	return submitBtn.Click()
}

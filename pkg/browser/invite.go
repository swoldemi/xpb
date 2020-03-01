package browser

import (
	"github.com/tebeka/selenium"
)

// ClickInvite clicks accept for the host's invite to the guest.
func (c *ChromeBrowser) ClickInvite() error {
	acceptBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, InviteAcceptXPATH)
	if err != nil {
		return err
	}
	return acceptBtn.Click()
}

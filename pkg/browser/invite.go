package browser

// AcceptInvite accept's the host's invite for the guest.
func (c *ChromeBrowser) AcceptInvite() error {
	acceptBtn, err := c.WebDriver.FindElement(selenium.ByXPATH, InviteAcceptXPATH)
	if err != nil {
		return err
	}

	err = acceptBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

package cmd

import (
	"github.com/swoldemi/xpb/pkg/accounts"
	"github.com/swoldemi/xpb/pkg/browser"
	"github.com/swoldemi/xpb/pkg/config"
	"github.com/swoldemi/xpb/pkg/log"
)

// Execute runs the XPB invocation using the provided configuration.
// TODO: Turn into actual CLI command using spf13/
func Execute(cfg *config.Config) error {
	log.SetLevel("trace")
	acc, err := accounts.New(cfg)
	if err != nil {
		return err
	}
	billingAccount, err := acc.Verify()
	if err != nil {
		return err
	}
	chrome, err := browser.New(cfg)
	if err != nil {
		return err
	}

	defer func() {
		if err := chrome.SeleniumService.Stop(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	defer func() {
		if err := chrome.WebDriver.Quit(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	if err := chrome.LoginHost(); err != nil {
		return err
	}
	if err := chrome.InviteGuest(); err != nil {
		return err
	}
	if err := chrome.AcceptInvite(); err != nil {
		return err
	}

	// Guest has been invited to the host's GCP console
	// and has accepted the invite. Link the guest's
	// valid billing account
	if err := acc.Link(billingAccount); err != nil {
		return err
	}
	return nil
}

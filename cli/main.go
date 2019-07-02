package main

import (
	"time"

	xpb "github.com/swoldemi/xpb"
	// "github.com/swoldemi/xpb/browser"
	"github.com/swoldemi/xpb/email"
)

const (
	version = "1.0.0-rc.1"
)

func main() {
	config := &xpb.Config{
		NamedHostEmail:      "nickelapi@gmail.com",
		NamedGuestEmail:     "machserve.io@gmail.com",
		IntermdiateTimeout:  time.Second * 10,
		HostPass:            "E5BB55AD2B9C54BB3264FC862513E",
		HostProjectID:       "nickel-api",
		HostKeyFilePath:     "cli/xpb-host.json",
		GuestGmailCredsPath: "cli/guest-gmail-credentials.json",
		SeleniumRemotePort:  8080,
		PollInterval:        time.Millisecond * 500,
		ChromeDriverPath:    "drivers/chromedriver.exe",
		SeleniumPath:        "drivers/selenium-server-standalone.jar",
		Debug:               true,
	}

	// Before doing anything, verify
	// that the guest has a valid billing account
	// _, err := xpb.New(config)
	// if err != nil {
	// 	xpb.Fataler(err)
	// }

	// err = x.ValidateGuest()
	// if err != nil {
	// 	xpb.Fataler(err)
	// }

	// chrome, err := browser.New(config)
	// if err != nil {
	// 	xpb.Fataler(err)
	// }

	// defer func() {
	// 	serr := chrome.SeleniumService.Stop()
	// 	if serr != nil {
	// 		xpb.Fataler(err)
	// 	}
	// }()

	// defer func() {
	// 	qerr := chrome.WebDriver.Quit()
	// 	if qerr != nil {
	// 		xpb.Fataler(err)
	// 	}
	// }()

	// err = chrome.LoginHost()
	// if err != nil {
	// 	xpb.Fataler(err)
	// }

	// err = chrome.InviteGuest()
	// if err != nil {
	// 	xpb.Fataler(err)
	// }

	inbox, err := email.New(config)
	if err != nil {
		xpb.Fataler(err)
	}
	err = inbox.FindInvite()
	if err != nil {
		xpb.Fataler(err)
	}
	err = inbox.ExtractInvite()
	if err != nil {
		xpb.Fataler(err)
	}
	// Guest has been invited to the host's GCP console
	// and has accepted the invite. Link the guest's
	// valid billing account
	// 	err = x.LinkBilling()
	// if err != nil {
	//	xpb.Fataler(err)
	// }
}

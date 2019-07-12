package cmd

import (
	"time"

	"github.com/swoldemi/xpb/accounts"
	"github.com/swoldemi/xpb/log"
	"github.com/swoldemi/xpb/browser"
)

// Execute runs the XPB invocation using the provided configuration.
func Execute() {
	config := &xpb.Config{
		HostEmail:            "nickelapi@gmail.com",
		HostPass:             "E5BB55AD2B9C54BB3264FC862513E",
		HostKeyFilePath:      "xpb-host.json",
		HostProjectID:        "nickel-api",
		GuestEmail:           "machserve.io@gmail.com",
		GuestPass:            "convnet1122",
		GuestKeyFilePath:     "xpb-guest.json",
		ChromeDriverPath:     "drivers/chromedriver.exe",
		SeleniumPath:         "drivers/selenium-server-standalone.jar",
		SeleniumRemotePort:   8080,
		SeleniumPollInterval: time.Millisecond * 500,
		SeleniumTimeout:      time.Second * 10,
		Debug:                true,
		Extensions:           xpb.ConfigExtensions{},
	}

	x, err := accounts.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// chrome, err := browser.New(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer func() {
	// 	serr := chrome.SeleniumService.Stop()
	// 	if serr != nil {
	// 		log.Fatal(serr)
	// 	}
	// }()

	// defer func() {
	// 	qerr := chrome.WebDriver.Quit()
	// 	if qerr != nil {
	// 		log.Fatal(qerr)
	// 	}
	// }()

	// err = chrome.LoginHost()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = chrome.InviteGuest()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = chrome.AcceptInvite()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Guest has been invited to the host's GCP console
	// and has accepted the invite. Link the guest's
	// valid billing account
	// 	err = x.LinkBilling()
	// if err != nil {
	//	log.Fatal(err)
	// }
}

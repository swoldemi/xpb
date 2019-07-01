package main

import (
	xpb "github.com/swoldemi/xpb"
	"github.com/swoldemi/xpb/browser"
)

const (
	version = "1.0.0-rc.1"
)

func main() {
	config := &xpb.Config{
		NamedHostEmail:  "nickelapi@gmail.com",
		NamedGuestEmail: "machserve.io@gmail.com",
		HostPass:        "E5BB55AD2B9C54BB3264FC862513E",
		HostProjectID:   "nickel-api",
		Debug:           true,
	}

	b, err := browser.New(config)
	if err != nil {
		xpb.Fataler(err)
	}

	defer func() {
		serr := b.SeleniumService.Stop()
		if serr != nil {
			xpb.Fataler(err)
		}
	}()

	defer func() {
		qerr := b.WebDriver.Quit()
		if qerr != nil {
			xpb.Fataler(err)
		}
	}()

	err = b.LoginHost()
	if err != nil {
		xpb.Fataler(err)
	}

	err = b.InviteGuest()
	if err != nil {
		xpb.Fataler(err)
	}

	// xpb.MustExecute(config)
}

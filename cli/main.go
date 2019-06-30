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
		NamedHostEmail: "nickelapi@gmail.com",
		HostPass:       "E5BB55AD2B9C54BB3264FC862513E",
		HostProjectID:  "nickel-api",
	}
	// xpb.MustExecute(config)

	b, err := browser.New(config)
	if err != nil {
		panic(err)
	}

	defer b.SeleniumService.Stop()
	defer b.WebDriver.Quit()

	err = b.LoginHost()
	if err != nil {
		panic(err)
	}

	err = b.InviteGuest()
	if err != nil {
		panic(err)
	}
}

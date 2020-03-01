package main

import (
	"log"
	"time"

	"github.com/swoldemi/xpb/cmd"
	"github.com/swoldemi/xpb/pkg/config"
)

func main() {
	log.Printf("Version: %s\nGitSHA: %s\n", Version, GitSHA)

	cfg := &config.Config{
		HostEmail:            "<Your Email>",
		HostPass:             "<Your Password>",
		HostKeyFilePath:      "<Path to your keyfile>",
		HostProjectID:        "<Your Project ID>",
		GuestEmail:           "<Your guest's email address>",
		GuestPass:            "<Your guest's password>",
		GuestKeyFilePath:     "<Your guest's keyfile>",
		ChromeDriverPath:     "bin/chromedriver.exe",
		SeleniumPath:         "bin/selenium-server-standalone.jar",
		SeleniumRemotePort:   8080,
		SeleniumPollInterval: time.Millisecond * 500,
		SeleniumTimeout:      time.Second * 10,
		SeleniumDebug:        true,
		Extensions:           config.Extensions{},
	}
	if err := cmd.Execute(cfg); err != nil {
		log.Fatal(err.Error())
	}
}

package browser

import (
	"fmt"
	"os"
	"time"

	"github.com/swoldemi/xpb"
	"github.com/tebeka/selenium"
)

const (
	googleLoginURL = "https://accounts.google.com"
	port           = 8080

	// FIXME: Path relative to the root of the execution directory
	seleniumPath     = "drivers/selenium-server-standalone.jar"
	chromeDriverPath = "drivers/chromedriver.exe"
)

// GCPBrowser encapsulates fields and methods for automated
// interation with the Google Cloud Platform console
type GCPBrowser struct {
	SeleniumService *selenium.Service
	WebDriver       selenium.WebDriver
	Config          *xpb.Config
}

// New connects the to a locally running Selenium instance and
// establishes a remote connection to chrome
func New(config *xpb.Config) (*GCPBrowser, error) {
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath), // Specify the path to ChromeDriver in order to use Chrome.
		selenium.Output(os.Stderr),              // Output debug information to STDERR.
	}

	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		return nil, err
	}

	// Connect to the WebDriver instance running locally.
	c := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(c, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return nil, err
	}

	// Remember to call *selenium.Service.Stop() and selenium.WebDriver.Quit()
	return &GCPBrowser{
		SeleniumService: service,
		WebDriver:       wd,
		Config:          config,
	}, nil
}

// LoginHost log in to the host's Google account
func (g *GCPBrowser) LoginHost() error {
	if err := g.WebDriver.Get(googleLoginURL); err != nil {
		return err
	}

	if err := g.typeHostEmail(); err != nil {
		return err
	}
	if err := g.submitEmail(); err != nil {
		return err
	}
	if err := g.typeHostPassword(); err != nil {
		return err
	}
	if err := g.submitPassword(); err != nil {
		return err
	}

	return nil
}

// InviteGuest invites the guest account, with credits
// to the host's selected project
func (g *GCPBrowser) InviteGuest() error {
	gcpConsoleURL := fmt.Sprintf("https://console.cloud.google.com/home/dashboard?project=%v", g.Config.HostProjectID)
	if err := g.WebDriver.Get(gcpConsoleURL); err != nil {
		return err
	}
	time.Sleep(time.Second * 10)
	return nil
}

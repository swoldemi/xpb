package browser

import (
	"fmt"
	"os"
	"time"

	"github.com/swoldemi/xpb"
	"github.com/tebeka/selenium"
)

// ChromeBrowser encapsulates fields and methods for automated
// interation with the Google Cloud Platform console via ChromeDriver
type ChromeBrowser struct {
	SeleniumService *selenium.Service
	WebDriver       selenium.WebDriver
	Config          *xpb.Config
}

// New starts and connects a locally running Selenium
// instance to the remote ChromeDriver
func New(config *xpb.Config) (*ChromeBrowser, error) {
	selenium.SetDebug(config.Debug)

	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(config.ChromeDriverPath), // Specify the path to ChromeDriver in order to use Chrome.
		selenium.Output(os.Stderr),                     // Output debug information to STDERR.
	}
	service, err := selenium.NewSeleniumService(config.SeleniumPath, config.SeleniumRemotePort, opts...)
	if err != nil {
		return nil, err
	}

	// Connect to the WebDriver instance running locally.
	// Need to set mapped chromeOptions to enable, non-W3C standard functionality (like IsDisplayed)
	// Reference: https://stackoverflow.com/questions/56111529/cannot-call-non-w3c-standard-command-while-in-w3c-mode-seleniumwebdrivererr
	c := selenium.Capabilities{"browserName": "chrome", "chromeOptions": map[string]bool{"w3c": false}}
	wd, err := selenium.NewRemote(c, fmt.Sprintf("http://localhost:%d/wd/hub", config.SeleniumRemotePort))
	if err != nil {
		return nil, err
	}

	// Remember to call *selenium.Service.Stop() and selenium.WebDriver.Quit()
	return &ChromeBrowser{
		SeleniumService: service,
		WebDriver:       wd,
		Config:          config,
	}, nil
}

// Wait is a reused method for blocking until the page has achived a completed ready state
func (c *ChromeBrowser) Wait() error {
	err := c.WebDriver.WaitWithTimeoutAndInterval(ReadyStateCond, c.Config.IntermdiateTimeout, c.Config.PollInterval)
	if err != nil {
		return err
	}
	return nil
}

// LoginHost log in to the host's Google account
func (c *ChromeBrowser) LoginHost() error {
	defer func() {
		err := c.Wait()
		if err != nil {
			xpb.Fataler(err)
		}
	}()

	gcpConsoleURL := fmt.Sprintf("https://console.cloud.google.com/iam-admin/iam?authuser=0&project=%v", c.Config.HostProjectID)
	if err := c.WebDriver.Get(gcpConsoleURL); err != nil {
		return err
	}

	if err := c.typeHostEmail(); err != nil {
		return err
	}
	if err := c.submitEmail(); err != nil {
		return err
	}
	if err := c.typeHostPassword(); err != nil {
		return err
	}
	if err := c.submitPassword(); err != nil {
		return err
	}
	return nil
}

// InviteGuest invites the guest account, with credits,
// to the host's selected project
func (c *ChromeBrowser) InviteGuest() error {
	defer func() {
		err := c.Wait()
		if err != nil {
			xpb.Fataler(err)
		}
	}()

	time.Sleep(time.Second * 10) // FIXME: Just wait for the add button to be visible
	if err := c.clickAdd(); err != nil {
		return err
	}
	if err := c.typeGuestEmail(); err != nil {
		return err
	}

	// A Material option object gets displayed after the email is typed.
	// If this element is visible before beginning role selection,
	// then it will prevent the WebDriver from selecting the role box by
	// intercepting the click. To solve this, remove focus from the email
	// field by clicking on an arbitrary element. In this case, a header
	if err := c.clickHeader(); err != nil {
		return err
	}

	// The owner must have the Owner and Project Billing Manager
	// roles in order to be able to change billing accounts
	// The drawer only displays one role entry field, by default
	if err := c.addFirstRole(); err != nil {
		return err
	}
	if err := c.clickAddAnother(); err != nil {
		return err
	}
	if err := c.addSecondRole(); err != nil {
		return err
	}
	if err := c.submitGuestInvite(); err != nil {
		return err
	}
	return nil
}

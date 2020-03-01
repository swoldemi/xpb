// Package browser enables Selenium driven automation
// of the Google Cloud Platform Console.
// TODO: Browser should encapsulate browser specific packages
// for both Chrome and Firefox (Gecko) clients
// A simple browser interface with `Wait() error`, `LoginHost() error`
// and `InviteGuest() error` methods will suffice
package browser // import "github.com/swoldemi/xpb/pkg/browser"

import (
	"fmt"
	"os"
	"time"

	"github.com/swoldemi/xpb/pkg/config"
	"github.com/swoldemi/xpb/pkg/log"
	"github.com/swoldemi/xpb/pkg/util"
	"github.com/tebeka/selenium"
)

// ChromeBrowser encapsulates fields and methods for automated
// interation with the Google Cloud Platform console via ChromeDriver.
type ChromeBrowser struct {
	SeleniumService *selenium.Service
	WebDriver       selenium.WebDriver
	Config          *config.Config
}

// New starts and connects a locally running Selenium
// instance to the remote ChromeDriver.
func New(config *config.Config) (*ChromeBrowser, error) {
	selenium.SetDebug(config.SeleniumDebug)

	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(config.ChromeDriverPath), // Specify the path to ChromeDriver in order to use Chrome
		selenium.Output(os.Stderr),                     // Output debug information to STDERR
	}
	service, err := selenium.NewSeleniumService(config.SeleniumPath, config.SeleniumRemotePort, opts...)
	if err != nil {
		return nil, err
	}

	// Connect to the WebDriver instance running locally.
	// Need to set mapped chromeOptions to enable, non-W3C standard functionality (like IsDisplayed)
	// Reference: https://stackoverflow.com/questions/56111529/cannot-call-non-w3c-standard-command-while-in-w3c-mode-seleniumwebdrivererr
	c := selenium.Capabilities{"browserName": "chrome", "chromeOptions": map[string]bool{"w3c": false, "start-maximized": true}}
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

// Wait is used for blocking until the page has achieved a completed ready state.
// Currently inconsistent and unreliable.
func (c *ChromeBrowser) Wait() error {
	err := c.WebDriver.WaitWithTimeoutAndInterval(util.ReadyStateCond, c.Config.SeleniumTimeout, c.Config.SeleniumPollInterval)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 5)
	return nil
}

// LoginHost logs in to the host's Google account.
func (c *ChromeBrowser) LoginHost() error {
	defer func() {
		err := c.Wait()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	gcpConsoleURL := fmt.Sprintf("https://console.cloud.google.com/iam-admin/iam?authuser=0&project=%v", c.Config.HostProjectID)
	if err := c.WebDriver.Get(gcpConsoleURL); err != nil {
		return err
	}
	if err := c.TypeLoginEmail(c.Config.HostEmail); err != nil {
		return err
	}
	if err := c.SubmitEmail(); err != nil {
		return err
	}
	if err := c.TypeLoginPassword(c.Config.HostPass); err != nil {
		return err
	}
	if err := c.SubmitPassword(); err != nil {
		return err
	}
	return nil
}

// InviteGuest invites the guest account, with credits,
// to the host's selected project.
func (c *ChromeBrowser) InviteGuest() error {
	defer func() {
		err := c.Wait()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	time.Sleep(time.Second * 10) // FIXME: Just wait for the add button to be visible
	if err := c.ClickAdd(); err != nil {
		return err
	}
	if err := c.TypeGuestEmail(); err != nil {
		return err
	}

	// A Material option object gets displayed after the email is typed.
	// If this element is visible before beginning role selection,
	// then it will prevent the WebDriver from selecting the role box by
	// intercepting the click. To solve this, remove focus from the email
	// field by clicking on an arbitrary element. In this case, a header
	if err := c.ClickHeader(); err != nil {
		return err
	}

	// The owner must have the Owner and Project Billing Manager
	// roles in order to be able to change billing accounts
	// The drawer only displays one role entry field, by default
	if err := c.AddFirstRole(); err != nil {
		return err
	}
	if err := c.ClickAddAnother(); err != nil {
		return err
	}
	if err := c.AddSecondRole(); err != nil {
		return err
	}
	if err := c.SubmitGuestInvite(); err != nil {
		return err
	}
	return nil
}

// AcceptInvite logs in to the guest's Google account and accepts the host's invite.
func (c *ChromeBrowser) AcceptInvite() error {
	defer func() {
		err := c.Wait()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	// Assume that the format of the login url is always the same
	// TODO: What happens if this is invoked too quickly?
	inviteURL := fmt.Sprintf(
		"https://console.cloud.google.com/invitation?project=%v&account=%v&memberEmail=%v",
		c.Config.HostProjectID,
		c.Config.GuestEmail,
		c.Config.GuestEmail,
	)
	if err := c.WebDriver.Get(inviteURL); err != nil {
		return err
	}
	if err := c.TypeLoginEmail(c.Config.GuestEmail); err != nil {
		return err
	}
	if err := c.SubmitEmail(); err != nil {
		return err
	}
	if err := c.TypeLoginPassword(c.Config.GuestPass); err != nil {
		return err
	}
	if err := c.SubmitPassword(); err != nil {
		return err
	}

	// Desperate wait to wait for modal bootstrapping...
	if err := c.Wait(); err != nil {
		return err
	}
	return c.ClickInvite()
}

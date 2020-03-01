// Package config provides configuration primitives and extensions.
package config // import "github.com/swoldemi/xpb/pkg/config"

import "time"

// Extensions defines extended config behavior for the CLI.
type Extensions struct {
	UseYAML  bool
	YAMLPath string
}

// Config encapsulates parameters and arguments for XPB execution.
type Config struct {
	// Required. The email address of the GCP account with a project
	// that you want to keep running
	HostEmail string

	// Required. The password for the email address used to login
	// to the GCP account
	HostPass string

	// Required. The path to the keyfile for a service account made in
	// the host's account
	HostKeyFilePath string

	// Required. The ID of the project you want to refresh the billing account for.
	HostProjectID string

	// Required. The email address of the secondary GCP account made
	// for a new trial account
	GuestEmail string

	// Required. The password for the email address used to login to
	// the GCP account
	GuestPass string

	// Required. The path to the keyfile for a service account made in
	// the guest's account
	GuestKeyFilePath string

	// Required. Path to the ChromeDriver WebDriver for Chrome
	// Download from here: https://chromedriver.chromium.org/downloads
	ChromeDriverPath string

	// Required. The path to the sandalone selenium server binary. One is
	// provided in the `bin` directory as "selenium-server-standalone.jar"
	SeleniumPath string

	// Required. The port that the server is configured to listen on
	SeleniumRemotePort int

	// Required. The poll duration for the server
	SeleniumPollInterval time.Duration

	// Required. The timeout duration for the server
	SeleniumTimeout time.Duration

	// Required. Should debug be enabled for the server?
	SeleniumDebug bool

	// Optional. Any extensions; mainly for configuration once XPB is
	// shipped as a portable binary.
	Extensions Extensions
}

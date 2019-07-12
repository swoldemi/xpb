package config

import "time"

// Extensions defines extended config behavior for the CLI.
type Extensions struct {
	UseYAML  bool
	YAMLPath string
}

// Config encapsulates parameters and arguments for XPB execution.
type Config struct {
	HostEmail            string
	HostPass             string
	HostKeyFilePath      string
	HostProjectID        string
	GuestEmail           string
	GuestPass            string
	GuestKeyFilePath     string
	ChromeDriverPath     string
	SeleniumPath         string
	SeleniumRemotePort   int
	SeleniumPollInterval time.Duration
	SeleniumTimeout      time.Duration
	SeleniumDebug                bool
	Extensions           ConfigExtensions
}

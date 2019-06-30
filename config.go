package xpb

import "time"

// ConfigExtensions defines extended config behavior for the CLI
type ConfigExtensions struct {
	UseYAML  bool
	YAMLPath string
}

// Config encapsulates parameters and arguments for XPB execution
type Config struct {
	AliasOverride      string
	NamedHostEmail     string
	NamedGuestEmail    string
	IntermdiateTimeout *time.Duration
	Extensions         *ConfigExtensions
	HostPass           string
	HostProjectID      string
	HostKeyFilePath    string
	GuestPass          string
}

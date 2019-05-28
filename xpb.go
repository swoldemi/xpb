package xpb

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	l *logrus.Logger
)

func init() {
	l = logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	l.SetOutput(os.Stdout)

	// log the trace severity
	// TODO: Make this a config flag
	l.SetLevel(logrus.TraceLevel)

	// Display caller in log trace
	l.SetReportCaller(false)
}

// Fataler calls logrus.Logger.Fatal on a non-nil error
func Fataler(err error) {
	if err != nil {
		l.Fatal(err)
	}
}

// MustExecute authenticates an XPB client and calls Execute
func MustExecute(config *Config) {
	client, err := New(l, config)
	Fataler(err)

	l.Info("Beginning XPB execution...")
	xerr := client.Execute()
	Fataler(xerr)
}

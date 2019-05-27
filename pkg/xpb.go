package xpb

import (
	"github.com/sirupsen/logrus"
)

// MustExecute authenticates an XPB client and calls Execute
func MustExecute(l *logrus.Entry, config *Config) {
	client, err := New(l, config)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("Beginning XPB execution...")
	xerr := client.Execute()
	if xerr != nil {
		l.Fatal(xerr)
	}
}

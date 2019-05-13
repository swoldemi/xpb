package xpb

import (
	"github.com/sirupsen/logrus"
)

// MustExecute authenticates an XPB client and calls Execute
func MustExecute(l *logrus.Logger) {
	client, err := New(l)
	if err != nil {
		l.Error(err)
		panic(err)
	}

	xerr := client.Execute()
	if xerr != nil {
		l.Error(xerr)
		panic(xerr)
	}
}

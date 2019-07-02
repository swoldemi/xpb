package xpb

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var (
	// ErrGuestNoBilling - the guest account List call returns empty BillingAccounts
	ErrGuestNoBilling = errors.New("xpb: guest has no valid billing accounts")

	// ErrAccountsIdentical - the host account's project is already associated with the guest account's billing account
	ErrAccountsIdentical = errors.New("xpb: host account is already set to guest's billing account")
)

// NewLogger creates a component divided logger
func NewLogger(name string) *logrus.Entry {
	l := logrus.New()

	// Log as JSON instead of the default ASCII formatter
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	l.SetOutput(os.Stdout)

	// log the trace severity
	// TODO: Make this a config flag
	l.SetLevel(logrus.TraceLevel)

	// Do not display caller in log trace
	l.SetReportCaller(false)

	return l.WithField("component", name)
}

// Fataler does a cheap fatal call (print error and exit with code 1)
func Fataler(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// GetADCEmail retrieve the email address of the user that is currently
// authenticated using Application Default Credentials via the
// local gcloud CLI configuration
func GetADCEmail() []byte {
	cmd := exec.Command("gcloud", "config", "get-value", "core/account")
	out, err := cmd.CombinedOutput()
	Fataler(err)
	return out
}

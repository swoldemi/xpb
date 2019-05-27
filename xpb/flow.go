package xpb

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	iam "google.golang.org/api/iam/v1"
)

var (
	// ErrGuestBillingDisabled - the guest account List call returns empty BillingAccounts
	ErrGuestBillingDisabled = errors.New("xpb: guest account's billing is disabled")

	// ErrAccountsIdentical - the host account's project is already associated with the guest account's billing account
	ErrAccountsIdentical = errors.New("xpb: host account is already set to guest's billing account")

	hostChan  = make(chan struct{}, 1)
	guestChan = make(chan struct{}, 1)
)

// Host encapsulates API's and services for authenticating and
// interacting with the host account's resources.
// This account owns the project that is in need of trial credits.
type Host struct {
	BillingService *cloudbilling.APIService
	IamService     *iam.Service
	ProjectID      string
	SvcEmail       string
}

// Guest encapsulates API's and services for authenticating and
// interacting with the guest account's resources.
// This account owns the trial billing account that has credits
// for the host account's.
type Guest struct {
	BillingService *cloudbilling.APIService
	IamService     *iam.Service
	ProjectID      string
	SvcEmail       string
}

// Flow encapsulates attributes necessary for
// invoking the cross-project billing flow.
// All clients are given minimum required OAuth scopes.
type Flow struct {
	Host  *Host
	Guest *Guest
	l     *logrus.Logger
}

// New creates a new XPB Client
func New(l *logrus.Logger, config *Config) (*Flow, error) {
	l.Debug("Creating new Flow instance...")

	host := &Host{}
	guest := &Guest{}
	go AuthenticateHost(config, l, host, hostChan)
	go AuthenticateGuest(config, l, guest, guestChan)
	<-hostChan
	<-guestChan

	return &Flow{
		Host:  host,
		Guest: guest,
		l:     l,
	}, nil
}

// Execute begins the XPB sequence
func (f *Flow) Execute() error {
	ctx := context.Background()

	f.l.Infof("Host's project is set to: %v", f.Host.ProjectID)
	f.l.Infof("Host's service account is: %v", f.Host.SvcEmail)
	f.l.Infof("Guest's project is set to: %v", f.Guest.ProjectID)
	f.l.Infof("Guest's service account is: %v", f.Guest.SvcEmail)

	// Get host account's billing information for the current project
	hostInfo, err := f.Host.BillingService.Projects.GetBillingInfo(projectsPrefix + f.Host.ProjectID).Context(ctx).Do()
	if err != nil {
		f.l.Error(err)
		return err
	}
	f.l.WithFields(logrus.Fields{
		"user":                 "Host",
		"Project ID":           hostInfo.ProjectId,
		"Billing Account Name": hostInfo.BillingAccountName,
		"Enabled":              hostInfo.BillingEnabled,
	}).Infof("Retrieved Host's billing information")
	if !hostInfo.BillingEnabled {
		f.l.Info(yikes)
	}

	f.l.Infof(projectsPrefix + f.Guest.ProjectID)
	// Get host account's billing information for the current project
	guestInfo, err := f.Guest.BillingService.Projects.GetBillingInfo(projectsPrefix + f.Guest.ProjectID).Context(ctx).Do()
	if err != nil {
		f.l.Error(err)
		return err
	}
	f.l.WithFields(logrus.Fields{
		"user":                 "Guest",
		"Project ID":           guestInfo.ProjectId,
		"Billing Account Name": guestInfo.BillingAccountName,
		"Enabled":              guestInfo.BillingEnabled,
	}).Infof("Retrieved Host's billing information")
	if !guestInfo.BillingEnabled {
		return ErrGuestBillingDisabled
	}

	// Make sure the two accounts aren't already the same
	if guestInfo.BillingAccountName == hostInfo.BillingAccountName {
		return ErrAccountsIdentical
	}

	// Set the host account's project billing account to the guest's billing account

	return nil
}

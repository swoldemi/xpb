package xpb

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	iam "google.golang.org/api/iam/v1"
)

var hostChan = make(chan struct{}, 1)
var guestChan = make(chan struct{}, 1)

var (
	// ErrNoGuestBilling - the guest account List call returns empty BillingAccounts
	ErrNoGuestBilling = errors.New("xpb: guest account does not have any open billing accounts")

	// ErrAccountsIdentical - the host account's project is already associated with the guest account's billing account
	ErrAccountsIdentical = errors.New("xpb: host account is already set to guest's billing account")
)

// Host encapsulates API's and services for authenticating and
// interacting with the host account's resources.
// This account owns the project that is in need of trial credits.
type Host struct {
	BillingService *cloudbilling.APIService
	IamService     *iam.Service
	ProjectID      string
}

// Guest encapsulates API's and services for authenticating and
// interacting with the guest account's resources.
// This account owns the trial billing account that has credits
// for the host account's.
type Guest struct {
	BillingService *cloudbilling.APIService
	IamService     *iam.Service
}

// Flow encapsulates attributes necessary for
// invoking the cross-project billing flow.
// All clients are given minimum required OAuth scopes.
type Flow struct {
	Host  *Host
	Guest *Guest
	l     *logrus.Entry
}

// New creates a new XPB Client
func New(l *logrus.Entry, config *Config) (*Flow, error) {
	l.Debug("Creating new XPBClient [GCP IAM, GCP Billing]...")

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

	// Get host account's billing information for the current project
	project := fmt.Sprintf("projects/%v", f.Host.ProjectID)
	hostInfo, err := f.Host.BillingService.Projects.GetBillingInfo(project).Context(ctx).Do()
	if err != nil {
		f.l.Error(err)
		return err
	}
	f.l.WithFields(logrus.Fields{
		"user":                 "Host",
		"ProjectID":            hostInfo.ProjectId,
		"Billing Account Name": hostInfo.BillingAccountName,
		"Enabled":              hostInfo.BillingEnabled,
	}).Infof("Retrieved Host's billing information")
	f.l.Infof("%+v", hostInfo)
	if !hostInfo.BillingEnabled {
		f.l.Info(yikes)
	}

	// Get the guest account's available billing accounts
	// Given, new billing account is enabled
	var accounts []*cloudbilling.BillingAccount

	call := f.Guest.BillingService.BillingAccounts.List()
	perr := call.Pages(ctx, func(page *cloudbilling.ListBillingAccountsResponse) error {
		for _, account := range page.BillingAccounts {
			f.l.Infof("%+v", account)
			accounts = append(accounts, account)
		}
		return nil
	})

	if perr != nil {
		f.l.Error(perr)
		return perr
	}

	if len(accounts) == 0 {
		return ErrNoGuestBilling
	}

	// TODO: Allow user to select preferred billing account
	guestInfo := accounts[0]

	// Make sure the two accounts aren't already the same
	if guestInfo.Name == hostInfo.BillingAccountName {
		return ErrAccountsIdentical
	}

	// Set the host account's project billing account to the guest's billing account

	return nil
}

package xpb

import (
	"context"

	"github.com/sirupsen/logrus"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
)

var (
	hostChan  = make(chan struct{}, 1)
	guestChan = make(chan struct{}, 1)
)

// Host encapsulates API's and services for authenticating and
// interacting with the host account's resources.
// This account owns the project that is in need of trial credits.
type Host struct {
	BillingService     *cloudbilling.APIService
	IamService         *iam.Service
	ResourceMgrService *cloudresourcemanager.Service
	ProjectID          string
	SvcEmail           string
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

// New creates a new Flow Client
func New(config *Config) (*Flow, error) {
	l := NewLogger("Flow")
	l.Debug("Creating new Flow instance...")

	host := &Host{}
	guest := &Guest{}
	go AuthenticateHost(config, l, host, hostChan)
	go AuthenticateGuest(l, guest, guestChan)
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

	// Get host account's billing information for the current project
	hostInfo, err := f.getHostBilling(ctx)
	if err != nil {
		f.l.Error(err)
		return err
	}

	guestInfo, err := f.getGuestBilling(ctx)
	if err != nil {
		f.l.Error(err)
		return err
	}

	// Make sure the two accounts aren't already the same
	if guestInfo.Name == hostInfo.Name {
		return ErrAccountsIdentical
	}

	// Per `package browser`, guest has permission to access the host's project.
	// Set the host account's project billing account to the guest's billing account
	association := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: guestInfo.Name,
	}

	updateResp, err := f.Guest.BillingService.Projects.UpdateBillingInfo("projects/"+f.Host.ProjectID, association).Context(ctx).Do()
	if updateResp != nil {
		f.l.Infof("%+v", updateResp)

	} else {
		Fataler(err)
	}

	f.l.Infof("%s's billing account changed from %s to %s. Exchange complete.", f.Host.ProjectID, hostInfo.BillingAccountName, guestInfo.Name)
	return nil
}

func (f *Flow) getHostBilling(ctx context.Context) (*cloudbilling.ProjectBillingInfo, error) {
	hostInfo, err := f.Host.BillingService.Projects.GetBillingInfo("projects/" + f.Host.ProjectID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	f.l.WithFields(logrus.Fields{
		"user":                 "Host",
		"Project ID":           hostInfo.ProjectId,
		"Billing Account Name": hostInfo.BillingAccountName,
		"Enabled":              hostInfo.BillingEnabled,
		"Resource Name":        hostInfo.Name,
	}).Infof("Retrieved Host's billing information")

	if !hostInfo.BillingEnabled {
		f.l.Info(yikes)
	}

	return hostInfo, err
}

func (f *Flow) getGuestBilling(ctx context.Context) (*cloudbilling.BillingAccount, error) {
	// Get guest account's billing information for the current project
	// Because we don't need any of the guest's ProjectIDs,
	// we need to list the guest's billing accounts
	guestAccounts, err := f.Guest.BillingService.BillingAccounts.List().Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var validAccounts []*cloudbilling.BillingAccount
	for _, account := range guestAccounts.BillingAccounts {
		if account.Open {
			validAccounts = append(validAccounts, account)
		}
	}
	if len(validAccounts) == 0 {
		return nil, ErrGuestNoBilling
	}

	// Pick the first account for now
	guestInfo := validAccounts[0]
	f.l.WithFields(logrus.Fields{
		"user":                   "Guest",
		"Billing Display Name":   guestInfo.DisplayName,
		"Master Billing Account": guestInfo.MasterBillingAccount,
		"Enabled":                guestInfo.Open,
		"Billing Account Name":   guestInfo.Name,
	}).Infof("Retrieved Guest's's billing information")

	return guestInfo, nil
}

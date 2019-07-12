package accounts

import (
	"context"

	"github.com/swoldemi/xpb/log"
	"github.com/swoldemi/xpb/util"
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
	ProjectID string
}

// Guest encapsulates API's and services for authenticating and
// interacting with the guest account's resources.
// This account owns the trial billing account that has credits
// for the host account's.
type Guest struct {
	BillingService *cloudbilling.APIService
	IamService     *iam.Service
}

// Accounts encapsulates attributes and methods 
// necessary for verifying the state of the host
// and guest accounts.
type Accounts struct {
	Host  *Host
	Guest *Guest
}

// New creates a new Accounts instance.
func New(config *Config) (*Accounts, error) {
	log.Debug("Creating new Accounts instance...")

	host := &Host{}
	guest := &Guest{}
	go AuthenticateHost(config, host, hostChan)
	go AuthenticateGuest(config, guest, guestChan)
	<-hostChan
	<-guestChan

	return &Accounts{
		Host:  host,
		Guest: guest,
	}, nil
}

// Verify begins the XPB sequence by
// checking that the guest and host are valid
// accounts and retrives the guest's billing account
// which will be linked to the host's project.
func (f *Flow) Verify() (string, error) {
	ctx := context.Background()

	log.Info("Host's project is set to: %v", f.Host.ProjectID)
	log.Info("Host's service account is: %v", f.Host.SvcEmail)

	hostInfo, err := f.getHostBilling(ctx)
	if err != nil {
		log.Error(err)
		return "", err
	}

	guestInfo, err := f.getGuestBilling(ctx)
	if err != nil {
		log.Error(err)
		return "", err
	}

	// Make sure the two accounts aren't already the same
	if guestInfo.Name == hostInfo.Name {
		return "", util.ErrAccountsIdentical
	}

	// State is good enough, return the guest's billing account
	return guestInfo.MasterBillingAccount, nil
}


// Link performs the link of the guest's billing account to the host.
// This assumes the guest has accepted the host's invite, per 
// `package browser`'s automation.
func (f *Flow) Link() error {
	association := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: guestInfo.Name,
	}

	updateResp, err := f.Guest.BillingService.Projects.UpdateBillingInfo("projects/"+f.Host.ProjectID, association).Context(ctx).Do()
	if updateResp != nil {
		log.Trace("%+v", updateResp)
	} else {
		log.Fatal(err)
	}

	log.Trace("%s's billing account changed from %s to %s. Exchange complete.", f.Host.ProjectID, hostInfo.BillingAccountName, guestInfo.Name)
}

 func (f *Flow) getHostBilling(ctx context.Context) (*cloudbilling.ProjectBillingInfo, error) {
	hostInfo, err := f.Host.BillingService.Projects.GetBillingInfo("projects/" + f.Host.ProjectID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	log.Trace("%+v", map[string]string{
		"user":                 "Host",
		"Project ID":           hostInfo.ProjectId,
		"Billing Account Name": hostInfo.BillingAccountName,
		"Enabled":              hostInfo.BillingEnabled,
		"Resource Name":        hostInfo.Name,
	})
	log.Trace("Retrieved Host's billing information")

	if !hostInfo.BillingEnabled {
		log.Warning("This billing account isn't enabled anymore!")
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
	log.Info("%+v", map[string]string{
		"user":                   "Guest",
		"Billing Display Name":   guestInfo.DisplayName,
		"Master Billing Account": guestInfo.MasterBillingAccount,
		"Enabled":                guestInfo.Open,
		"Billing Account Name":   guestInfo.Name,
	})
	log.Trace("Retrieved Guest's billing information")

	return guestInfo, nil
}

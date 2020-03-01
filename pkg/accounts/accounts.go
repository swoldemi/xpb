package accounts

import (
	"context"

	"github.com/swoldemi/xpb/pkg/config"
	"github.com/swoldemi/xpb/pkg/log"
	"github.com/swoldemi/xpb/pkg/util"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
)

// Host encapsulates API's and services for authenticating and
// interacting with the host account's resources.
// This account owns the project that is in need of trial credits.
type Host struct {
	BillingService *cloudbilling.APIService
	ProjectID      string
	SvcEmail       string
}

// Guest encapsulates API's and services for authenticating and
// interacting with the guest account's resources.
// This account owns the trial billing account that has credits
// for the host account's.
type Guest struct {
	BillingService *cloudbilling.APIService
	ProjectID      string
	SvcEmail       string
}

// Accounts encapsulates attributes and methods
// necessary for verifying the state of the host
// and guest accounts.
type Accounts struct {
	Host
	Guest
}

// New creates a new Accounts instance.
func New(c *config.Config) (*Accounts, error) {
	log.Debug("Creating new Accounts instance...")

	host, err := AuthenticateHost(c)
	if err != nil {
		return nil, err
	}

	guest, err := AuthenticateGuest(c)
	if err != nil {
		return nil, err
	}

	return &Accounts{
		Host:  host,
		Guest: guest,
	}, nil
}

// Verify begins the XPB sequence by
// checking that the guest and host are valid
// accounts and retrives the guest's billing account
// which will be linked to the host's project.
func (a *Accounts) Verify() (string, error) {
	log.Info("Host's project is set to: %v", a.Host.ProjectID)
	log.Info("Host's service account is: %v", a.Host.SvcEmail)
	log.Info("Guest's project is set to: %v", a.Guest.ProjectID)
	log.Info("Guest's service account is: %v", a.Guest.SvcEmail)

	hostInfo, err := a.GetHostBilling()
	if err != nil {
		return "", err
	}

	guestInfo, err := a.GetGuestBilling()
	if err != nil {
		return "", err
	}

	// Make sure the two accounts aren't already the same
	if guestInfo.Name == hostInfo.Name {
		return "", util.ErrAccountsIdentical
	}

	// State is good enough, return the guest's billing account
	return guestInfo.BillingAccountName, nil
}

// Link performs the link of the guest's billing account to the host.
// This assumes the guest has accepted the host's invite, per
// `package browser`'s automation.
func (a *Accounts) Link(billingAccount string) error {
	ctx := context.Background()
	association := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: billingAccount,
	}
	updateResp, err := a.Guest.BillingService.Projects.UpdateBillingInfo("projects/"+a.Host.ProjectID, association).Context(ctx).Do()
	if updateResp != nil {
		return err
	}
	log.Trace("%s's billing account changed to %s. Exchange complete: %+v", a.Host.ProjectID, billingAccount, updateResp)
	return nil
}

// GetHostBilling retrieves the host's billing account information.
func (a *Accounts) GetHostBilling() (*cloudbilling.ProjectBillingInfo, error) {
	ctx := context.Background()
	hostInfo, err := a.Host.BillingService.Projects.GetBillingInfo("projects/" + a.Host.ProjectID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	log.Trace("%+v", map[string]string{
		"user":                 "Host",
		"Project ID":           hostInfo.ProjectId,
		"Billing Account Name": hostInfo.BillingAccountName,
		"Resource Name":        hostInfo.Name,
	})
	log.Trace("Retrieved Host's project billing information")

	if !hostInfo.BillingEnabled {
		log.Warning("Host's billing account isn't enabled anymore!")
	} else {
		log.Warning("Host's billing account is enabled! It is recommended to only switch accounts when you have low credits.")
	}
	return hostInfo, nil
}

// GetGuestBilling retrieves the guest's billing account information.
func (a *Accounts) GetGuestBilling() (*cloudbilling.ProjectBillingInfo, error) {
	ctx := context.Background()
	guestInfo, err := a.Guest.BillingService.Projects.GetBillingInfo("projects/" + a.Guest.ProjectID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	log.Trace("%+v", map[string]string{
		"user":                 "Guest",
		"Project ID":           guestInfo.ProjectId,
		"Billing Account Name": guestInfo.BillingAccountName,
		"Resource Name":        guestInfo.Name,
	})
	log.Trace("Retrieved Guest's project billing information")

	if !guestInfo.BillingEnabled {
		return nil, util.ErrGuestNoBilling
	}

	return guestInfo, nil
}

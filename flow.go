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
		"Resource Name":        hostInfo.Name,
	}).Infof("Retrieved Host's billing information")
	if !hostInfo.BillingEnabled {
		f.l.Info(yikes)
	}

	// Get guest account's billing information for the current project
	// Because we don't have any ProjectID, we need to list the guest's billing accounts
	guestAccounts, err := f.Guest.BillingService.BillingAccounts.List().Context(ctx).Do()
	if err != nil {
		f.l.Error(err)
		return err
	}
	var validAccounts []*cloudbilling.BillingAccount
	for _, account := range guestAccounts.BillingAccounts {
		if account.Open {
			validAccounts = append(validAccounts, account)
		}
	}

	if len(validAccounts) == 0 {
		return ErrGuestNoBilling
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

	// Make sure the two accounts aren't already the same
	if guestInfo.Name == hostInfo.Name {
		return ErrAccountsIdentical
	}

	// Allow the guest to change the host's billing account by creating an IAM
	// binding on the host's account.
	// NOTE: This must be done using cloudresourcemanager/v1.
	// See why here: https://github.com/googleapis/google-api-go-client/blob/master/iam/v1/iam-gen.go#L5848-L5875
	policies, err := f.Host.ResourceMgrService.Projects.GetIamPolicy(
		f.Host.ProjectID,
		&cloudresourcemanager.GetIamPolicyRequest{},
	).Do()
	Fataler(err)

	// Grant the guest's service account owner and billing.Projectmanager roles.
	// Note: these are both necessary for a service account to change billing state.
	// See: https://support.google.com/cloud/answer/7283646?hl=en
	roles := []string{"roles/billing.admin"}
	member := "user:backtracksimon@gmail.com"
	f.l.Info(member)
	for _, role := range roles {
		policies.Bindings = append(policies.Bindings, &cloudresourcemanager.Binding{
			Role:    role,
			Members: []string{member},
		})
	}
	rb := &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policies,
	}
	setResp, err := f.Host.ResourceMgrService.Projects.SetIamPolicy(f.Host.ProjectID, rb).Context(ctx).Do()
	if setResp != nil {
		f.l.Infof("%+v", setResp.ServerResponse.Header["Date"])
	} else {
		Fataler(err)
	}

	// // Guest now has permission to access the host's project.
	// // Set the host account's project billing account to the guest's billing account
	// association := &cloudbilling.ProjectBillingInfo{
	// 	BillingAccountName: guestInfo.BillingAccountName,
	// }

	// updateResp, err := f.Guest.BillingService.Projects.UpdateBillingInfo(projectsPrefix+f.Host.ProjectID, association).Context(ctx).Do()
	// if updateResp != nil {
	// 	f.l.Infof("%+v", updateResp)

	// } else {
	// 	Fataler(err)
	// }

	// f.l.Infof("%s's billing account changed from %s to %s. Exchange complete.", f.Host.ProjectID, hostInfo.BillingAccountName, guestInfo.BillingAccountName)
	return nil
}

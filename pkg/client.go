package xpb

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	gmail "google.golang.org/api/gmail/v1"
	iam "google.golang.org/api/iam/v1"
)

// Client encapsulates attributes and http.Client's necessary for
// invoking cross-project billing. All clients are given minimum required
// OAuth scopes.
type Client struct {
	GmailService   *gmail.Service
	BillingService *cloudbilling.APIService
	IamService     *iam.Service
	l              *logrus.Logger
}

// New creates a new XPB Client
func New(l *logrus.Logger) (*Client, error) {
	l.Debug("Creating new XPBClient [Gmail, GCP IAM, GCP Billing]...")
	ctx := context.Background()

	l.Debug("Attempting to authenticate using application default credentials...")
	// NOTE: package cloudbilling and package iam both export CloudPlatformScope
	client, err := google.DefaultClient(ctx, gmail.GmailReadonlyScope, iam.CloudPlatformScope)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	l.Debug("Instantiating Gmail service client...")
	gmailService, err := gmail.New(client)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	l.Debug("Instantiating billing service client...")
	billingService, err := cloudbilling.New(client)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	l.Debug("Instantiating IAM service client...")
	iamService, err := iam.New(client)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return &Client{
		GmailService:   gmailService,
		BillingService: billingService,
		IamService:     iamService,
		l:              l,
	}, nil
}

// Execute begins the XPB sequence
func (c *Client) Execute() error {
	c.l.Info("Authentication succeeded. Beginning XPB execution...")
	b, err := c.BillingService.Projects.GetBillingInfo("projects/nickel-api").Do()
	if err != nil {
		return err
	}

	c.l.Infof("Project is set to: %v", b.ProjectId)

	// result, err := c.RolesService.
	// Invite the `guest` account to the current project
	reqCall := c.BillingService.BillingAccounts.List()
	result, err := reqCall.Do()
	if err != nil {
		c.l.Error(err)
		return err
	}
	for _, account := range result.BillingAccounts {
		fmt.Println(account.Name, account.Open)
	}
	return nil
}

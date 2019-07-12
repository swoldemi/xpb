package accounts

import (
	"context"

	"github.com/swoldemi/xpb/pkg/log"
	"github.com/swoldemi/xpb/pkg/util"
	"github.com/swoldemi/xpb/pkg/config"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

// AuthenticateHost fulfills xpb.Host
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateHost(c *config.Config, host *Host, wait chan struct{}) {
	log.Trace("Attempting to authenticate host using provided service account...")
	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithCredentialsFile(c.HostKeyFilePath),
	}

	log.Trace("Instantiating IAM service client for host account...")
	iamsvc, err := iam.NewService(ctx, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Trace("Instantiating Cloud Billing service client for host account...")
	opts = append(opts, option.WithScopes("https://www.googleapis.com/auth/cloud-billing"))
	billingsvc, err := cloudbilling.NewService(ctx, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Trace("Instantiating Cloud Resource Manager service client for host account...")
	opts = append(opts, option.WithScopes(cloudresourcemanager.CloudPlatformScope))
	resourcemgrsvc, err := cloudresourcemanager.NewService(ctx, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Trace("Storing host's clients and service account information...")
	keyFile := util.ReadKeyFile(c.HostKeyFilePath)
	host.BillingService = billingsvc
	host.IamService = iamsvc
	host.ResourceMgrService = resourcemgrsvc
	host.ProjectID = keyFile["project_id"].(string)
	host.SvcEmail = keyFile["client_email"].(string)
	close(wait)
}

// AuthenticateGuest fulfills accounts.Guest
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateGuest(c *config.Config, guest *Guest, wait chan struct{}) {
	log.Trace("Attempting to authenticate guest using provided service account...")
	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithCredentialsFile(c.GuestKeyFilePath),
	}

	log.Trace("Instantiating IAM service client for guest account...")
	iamsvc, err := iam.NewService(ctx, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Trace("Instantiating Cloud Billing service client for guest account...")
	opts = append(opts, option.WithScopes("https://www.googleapis.com/auth/cloud-billing"))
	billingsvc, err := cloudbilling.NewService(ctx, opts...)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Trace("Storing guest's clients...")
	guest.BillingService = billingsvc
	guest.IamService = iamsvc
	close(wait)
}

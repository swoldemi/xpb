package accounts

import (
	"context"

	"github.com/swoldemi/xpb/pkg/config"
	"github.com/swoldemi/xpb/pkg/log"
	"github.com/swoldemi/xpb/pkg/util"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

// AuthenticateHost fulfills xpb.Host
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateHost(c *config.Config) (Host, error) {
	log.Trace("Attempting to authenticate host using provided service account...")
	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithCredentialsFile(c.HostKeyFilePath),
	}

	log.Trace("Instantiating Cloud Billing service client for host account...")
	billingsvc, err := cloudbilling.NewService(ctx, opts...)
	if err != nil {
		return Host{}, err
	}

	log.Trace("Storing host's clients and service account information...")
	keyFile := util.ReadKeyFile(c.HostKeyFilePath)
	return Host{
		BillingService: billingsvc,
		ProjectID:      keyFile["project_id"].(string),
		SvcEmail:       keyFile["client_email"].(string),
	}, nil
}

// AuthenticateGuest fulfills accounts.Guest
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateGuest(c *config.Config) (Guest, error) {
	log.Trace("Attempting to authenticate guest using provided service account...")
	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithCredentialsFile(c.GuestKeyFilePath),
	}

	log.Trace("Instantiating Cloud Billing service client for guest account...")
	billingsvc, err := cloudbilling.NewService(ctx, opts...)
	if err != nil {
		return Guest{}, err
	}

	log.Trace("Storing guest's clients and service account information...")
	keyFile := util.ReadKeyFile(c.GuestKeyFilePath)

	return Guest{
		BillingService: billingsvc,
		ProjectID:      keyFile["project_id"].(string),
		SvcEmail:       keyFile["client_email"].(string),
	}, nil
}

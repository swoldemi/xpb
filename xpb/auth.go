package xpb

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func readKeyFile(filepath string, l *logrus.Logger) map[string]interface{} {
	l.Debug("Storing guest's service account information...")

	// Read key file only to extract project_id and client_email
	data, err := ioutil.ReadFile(filepath)
	Fataler(err)

	var svcAccount map[string]interface{}
	jerr := json.Unmarshal(data, &svcAccount)
	Fataler(jerr)

	return svcAccount
}

// AuthenticateHost fulfills xpb.Host
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateHost(config *Config, l *logrus.Logger, host *Host, c chan struct{}) {
	l.Debug("Attempting to authenticate host using provided service account...")
	ctx := context.Background()

	opts := []option.ClientOption{option.WithCredentialsFile(config.HostKeyFilePath)}
	l.Debug("Instantiating billing service client for host account...")
	billing, err := cloudbilling.NewService(ctx, opts...)
	Fataler(err)

	l.Debug("Instantiating IAM service client for host account...")
	iam, err := iam.NewService(ctx, opts...)
	Fataler(err)

	keyFile := readKeyFile(config.HostKeyFilePath, l)
	host.BillingService = billing
	host.IamService = iam
	host.ProjectID = keyFile["project_id"].(string)
	host.SvcEmail = keyFile["client_email"].(string)
	close(c)
}

// AuthenticateGuest fulfills xpb.Guest
// with services authenticated and instantiated using
// the provided service account key file.
// FIXME: Do logrus.Logger.Fatal calls even allow the close(chan) calls to run?
func AuthenticateGuest(config *Config, l *logrus.Logger, guest *Guest, c chan struct{}) {
	l.Debug("Attempting to authenticate guest using provided service account...")
	ctx := context.Background()

	opts := []option.ClientOption{option.WithCredentialsFile(config.GuestKeyFilePath)}
	l.Debug("Instantiating billing service client for guest account...")
	billing, err := cloudbilling.NewService(ctx, opts...)
	Fataler(err)

	l.Debug("Instantiating IAM service client for guest account...")
	iam, err := iam.NewService(ctx, opts...)
	Fataler(err)

	keyFile := readKeyFile(config.GuestKeyFilePath, l)
	guest.BillingService = billing
	guest.IamService = iam
	guest.ProjectID = keyFile["project_id"].(string)
	guest.SvcEmail = keyFile["client_email"].(string)
	close(c)
}

// FIXME: How do we use Application Default Credentials (ADC) for both the
// Host and the Guest? Maybe: Authenticate as Host and Guest using gcloud CLI
// programmatically. One channel will trigger `gcloud auth login`, user will login
// via browser to the Host and channel will retrieve ADC. Next, another channel
// will do the same, but the user will login to the Guest account.
// Doing this might require that the filesystem is able to update
// %APPDATA%/gcloud/application_default_credentials.json (Windows) or
// $HOME/.config/gcloud/application_default_credentials.json (Other OS') after
// the program has read the JSON for the Host.
func authenticateDefault(config *Config, l *logrus.Logger) error {
	// NOTE: package cloudbilling and package iam both export CloudPlatformScope
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		l.Error(err)
		return err
	}

	l.Debug("Instantiating billing service client...")
	billing, err := cloudbilling.New(client)
	if err != nil {
		l.Error(err)
		return err
	}

	l.Debug("Instantiating IAM service client...")
	iam, err := iam.New(client)
	if err != nil {
		l.Error(err)
		return err
	}

	_ = billing
	_ = iam
	return nil
}

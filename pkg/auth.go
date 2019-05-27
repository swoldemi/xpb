package xpb

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	iam "google.golang.org/api/iam/v1"
)

var scopes = []string{cloudbilling.CloudPlatformScope, iam.CloudPlatformScope}

// FIXME: How do we use Application Default Credentials (ADC) for both the
// Host and the Guest.
// Maybe: Authenticate as Host using gcloud CLI, get ADC, wait on a channel
// for the user to, using the gcloud CLI, authenticate their Guest account and
// retrieve ADC again? This assumes the filesystem is able to update
// %APPDATA%/gcloud/application_default_credentials.json (Windows) or
// $HOME/.config/gcloud/application_default_credentials.json (Other OS') after
// the program has read the JSON for the Host.
func authenticateDefault(config *Config, l *logrus.Entry) error {
	// NOTE: package cloudbilling and package iam both export CloudPlatformScope
	client, err := google.DefaultClient(context.Background(), scopes...)
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

// AuthenticateHost returns an instance of xpb.Host
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateHost(config *Config, l *logrus.Entry, host *Host, c chan struct{}) {
	l.Debug("Attempting to authenticate host using provided service account...")
	ctx := context.Background()
	data, err := ioutil.ReadFile(config.HostKeyFilePath)
	if err != nil {
		l.Fatal(err)
		close(c)
	}
	creds, err := google.CredentialsFromJSON(ctx, data, scopes...)
	if err != nil {
		l.Fatal(err)
		close(c)

	}
	client := oauth2.NewClient(ctx, creds.TokenSource)
	l.Info("Succcessfully created OAuth2 client for host account")

	if config.AddressedProjectID != creds.ProjectID {
		l.Warn("Warning: Configured project ID does not match service account's project ID")
	}

	l.Debug("Instantiating billing service client for host account...")
	billing, err := cloudbilling.New(client)
	if err != nil {
		l.Fatal(err)
		close(c)
	}

	l.Debug("Instantiating IAM service client for host account...")
	iam, err := iam.New(client)
	if err != nil {
		l.Fatal(err)
		close(c)
	}

	host.BillingService = billing
	host.IamService = iam
	host.ProjectID = config.AddressedProjectID
	close(c)
}

// AuthenticateGuest returns an instance of xpb.Guest
// with services authenticated and instantiated using
// the provided service account key file.
func AuthenticateGuest(config *Config, l *logrus.Entry, host *Guest, c chan struct{}) {
	l.Debug("Attempting to authenticate guest using provided service account...")
	ctx := context.Background()
	data, err := ioutil.ReadFile(config.GuestKeyFilePath)
	if err != nil {
		l.Fatal(err)
		close(c)
	}
	creds, err := google.CredentialsFromJSON(ctx, data, iam.CloudPlatformScope)
	l.Info(string(creds.JSON))
	if err != nil {
		l.Fatal(err)
		close(c)

	}
	client := oauth2.NewClient(ctx, creds.TokenSource)
	l.Info("Succcessfully created OAuth2 client for guest account")

	l.Debug("Instantiating billing service client for guest account...")
	billing, err := cloudbilling.New(client)
	if err != nil {
		l.Fatal(err)
		close(c)
	}
	l.Debug("Instantiating IAM service client for guest account...")
	iam, err := iam.New(client)
	if err != nil {
		l.Fatal(err)
		close(c)
	}

	host.BillingService = billing
	host.IamService = iam
	close(c)
}

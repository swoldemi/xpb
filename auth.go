package xpb

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	cloudbilling "google.golang.org/api/cloudbilling/v1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func readKeyFile(filepath string) map[string]interface{} {
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
func AuthenticateHost(config *Config, l *logrus.Entry, host *Host, c chan struct{}) {
	l.Debug("Attempting to authenticate host using provided service account...")
	ctx := context.Background()

	opts := []option.ClientOption{option.WithCredentialsFile(config.HostKeyFilePath)}
	l.Debug("Instantiating Cloud Billing service client for host account...")
	billingsvc, err := cloudbilling.NewService(ctx, append(opts, option.WithScopes("https://www.googleapis.com/auth/cloud-billing"))...)
	Fataler(err)

	l.Debug("Instantiating IAM service client for host account...")
	iamsvc, err := iam.NewService(ctx, opts...)
	Fataler(err)

	l.Debug("Instantiating Cloud Resource Manager service client for host account...")
	resourcemgrsvc, err := cloudresourcemanager.NewService(ctx, append(opts, option.WithScopes(cloudresourcemanager.CloudPlatformScope))...)
	Fataler(err)

	l.Debug("Storing guest's service account information...")
	keyFile := readKeyFile(config.HostKeyFilePath)
	host.BillingService = billingsvc
	host.IamService = iamsvc
	host.ResourceMgrService = resourcemgrsvc
	host.ProjectID = keyFile["project_id"].(string)
	host.SvcEmail = keyFile["client_email"].(string)
	close(c)
}

// AuthenticateGuest fulfills xpb.Guest
// with services authenticated and instantiated using default credentials
func AuthenticateGuest(l *logrus.Entry, guest *Guest, c chan struct{}) {
	l.Debug("Attempting to authenticate guest using using ADC...")
	ctx := context.Background()

	l.Debug("Instantiating Cloud Billing service client for guest account...")
	billingsvc, err := cloudbilling.NewService(ctx)
	Fataler(err)

	l.Debug("Instantiating IAM service client for guest account...")
	iamsvc, err := iam.NewService(ctx)
	Fataler(err)

	guest.BillingService = billingsvc
	guest.IamService = iamsvc
	close(c)
}

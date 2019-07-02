package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	// ErrNoMessages is returned when the guest has no emails in their inbox
	ErrNoMessages = errors.New("email: guest has no emails in their inbox")

	// ErrInviteNotFound is returned when the host's invite is not found in the guest's inbox
	ErrInviteNotFound = errors.New("email: unable to find host's invite in guest's inbox")

	// ErrCredentialsNotFound is returned when the provided credentials file is invalid
	ErrCredentialsNotFound = errors.New("email: unable to read credentials file")

	// ErrCantParseCredentials is returned when the provided credentials file is invalid
	ErrCantParseCredentials = errors.New("email: unable to parse credentials file")

	// ErrFoundInvite is returned to halt page iteration when the desired email is found
	ErrFoundInvite = errors.New("found invite email")
)

func newGmailService(credentials string) (*gmail.Service, error) {
	b, err := ioutil.ReadFile(credentials)
	if err != nil {
		return nil, ErrCredentialsNotFound
	}

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		return nil, ErrCantParseCredentials
	}
	client := getClient(config)

	opts := []option.ClientOption{
		option.WithHTTPClient(client),
	}
	svc, err := gmail.NewService(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "guest-gmail-token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

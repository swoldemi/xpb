package email

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/swoldemi/xpb"
	gmail "google.golang.org/api/gmail/v1"
)

// GuestInbox encapsulates methods
// for interacting with the guest's inbox and accepting
// the GCP Project invite from the host
type GuestInbox struct {
	GmailService *gmail.Service
	Invite       *gmail.Message
	l            *logrus.Entry
}

// New creates a new GuestInbox instance
func New(config *xpb.Config) (*GuestInbox, error) {
	l := xpb.NewLogger("Email")

	l.Debug("Creating Gmail service using provided credentials file...")
	gmailService, err := newGmailService(config.GuestGmailCredsPath)
	if err != nil {
		return nil, err
	}

	l.Debug("Successfully created Gmail service")
	return &GuestInbox{
		GmailService: gmailService,
		l:            l,
	}, nil
}

// FindInvite does a linear search to retrieve messages
// from the guest's inbox and stores the invite message
func (g *GuestInbox) FindInvite() error {
	g.l.Info("Searching for invite in guest's Gmail inbox...")
	inviteSubject := "Join my project on Google Developers Console"

	svc := gmail.NewUsersMessagesService(g.GmailService)
	request := svc.List("me").Q("label:inbox")
	err := request.Pages(context.Background(), func(l *gmail.ListMessagesResponse) error {
		if l.Messages == nil {
			return ErrNoMessages
		}

		for _, message := range l.Messages {
			m, err := svc.Get("me", message.Id).Do()
			if err != nil {
				return err
			}

			if m.Payload == nil {
				continue
			}

			// Determine that we have found the email by
			// matching against an expected Subject line
			for _, h := range m.Payload.Headers {
				if h.Name == "Subject" {
					if strings.Contains(h.Value, inviteSubject) {
						g.l.Infof("Found invite email with ID: %v!", message.Id)
						// Because supplying format= to the api
						// significally slows down message exchanges,
						// only provide it for the message we care about
						invite, err := svc.Get("me", message.Id).Format("RAW").Do()
						if err != nil {
							return err
						}

						g.Invite = invite
						return ErrFoundInvite // Return an error to halt interation
					}
				}
			}
		}
		return nil
	})

	if err != nil && err != ErrFoundInvite {
		return err
	}

	if g.Invite == nil {
		return ErrInviteNotFound
	}

	g.l.Info("Search completed successfully.")
	return nil
}

// ExtractInvite finds the invite url in the message
// found by FindInvite and accepts the invite
func (g *GuestInbox) ExtractInvite() error {
	g.l.Infof("%+v", g.Invite.Payload.Body)
	return nil
}

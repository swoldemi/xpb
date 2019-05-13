package email

// Package email encapsulates functions necessary
// for retrieving the IAM invite sent from `guest`
// to `host`. Accounts are assumed to be Gmail bound.

type HostAccount struct {
	Address  string
	Password string
}

func (h *HostAccount) AcceptInvite() error {
	return nil
}

package util

import (
	"errors"
)

var (
	// ErrGuestNoBilling - the guest account List call returns empty BillingAccounts
	ErrGuestNoBilling = errors.New("xpb: guest has no valid billing accounts")

	// ErrAccountsIdentical - the host account's project is already associated with the guest account's billing account
	ErrAccountsIdentical = errors.New("xpb: host account is already set to guest's billing account")

	// ErrInvalidReply is returned when an expected type assertion fails on an ExectueScript call.
	ErrInvalidReply = errors.New("browser: ExecuteScript returned unexpected reply")
)

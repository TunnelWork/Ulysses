package server

type ServerConfigurables map[string]string
type AccountConfigurables map[string]string

// It is up to module designer to parse/utilize the AccountUsage.
type AccountUsage interface {
	ForClient() (usage string)
	ForAdminUsage() (usage string)
}

// It is up to module designer to parse/utilize the Credential.
type Credential interface {
	ForClient() (credential string)
	ForAdmin() (credential string)
}

type Server interface {
	//////// Mandatory Functions: will be called by Ulysses Core

	// AddAccount() takes in:
	// - (a ptr to) a ServerConfigurables specifically designed for the target server. (e.g., Database credentials, IP addresses)
	// - (a ptr to) a AccountConfigurables specifically designed for the account to be created. (e.g., Password, Service Port, Quota)
	// And returns an integer & an error
	// - if err == nil: int as the Account ID, aka Service ID/Product ID, for the newly created account. It's caller's responsibility to store the accid and (possibly) associate it with a user.
	// - otherwise, discard the int and check the err
	AddAccount(sconf *ServerConfigurables, aconf *AccountConfigurables) (accid int, err error)

	// UpdateAccount() takes in:
	// - an int as the Account ID specifying the exact account needs to be updated
	// - (a ptr to) a ServerConfigurables specifically designed for the target server. (e.g., Database credentials, IP addresses)
	// - (a ptr to) a AccountConfigurables specifically designed for the account to be created. (e.g., Password, Service Port, Quota)
	// And returns an error, if nil, success. Otherwise, check the error.
	UpdateAccount(accid int, sconf *ServerConfigurables, aconf *AccountConfigurables) (err error)

	// DeleteAccount() takes in:
	// - an int as the Account ID specifying the exact account needs to be deleted
	// And returns an error, if nil, success. Otherwise, check the error.
	DeleteAccount(accid int) (err error)

	//////// Optional Functions: may be called by Ulysses Extensions

	// GetCredentials() fetch the Credential in JSON string format for a specific Account specified by accid.
	GetCredentials(accid int) (credential []Credential)

	// GetUsage() fetch the history usage of a service
	GetUsage(accid int) (usage AccountUsage)
}

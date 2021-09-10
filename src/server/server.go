package server

// It is up to module designer to parse/utilize the AccountUsage.
type AccountUsage interface {
	ForClient() (usage string)
	ForAdmin() (usage string)
}

// It is up to module designer to parse/utilize the Credential.
type Credential interface {
	ForClient() (credential string)
	ForAdmin() (credential string)
}

type Server interface {
	//////// Mandatory Functions: will be called by Ulysses Core

	// UpdateServer() takes in:
	// - a ServerConfigurables specifically designed for the target server. (e.g., Database credentials, IP addresses)
	UpdateServer(sconf ServerConfigurables) (err error)

	// AddAccount() takes in:
	// - a slice of AccountConfigurables specifically designed for each accounts to be created. (e.g., Password, Service Port, Quota)
	// And returns a slice of integer & an error
	// - if err == nil: accID includes the Account ID for each account added, aka Service ID/Product ID, for the newly created account.
	// It's caller's responsibility to store the accid and (possibly) associate it with a user.
	// - otherwise, accID includes the Account ID for each account added BEFORE the err occurs. (No more operation after err)
	AddAccount(aconf []AccountConfigurables) (accID []int, err error)

	// UpdateAccount() takes in:
	// - a slice of int as the Account ID specifying each account needs to be updated
	// - a slice of AccountConfigurables specifically designed for the account to be created. (e.g., Password, Service Port, Quota)
	// And returns a slice of integer & an error
	// - if err == nil: accID includes the Account ID for each account updated, aka Service ID/Product ID, for the updated account.
	// - otherwise, accID includes the Account ID for each account updated BEFORE the err occurs. (No more operation after err)
	UpdateAccount(accID []int, aconf []AccountConfigurables) (successAccID []int, err error)

	// DeleteAccount() takes in:
	// - a slice of int as the Account ID specifying each account needs to be deleted
	// And returns a slice of integer & an error
	// - if err == nil: accID includes the Account ID for each account deleted, aka Service ID/Product ID, for the deleted account.
	// - otherwise, accID includes the Account ID for each account deleted BEFORE the err occurs. (No more operation after err)
	DeleteAccount(accID []int) (successAccID []int, err error)

	//////// Optional Functions: may be called by Ulysses 3rd-party Extensions

	// GetCredentials() fetch Credentials in JSON string format for each Account specified by accID.
	GetCredentials(accID []int) (credentials []Credential, err error)

	// GetUsage() fetch the history usages of each service specified by accID
	GetUsage(accID []int) (usages []AccountUsage, err error)
}

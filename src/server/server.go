package server

// Server interface-compatible structs should be copyable.
// Recommended design:
// - Pointer to struct
// - Member pointers in struct
type Server interface {
	//// Server Settings: Alter the local information needed to interact with a backend server.

	// UpdateServer() should update the internal variable of a Server to reflect the new state.
	UpdateServer(sconf Configurables) (err error)

	//// Account Operations: Connects to the backend server to perform operations for user accounts.

	// AddAccount() utilizes internal server configuration and aconf pased in to create a series of accounts with
	// same sconf and variable aconf in order.
	// This function returns immediately upon an error has occured. The returned accID should contain IDs for
	// all successfully created accounts.
	AddAccount(aconf []Configurables) (accID []int, err error)

	// UpdateAccount() utilizes internal server configuration and aconf pased in to update a series of accounts
	// specified by accID.
	// This function returns immediately upon an error has occured. The returned successAccID should contain
	// IDs for all successfully updated accounts.
	UpdateAccount(accID []int, aconf []Configurables) (successAccID []int, err error)

	// DeleteAccount() utilizes internal server configuration to delete a series of accounts specified by accID.
	// This function returns immediately upon an error has occured. The returned successAccID should contain
	// IDs for all successfully deleted accounts.
	DeleteAccount(accID []int) (successAccID []int, err error)

	//// User Interface Helpers: Acquire needed informations from Server for the User Interface or Admin Panel.

	// GetCredentials() fetch Credentials in JSON string format for each Account specified by accID.
	GetCredentials(accID []int) ([]Credential, error)

	// GetUsage() fetch the history usages of each service specified by accID
	GetUsage(accID []int) ([]AccountUsage, error)
}

package server

// It is up to module designer to parse/utilize the Credential.
type Credential interface {
	ForClient() (credential string)
	ForAdmin() (credential string)
}

package server

// It is up to module designer to parse/utilize the AccountUsage.
type AccountUsage interface {
	ForClient() (usage string)
	ForAdmin() (usage string)
}

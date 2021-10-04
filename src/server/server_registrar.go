package server

var (
	regManagers ServerRegistrarMap = ServerRegistrarMap{}
)

// ServerRegistrar interface-compatible structs should be copyable.
// Recommended design:
// - Pointer to struct
// - Member pointers in struct
type ServerRegistrar interface {
	// NewServer returns a Server interface with internal state set to reflect sconf.
	NewServer(sconf Configurables) (Server, error)
}

type ServerRegistrarMap map[string]ServerRegistrar

// AddServerRegistrar adds a registrar to the global ServerRegistrarMap
func AddServerRegistrar(serverTypeName string, serverReg ServerRegistrar) {
	if regManagers == nil {
		regManagers = ServerRegistrarMap{}
	}
	regManagers[serverTypeName] = serverReg
}

// NewServerByType returns a Server interface specified by serverType according to the ServerRegistrarMap
// the internal state of the returned Server interface should reflect sconf.
func NewServerByType(serverType string, sconf Configurables) (Server, error) {
	var newServer Server
	var err error

	if sr, ok := regManagers[serverType]; ok {
		newServer, err = sr.NewServer(sconf)
	} else {
		err = ErrServerUnknown
	}

	return newServer, err
}

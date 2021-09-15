package server

var (
	regManagers ServerRegistrarMap
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
	regManagers[serverTypeName] = serverReg
}

// NewServerByType returns a Server interface specified by serverType according to the ServerRegistrarMap
// the internal state of the returned Server interface should reflect sconf.
func (srm ServerRegistrarMap) NewServerByType(serverType string, sconf Configurables) (Server, error) {
	return srm[serverType].NewServer(sconf)
}

package server

var (
	regManagers ServerRegistrarMap
)

// ServerRegistrar interface-compatible structs should be copyable.
// Recommended design:
// - Pointer to struct
// - Member pointers in struct
type ServerRegistrar interface {
	NewServer(sconf ServerConfigurables) (Server, error)
}

type ServerRegistrarMap map[string]ServerRegistrar

func AddServerRegistrar(serverTypeName string, serverReg ServerRegistrar) {
	regManagers[serverTypeName] = serverReg
}

func (srm ServerRegistrarMap) NewServerByType(serverType string, sconf ServerConfigurables) (Server, error) {
	return srm[serverType].NewServer(sconf)
}

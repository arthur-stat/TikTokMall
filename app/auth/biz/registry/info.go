package registry

import "net"

// Info contains the information for service registration
type Info struct {
	ServiceName string
	Addr        net.Addr
	Weight      int
	Tags        []string
}

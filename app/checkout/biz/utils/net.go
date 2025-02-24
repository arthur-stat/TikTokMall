package utils

import "net"

// NewNetAddr creates a new net.Addr
func NewNetAddr(network, addr string) net.Addr {
	return &netAddr{network: network, addr: addr}
}

type netAddr struct {
	network string
	addr    string
}

func (n *netAddr) Network() string { return n.network }
func (n *netAddr) String() string  { return n.addr }

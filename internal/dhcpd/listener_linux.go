//go:build linux
// +build linux

package dhcpd

import (
	"net"

	"github.com/krolaw/dhcp4/conn"
)

// udpListener returns a PacketConn bound to the specified interface.
func udpListener(iface, addr string) (net.PacketConn, error) {
	return conn.NewUDP4BoundListener(iface, addr)
}

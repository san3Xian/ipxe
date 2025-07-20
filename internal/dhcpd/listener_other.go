//go:build !linux
// +build !linux

package dhcpd

import (
	"net"

	"github.com/krolaw/dhcp4/conn"
)

// udpListener creates a PacketConn filtered for the specified interface.
func udpListener(iface, addr string) (net.PacketConn, error) {
	return conn.NewUDP4FilterListener(iface, addr)
}

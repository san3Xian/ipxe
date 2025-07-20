//go:build !linux
// +build !linux

package dhcpd

import (
	"net"
	"time"

	"github.com/krolaw/dhcp4/conn"
)

// serveIfPacket wraps conn.serveIfConn to satisfy net.PacketConn.
type serveIfPacket struct {
	*conn.serveIfConn
	laddr net.Addr
}

func (s *serveIfPacket) LocalAddr() net.Addr                { return s.laddr }
func (s *serveIfPacket) SetDeadline(t time.Time) error      { return s.conn.SetDeadline(t) }
func (s *serveIfPacket) SetReadDeadline(t time.Time) error  { return s.conn.SetReadDeadline(t) }
func (s *serveIfPacket) SetWriteDeadline(t time.Time) error { return s.conn.SetWriteDeadline(t) }

// udpListener creates a PacketConn filtered for the specified interface.
func udpListener(iface, addr string) (net.PacketConn, error) {
	c, err := conn.NewUDP4FilterListener(iface, addr)
	if err != nil {
		return nil, err
	}
	laddr, _ := net.ResolveUDPAddr("udp4", addr)
	return &serveIfPacket{serveIfConn: c, laddr: laddr}, nil
}

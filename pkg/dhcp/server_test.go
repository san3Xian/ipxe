package dhcp

import (
	dhcp4 "github.com/krolaw/dhcp4"
	"net"
	"testing"
	"time"
)

func TestServeDHCP(t *testing.T) {
	s := NewServer(net.IPv4(192, 168, 1, 100), net.IPv4(192, 168, 1, 200), net.IPv4(192, 168, 1, 1), time.Hour)
	req := dhcp4.RequestPacket(dhcp4.Discover, net.HardwareAddr{1, 2, 3, 4, 5, 6}, nil, []byte{1, 2, 3, 4}, false, nil)
	resp := s.ServeDHCP(req, dhcp4.Discover, nil)
	if resp == nil {
		t.Fatal("no response")
	}
}

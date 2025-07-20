package dhcp

import (
	"net"
	"time"

	dhcp4 "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"
)

// Lease holds info about a DHCP lease

type Lease struct {
	IP        net.IP           `json:"ip"`
	MAC       net.HardwareAddr `json:"mac"`
	ExpiresAt time.Time        `json:"expires_at"`
}

// DHCPServer implements dhcp4.Handler and stores leases in memory

func NewServer(rangeStart, rangeEnd, router net.IP, leaseTime time.Duration) *DHCPServer {
	return &DHCPServer{start: rangeStart, end: rangeEnd, router: router, leaseTime: leaseTime, leases: make(map[int]Lease)}
}

type DHCPServer struct {
	start, end net.IP
	router     net.IP
	leaseTime  time.Duration
	leases     map[int]Lease
}

func (d *DHCPServer) ServeDHCP(p dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) dhcp4.Packet {
	mac := p.CHAddr()
	switch msgType {
	case dhcp4.Discover:
		log.Infof("DHCP discover from %s", mac)
		ip := d.start.To4()
		return dhcp4.ReplyPacket(p, dhcp4.Offer, d.router, ip, d.leaseTime,
			d.options())
	case dhcp4.Request:
		log.Infof("DHCP request from %s", mac)
		ip := d.start.To4()
		d.leases[int(ip[3])] = Lease{IP: ip, MAC: mac, ExpiresAt: time.Now().Add(d.leaseTime)}
		return dhcp4.ReplyPacket(p, dhcp4.ACK, d.router, ip, d.leaseTime, d.options())
	default:
		return nil
	}
}

func (d *DHCPServer) options() []dhcp4.Option {
	return []dhcp4.Option{
		{Code: dhcp4.OptionSubnetMask, Value: []byte{255, 255, 255, 0}},
		{Code: dhcp4.OptionRouter, Value: []byte(d.router)},
	}
}

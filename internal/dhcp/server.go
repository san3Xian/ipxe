package dhcp

import (
	"bytes"
	"net"
	"time"

	dhcp4 "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"

	"github.com/example/dpxe/internal/state"
)

// Server implements dhcp4.Handler interface

type Server struct {
	iface      string
	startIP    net.IP
	endIP      net.IP
	router     net.IP
	dns        net.IP
	leaseTime  time.Duration
	subnetMask net.IP
	leases     *state.State
	pxeURL     string
}

func New(iface string, startIP, endIP, router, dns, subnetMask net.IP, leaseTime time.Duration, pxeURL string, st *state.State) *Server {
	return &Server{iface: iface, startIP: startIP, endIP: endIP, router: router, dns: dns, leaseTime: leaseTime, subnetMask: subnetMask, leases: st, pxeURL: pxeURL}
}

func (s *Server) Serve() error {
	log.Infof("DHCP server listening on :67")
	return dhcp4.ListenAndServe(s)
}

func (s *Server) ServeDHCP(p dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) dhcp4.Packet {
	mac := p.CHAddr().String()
	switch msgType {
	case dhcp4.Discover:
		log.Infof("DHCP discover from %s", mac)
		ip := s.nextIP(mac)
		if ip == nil {
			return nil
		}
		opts := []dhcp4.Option{
			{Code: dhcp4.OptionSubnetMask, Value: []byte(s.subnetMask.To4())},
			{Code: dhcp4.OptionRouter, Value: []byte(s.router.To4())},
			{Code: dhcp4.OptionDomainNameServer, Value: []byte(s.dns.To4())},
			{Code: dhcp4.OptionBootFileName, Value: []byte(s.pxeURL)},
			{Code: dhcp4.OptionIPAddressLeaseTime, Value: dhcp4.OptionsLeaseTime(s.leaseTime)},
		}
		return dhcp4.ReplyPacket(p, dhcp4.Offer, s.router, ip, s.leaseTime, opts)
	case dhcp4.Request:
		log.Infof("DHCP request from %s", mac)
		reqIP := net.IP(options[dhcp4.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = p.CIAddr()
		}
		if reqIP == nil {
			return nil
		}
		if !s.ipInRange(reqIP) {
			return dhcp4.ReplyPacket(p, dhcp4.NAK, s.router, nil, 0, nil)
		}
		expiry := time.Now().Add(s.leaseTime)
		s.leases.AddLease(p.CHAddr(), reqIP, expiry)
		opts := []dhcp4.Option{
			{Code: dhcp4.OptionSubnetMask, Value: []byte(s.subnetMask.To4())},
			{Code: dhcp4.OptionRouter, Value: []byte(s.router.To4())},
			{Code: dhcp4.OptionDomainNameServer, Value: []byte(s.dns.To4())},
			{Code: dhcp4.OptionBootFileName, Value: []byte(s.pxeURL)},
			{Code: dhcp4.OptionIPAddressLeaseTime, Value: dhcp4.OptionsLeaseTime(s.leaseTime)},
		}
		return dhcp4.ReplyPacket(p, dhcp4.ACK, s.router, reqIP, s.leaseTime, opts)
	}
	return nil
}

func (s *Server) ipInRange(ip net.IP) bool {
	return bytesCompare(ip, s.startIP) >= 0 && bytesCompare(ip, s.endIP) <= 0
}

func bytesCompare(a, b net.IP) int {
	return bytes.Compare(a.To4(), b.To4())
}

func (s *Server) nextIP(mac string) net.IP {
	// if already leased
	s.leases.Mu.Lock()
	for _, lease := range s.leases.Leases {
		if lease.MAC.String() == mac {
			s.leases.Mu.Unlock()
			return lease.IP
		}
	}
	for ip := s.startIP.To4(); bytesCompare(ip, s.endIP) <= 0; ip[3]++ {
		used := false
		for _, lease := range s.leases.Leases {
			if lease.IP.Equal(ip) {
				used = true
				break
			}
		}
		if !used {
			s.leases.Mu.Unlock()
			return net.IP{ip[0], ip[1], ip[2], ip[3]}
		}
	}
	s.leases.Mu.Unlock()
	return nil
}

package dhcpd

import (
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"
)

// Lease represents a leased IP address and its expiry time.
type Lease struct {
	IP     net.IP
	Expiry time.Time
}

// DHCPServer provides a minimal DHCP server implementation.
type DHCPServer struct {
	iface         string
	ipRangeStart  net.IP
	ipRangeEnd    net.IP
	router        net.IP
	dns           net.IP
	leaseDuration time.Duration
	leases        map[string]Lease
}

// NewDHCPServer creates a DHCPServer using the provided network settings.
func NewDHCPServer(iface string, start, end, router, dns net.IP, leaseDur time.Duration) *DHCPServer {
	return &DHCPServer{
		iface:         iface,
		ipRangeStart:  start,
		ipRangeEnd:    end,
		router:        router,
		dns:           dns,
		leaseDuration: leaseDur,
		leases:        make(map[string]Lease),
	}
}

// availableIP returns the first free IP within the configured range.
func (d *DHCPServer) availableIP() net.IP {
	for ip := d.ipRangeStart.To4(); !ip.Equal(d.ipRangeEnd); incIP(ip) {
		used := false
		for _, l := range d.leases {
			if l.IP.Equal(ip) && time.Now().Before(l.Expiry) {
				used = true
				break
			}
		}
		if !used {
			return dhcp.IPAdd(ip, 0)
		}
	}
	return nil
}

// incIP increments an IPv4 address in place.
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

// ServeDHCP implements dhcp.Handler.
func (d *DHCPServer) ServeDHCP(pkt dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	mac := pkt.CHAddr().String()
	switch msgType {
	case dhcp.Discover:
		log.Infof("DHCP DISCOVER from %s", mac)
		ip := d.availableIP()
		if ip == nil {
			log.Warn("No IP available for", mac)
			return nil
		}
		return dhcp.ReplyPacket(pkt, dhcp.Offer, d.router, ip, d.leaseDuration,
			d.options())
	case dhcp.Request:
		log.Infof("DHCP REQUEST from %s", mac)
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(pkt.CIAddr())
		}
		if reqIP == nil || reqIP.Equal(net.IPv4zero) {
			reqIP = d.availableIP()
		}
		d.leases[mac] = Lease{IP: reqIP, Expiry: time.Now().Add(d.leaseDuration)}
		log.Infof("Offering %s to %s", reqIP, mac)
		return dhcp.ReplyPacket(pkt, dhcp.ACK, d.router, reqIP, d.leaseDuration,
			d.options())
	case dhcp.Release, dhcp.Decline:
		log.Infof("DHCP %s from %s", msgType.String(), mac)
		delete(d.leases, mac)
	}
	return nil
}

// options returns static DHCP options used for replies.
func (d *DHCPServer) options() []dhcp.Option {
	return []dhcp.Option{
		{Code: dhcp.OptionSubnetMask, Value: net.IP(net.CIDRMask(24, 32))},
		{Code: dhcp.OptionRouter, Value: []byte(d.router)},
		{Code: dhcp.OptionDomainNameServer, Value: []byte(d.dns)},
	}
}

// Serve starts listening and processing DHCP packets.
func (d *DHCPServer) Serve() error {
	var (
		c   net.PacketConn
		err error
	)
	if d.iface != "" {
		c, err = udpListener(d.iface, ":67")
	} else {
		c, err = net.ListenPacket("udp4", ":67")
	}
	if err != nil {
		return err
	}
	log.Infof("DHCP server listening on %s", c.LocalAddr())
	return dhcp.Serve(c, d)
}

// Leases returns a copy of currently active leases.
func (d *DHCPServer) Leases() map[string]Lease {
	out := make(map[string]Lease)
	for k, v := range d.leases {
		if time.Now().Before(v.Expiry) {
			out[k] = v
		}
	}
	return out
}

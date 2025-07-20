package dhcp

import (
	"net"
	"time"

	dhcp4 "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"
)

// Lease represents a DHCP lease
type Lease struct {
	MAC net.HardwareAddr
	IP  net.IP
	Exp time.Time
}

// Server represents a simple DHCP server
type Server struct {
	ipRangeStart net.IP
	ipRangeEnd   net.IP
	leases       map[string]*Lease
	iface        *net.UDPConn
}

// New creates a new DHCP server
func New(start, end net.IP) *Server {
	return &Server{
		ipRangeStart: start,
		ipRangeEnd:   end,
		leases:       make(map[string]*Lease),
	}
}

// Serve starts serving DHCP requests on the specified addr
func (s *Server) Serve(addr string) error {
	conn, err := net.ListenPacket("udp4", addr)
	if err != nil {
		return err
	}
	log.Infof("DHCP server listening on %s", addr)
	s.iface = conn.(*net.UDPConn)
	dhcpServer := &dhcpHandler{s}
	return dhcp4.Serve(conn, dhcpServer)
}

type dhcpHandler struct{ s *Server }

func (h *dhcpHandler) ServeDHCP(pkt dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) dhcp4.Packet {
	switch msgType {
	case dhcp4.Discover:
		ip := h.s.nextIP(pkt.CHAddr())
		if ip == nil {
			log.Warn("No free IPs")
			return nil
		}
		log.Infof("Offering %s to %s", ip, pkt.CHAddr())
		return dhcp4.ReplyPacket(pkt, dhcp4.Offer, h.s.ipRangeStart, ip, 12*time.Hour, nil)
	case dhcp4.Request:
		ip := h.s.nextIP(pkt.CHAddr())
		if ip == nil {
			return nil
		}
		h.s.leases[pkt.CHAddr().String()] = &Lease{MAC: pkt.CHAddr(), IP: ip, Exp: time.Now().Add(12 * time.Hour)}
		log.Infof("Acknowledging %s to %s", ip, pkt.CHAddr())
		return dhcp4.ReplyPacket(pkt, dhcp4.ACK, h.s.ipRangeStart, ip, 12*time.Hour, nil)
	}
	return nil
}

func (s *Server) nextIP(mac net.HardwareAddr) net.IP {
	if l, ok := s.leases[mac.String()]; ok && time.Now().Before(l.Exp) {
		return l.IP
	}
	for ip := s.ipRangeStart.To4(); !ip.Equal(s.ipRangeEnd); inc(ip) {
		taken := false
		for _, l := range s.leases {
			if l.IP.Equal(ip) && time.Now().Before(l.Exp) {
				taken = true
				break
			}
		}
		if !taken {
			return ip
		}
	}
	return nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (s *Server) Leases() []*Lease {
	res := []*Lease{}
	for _, l := range s.leases {
		res = append(res, l)
	}
	return res
}

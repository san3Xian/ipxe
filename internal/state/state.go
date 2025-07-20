package state

import (
	"net"
	"sync"
	"time"
)

type Lease struct {
	IP     net.IP
	MAC    net.HardwareAddr
	Expiry time.Time
}

// State holds current leases
// It is safe for concurrent use by multiple goroutines

type State struct {
	Mu     sync.Mutex
	Leases map[string]Lease // keyed by MAC address string
}

func New() *State {
	return &State{Leases: make(map[string]Lease)}
}

// AddLease records lease for mac/ip with expiry
func (s *State) AddLease(mac net.HardwareAddr, ip net.IP, expiry time.Time) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Leases[mac.String()] = Lease{IP: ip, MAC: mac, Expiry: expiry}
}

// GetLease returns lease for mac if exists
func (s *State) GetLease(mac net.HardwareAddr) (Lease, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	l, ok := s.Leases[mac.String()]
	return l, ok
}

// AllLeases returns snapshot of leases
func (s *State) AllLeases() []Lease {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	leases := make([]Lease, 0, len(s.Leases))
	for _, l := range s.Leases {
		leases = append(leases, l)
	}
	return leases
}

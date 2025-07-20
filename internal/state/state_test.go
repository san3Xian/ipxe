package state

import (
	"net"
	"testing"
	"time"
)

func TestState(t *testing.T) {
	st := New()
	mac, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	ip := net.ParseIP("192.168.1.100")
	st.AddLease(mac, ip, time.Now().Add(time.Hour))

	if _, ok := st.GetLease(mac); !ok {
		t.Fatal("expected lease")
	}
	leases := st.AllLeases()
	if len(leases) != 1 {
		t.Fatalf("expected 1 lease, got %d", len(leases))
	}
}

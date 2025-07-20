package ipc

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/example/dpxe/internal/dhcpd"
)

type fakeServer struct{ leases map[string]dhcpd.Lease }

func (f *fakeServer) Leases() map[string]dhcpd.Lease { return f.leases }

func TestIPC(t *testing.T) {
	os.Remove("/tmp/dpxe.sock")
	fs := &fakeServer{leases: map[string]dhcpd.Lease{"aa:bb": {}}}
	go Serve(fs)
	time.Sleep(100 * time.Millisecond)
	resp, err := Query("status")
	if err != nil {
		t.Fatal(err)
	}
	var st Status
	if err := json.Unmarshal([]byte(resp), &st); err != nil {
		t.Fatal(err)
	}
	if len(st.Leases) != 1 {
		t.Fatalf("expected 1 lease")
	}
}

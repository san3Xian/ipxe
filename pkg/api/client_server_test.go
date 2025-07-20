package api

import (
	"os"
	"testing"
	"time"

	"github.com/example/dpxe/pkg/dhcp"
)

type dummyDHCP struct{}

func (d *dummyDHCP) Leases() []*dhcp.Lease {
	return []*dhcp.Lease{{}}
}

func TestAPI(t *testing.T) {
	sock := os.TempDir() + "/dpxe_test.sock"
	os.Remove(sock)
	srv := APIServer{SockPath: sock, DHCP: &dhcp.Server{}}
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	// give server time to listen
	time.Sleep(100 * time.Millisecond)

	c := Client{SockPath: sock}
	_, err := c.Query("leases")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
}

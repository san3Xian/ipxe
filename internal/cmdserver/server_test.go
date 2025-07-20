package cmdserver

import (
	"os"
	"testing"
	"time"

	"github.com/example/dpxe/internal/state"
)

func TestServerClient(t *testing.T) {
	path := "/tmp/dpxe_test.sock"
	st := state.New()
	srv := New(path, st)

	done := make(chan struct{})
	go func() {
		srv.Serve()
		close(done)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	resp, err := ClientQuery(path, "leases")
	if err != nil {
		t.Fatal(err)
	}
	if resp == "" {
		t.Fatal("expected response")
	}
	srv.Close()
	os.Remove(path)
	<-done
}

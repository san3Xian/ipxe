package ipc

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"os"

	"github.com/example/dpxe/internal/dhcpd"
	log "github.com/sirupsen/logrus"
)

// socketPath is the UNIX domain socket used for IPC.
const socketPath = "/tmp/dpxe.sock"

// LeaseProvider provides DHCP lease information for IPC queries.
type LeaseProvider interface {
	Leases() map[string]dhcpd.Lease
}

// Status represents the JSON structure returned to clients.
type Status struct {
	Leases map[string]dhcpd.Lease `json:"leases"`
}

// Serve starts listening on the UNIX socket and handles requests.
func Serve(p LeaseProvider) {
	os.Remove(socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Errorf("IPC listen error: %v", err)
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		go handleConn(conn, p)
	}
}

// handleConn processes a single IPC connection.
func handleConn(c net.Conn, p LeaseProvider) {
	defer c.Close()
	r := bufio.NewReader(c)
	line, _ := r.ReadString('\n')
	switch line {
	case "status\n":
		st := Status{Leases: p.Leases()}
		b, _ := json.MarshalIndent(st, "", "  ")
		c.Write(b)
	default:
		c.Write([]byte("unknown command"))
	}
}

// Query connects to the running daemon and executes a command.
func Query(cmd string) (string, error) {
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		return "", err
	}
	defer c.Close()
	_, err = c.Write([]byte(cmd + "\n"))
	if err != nil {
		return "", err
	}
	resp, err := io.ReadAll(c)
	return string(resp), err
}

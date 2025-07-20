package api

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/example/dpxe/pkg/dhcp"
	log "github.com/sirupsen/logrus"
)

// APIServer handles client requests over unix socket
type APIServer struct {
	SockPath string
	DHCP     *dhcp.Server
}

// Start starts the API server
func (s *APIServer) Start() error {
	l, err := net.Listen("unix", s.SockPath)
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go s.handle(conn)
		}
	}()
	log.Infof("API server listening on %s", s.SockPath)
	return nil
}

func (s *APIServer) handle(c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "leases":
			for _, l := range s.DHCP.Leases() {
				fmt.Fprintf(c, "%s %s %s\n", l.MAC, l.IP, l.Exp.Format(time.RFC3339))
			}
		default:
			fmt.Fprintf(c, "unknown command\n")
		}
	}
}

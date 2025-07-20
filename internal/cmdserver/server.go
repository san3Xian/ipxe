package cmdserver

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/example/dpxe/internal/state"
)

type Server struct {
	Path     string
	State    *state.State
	listener net.Listener
}

func New(path string, st *state.State) *Server {
	return &Server{Path: path, State: st}
}

func (s *Server) Serve() error {
	if err := os.RemoveAll(s.Path); err != nil {
		return err
	}
	l, err := net.Listen("unix", s.Path)
	if err != nil {
		return err
	}
	s.listener = l
	log.Infof("command server listening on %s", s.Path)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

// Close stops the server
func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handle(c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case "leases":
			leases := s.State.AllLeases()
			data, _ := json.MarshalIndent(leases, "", "  ")
			c.Write(append(data, '\n'))
		default:
			c.Write([]byte("unknown command\n"))
		}
	}
}

// ClientQuery connects to socket and sends command, printing response
func ClientQuery(path string, cmd string) (string, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if _, err := conn.Write([]byte(cmd + "\n")); err != nil {
		return "", err
	}
	resp, err := bufio.NewReader(conn).ReadBytes('\n')
	if err == nil {
		return string(resp), nil
	}
	// if not newline terminated, read rest
	b, _ := io.ReadAll(conn)
	return string(resp) + string(b), err
}

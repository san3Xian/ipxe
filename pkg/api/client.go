package api

import (
	"bufio"
	"fmt"
	"net"
)

// Client is used to communicate with the running dpxe server
type Client struct {
	SockPath string
}

// Query sends a command and returns the response
func (c *Client) Query(cmd string) (string, error) {
	conn, err := net.Dial("unix", c.SockPath)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	fmt.Fprintln(conn, cmd)
	if uc, ok := conn.(*net.UnixConn); ok {
		uc.CloseWrite()
	}
	scanner := bufio.NewScanner(conn)
	var resp string
	for scanner.Scan() {
		resp += scanner.Text() + "\n"
	}
	return resp, scanner.Err()
}

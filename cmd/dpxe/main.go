package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"dpxe/pkg/config"
	"dpxe/pkg/dhcp"
	"dpxe/pkg/httpd"
	"dpxe/pkg/ipxe"
	dhcp4 "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"
)

var configPath = flag.String("config", "config.yaml", "config file")
var client = flag.Bool("client", false, "query server status")

func main() {
	flag.Parse()
	if *client {
		resp, err := http.Get("http://localhost:8080/status")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		fmt.Println(string(b))
		return
	}

	c, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	script := ipxe.Script(c)
	os.WriteFile("boot.ipxe", []byte(script), 0644)

	go func() {
		s := httpd.New(c.HTTPRoot)
		log.Fatal(s.Start(c.Addr))
	}()

	dhcpServer := dhcp.NewServer(net.IPv4(192, 168, 1, 100), net.IPv4(192, 168, 1, 200), net.IPv4(192, 168, 1, 1), time.Hour)
	log.Infof("DHCP server listening on %s", c.Addr)
	l, err := net.ListenPacket("udp4", ":67")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(dhcp4.Serve(l, dhcpServer))
}

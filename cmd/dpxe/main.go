package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/example/dpxe/pkg/api"
	"github.com/example/dpxe/pkg/config"
	"github.com/example/dpxe/pkg/dhcp"
	"github.com/example/dpxe/pkg/ipxe"
)

func main() {
	clientMode := flag.Bool("client", false, "client mode")
	cfgPath := flag.String("config", "config.yaml", "config file")
	sock := flag.String("sock", filepath.Join(os.TempDir(), "dpxe.sock"), "unix socket path")
	flag.Parse()

	if *clientMode {
		c := api.Client{SockPath: *sock}
		resp, err := c.Query("leases")
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.WriteString(resp)
		return
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dhcpServer := dhcp.New(net.IP{192, 168, 0, 100}, net.IP{192, 168, 0, 200})
	go func() {
		if err := dhcpServer.Serve(":6767"); err != nil {
			log.Fatal(err)
		}
	}()

	ipxeServer := ipxe.New(cfg)
	go func() {
		if err := ipxeServer.Serve(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	apiServer := api.APIServer{SockPath: *sock, DHCP: dhcpServer}
	if err := apiServer.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}

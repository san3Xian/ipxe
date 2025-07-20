package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/example/dpxe/internal/config"
	"github.com/example/dpxe/internal/dhcpd"
	"github.com/example/dpxe/internal/httpd"
	"github.com/example/dpxe/internal/ipc"
	"github.com/example/dpxe/internal/ipxe"
)

func main() {
	cfgPath := flag.String("config", "dpxe.yaml", "config file")
	client := flag.Bool("client", false, "query running daemon")
	flag.Parse()

	if *client {
		resp, err := ipc.Query("status")
		if err != nil {
			fmt.Println("query error:", err)
			os.Exit(1)
		}
		fmt.Println(resp)
		return
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	start := net.ParseIP(cfg.DHCP.StartIP)
	end := net.ParseIP(cfg.DHCP.EndIP)
	router := net.ParseIP(cfg.DHCP.Router)
	dns := net.ParseIP(cfg.DHCP.DNS)
	dhcp := dhcpd.NewDHCPServer("", start, end, router, dns, time.Duration(cfg.DHCP.LeaseDuration)*time.Second)

	go func() {
		if err := dhcp.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	go ipc.Serve(dhcp)

	srv := httpd.New(":8080", cfg.HTTPRoot, ipxe.New(cfg.IPXEScript))
	log.Infof("serving ipxe script %s", cfg.IPXEScript)
	log.Infof("serving http root %s", cfg.HTTPRoot)
	if err := srv.Serve(); err != nil {
		log.Fatal(err)
	}
}

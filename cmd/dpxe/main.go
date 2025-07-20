package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/example/dpxe/internal/cmdserver"
	"github.com/example/dpxe/internal/config"
	"github.com/example/dpxe/internal/dhcp"
	"github.com/example/dpxe/internal/httpserver"
	"github.com/example/dpxe/internal/ipxe"
	"github.com/example/dpxe/internal/state"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "config file")
	client := flag.Bool("client", false, "client mode")
	flag.Parse()

	if *client {
		c, err := config.Load(*cfgPath)
		if err != nil {
			fmt.Println("failed to load config:", err)
			os.Exit(1)
		}
		resp, err := cmdserver.ClientQuery(c.Socket, "leases")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(resp)
		return
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	st := state.New()

	mux := http.NewServeMux()
	mux.Handle("/ipxe", ipxe.Handler(cfg.HTTPRoot, cfg.PXEFile))
	fs := http.FileServer(http.Dir(cfg.HTTPRoot))
	mux.Handle("/", fs)

	httpSrv := httpserver.New(":8080", mux)

	// Build pxe url from http server address
	pxeURL := fmt.Sprintf("http://%s/ipxe", httpSrv.Addr)

	startIP := net.ParseIP(cfg.DHCP.StartIP)
	endIP := net.ParseIP(cfg.DHCP.EndIP)
	router := net.ParseIP(cfg.DHCP.Router)
	dns := net.ParseIP(cfg.DHCP.DNS)
	mask := net.ParseIP(cfg.DHCP.SubnetMask)

	dhcpSrv := dhcp.New(cfg.DHCP.Interface, startIP, endIP, router, dns, mask, cfg.DHCP.LeaseTime, pxeURL, st)

	cmdSrv := cmdserver.New(cfg.Socket, st)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := dhcpSrv.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := cmdSrv.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Info("shutting down")
}

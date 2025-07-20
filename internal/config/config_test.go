package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	data := []byte(`
http_root: /tmp
pxe_file: boot.ipxe
socket: /tmp/dpxe.sock
dhcp:
  interface: eth0
  start_ip: 192.168.1.100
  end_ip: 192.168.1.110
  router: 192.168.1.1
  dns: 8.8.8.8
  lease_time: 1h
  subnet_mask: 255.255.255.0
`)
	tmp, err := os.CreateTemp("", "cfg*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	cfg, err := Load(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DHCP.Interface != "eth0" || cfg.PXEFile != "boot.ipxe" {
		t.Fatal("unexpected values")
	}
}

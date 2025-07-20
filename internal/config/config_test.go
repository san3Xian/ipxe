package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	data := []byte(`
iso_path: /tmp/test.iso
http_root: /tmp/root
ipxe_script: /tmp/script.ipxe
dhcp:
  start_ip: 192.168.0.10
  end_ip: 192.168.0.20
  router: 192.168.0.1
  dns: 8.8.8.8
  lease_duration: 60
`)
	f, err := ioutil.TempFile("", "cfg*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if cfg.ISOPath != "/tmp/test.iso" || cfg.HTTPRoot != "/tmp/root" {
		t.Fatalf("unexpected config %+v", cfg)
	}
	if cfg.DHCP.StartIP != "192.168.0.10" {
		t.Fatalf("dhcp start ip wrong")
	}
}

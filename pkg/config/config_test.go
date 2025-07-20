package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	data := []byte("addr: :8080\nhttp_root: /tmp\niso_path: /tmp/a.iso\nboot: sanboot\ndhcp_range: 192.168.1.10-192.168.1.20\n")
	f, err := ioutil.TempFile("", "cfg*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	c, err := Load(f.Name())
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.Addr != ":8080" || c.Boot != Sanboot {
		t.Fatalf("unexpected config: %+v", c)
	}
}

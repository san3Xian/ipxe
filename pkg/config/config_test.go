package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "cfg*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString("isos:\n- name: test\n  kernel: k\n  initrd: i\n  cmdline: c\n")
	tmp.Close()

	cfg, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(cfg.ISOs) != 1 || cfg.ISOs[0].Name != "test" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

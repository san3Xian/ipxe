package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// BootType defines how to boot the system
// either sanboot or memdisk

type BootType string

const (
	Sanboot BootType = "sanboot"
	Memdisk BootType = "memdisk"
)

// Config holds dpxe configuration loaded from YAML
// ISOPath is path to ISO image
// Kernel and Initrd are used when Boot is memdisk
// Addr is server listen address
// HTTPRoot is directory to serve over HTTP
// DHCPRange is IP range for DHCP leases

// ClientInfo contains information about connected clients for status query

type Config struct {
	Addr      string   `yaml:"addr"`
	HTTPRoot  string   `yaml:"http_root"`
	ISOPath   string   `yaml:"iso_path"`
	Kernel    string   `yaml:"kernel"`
	Initrd    string   `yaml:"initrd"`
	Boot      BootType `yaml:"boot"`
	DHCPRange string   `yaml:"dhcp_range"`
}

// Load parses YAML config from file path
func Load(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

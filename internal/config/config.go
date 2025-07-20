package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

// Config defines application configuration loaded from YAML
// http_root: path to static files for HTTP server
// pxe_file: iPXE script relative to http_root
// socket: unix domain socket path for command server
//
// DHCP config contains IP range and network settings

type Config struct {
	HTTPRoot string     `yaml:"http_root"`
	PXEFile  string     `yaml:"pxe_file"`
	Socket   string     `yaml:"socket"`
	DHCP     DHCPConfig `yaml:"dhcp"`
}

type DHCPConfig struct {
	Interface  string        `yaml:"interface"`
	StartIP    string        `yaml:"start_ip"`
	EndIP      string        `yaml:"end_ip"`
	Router     string        `yaml:"router"`
	DNS        string        `yaml:"dns"`
	LeaseTime  time.Duration `yaml:"lease_time"`
	SubnetMask string        `yaml:"subnet_mask"`
}

// Load reads config from a yaml file
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.Socket == "" {
		c.Socket = "/tmp/dpxe.sock"
	}
	return &c, nil
}

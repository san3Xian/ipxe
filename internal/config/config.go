package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// DHCPConfig holds settings for the DHCP server.
type DHCPConfig struct {
	StartIP       string `yaml:"start_ip"`
	EndIP         string `yaml:"end_ip"`
	Router        string `yaml:"router"`
	DNS           string `yaml:"dns"`
	LeaseDuration int    `yaml:"lease_duration"`
}

// Config is the top-level configuration parsed from YAML.
type Config struct {
	ISOPath    string     `yaml:"iso_path"`
	HTTPRoot   string     `yaml:"http_root"`
	IPXEScript string     `yaml:"ipxe_script"`
	DHCP       DHCPConfig `yaml:"dhcp"`
}

// Load reads the YAML configuration from the provided path.
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

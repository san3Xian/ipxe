package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type DHCPConfig struct {
	StartIP       string `yaml:"start_ip"`
	EndIP         string `yaml:"end_ip"`
	Router        string `yaml:"router"`
	DNS           string `yaml:"dns"`
	LeaseDuration int    `yaml:"lease_duration"`
}

type Config struct {
	ISOPath    string     `yaml:"iso_path"`
	HTTPRoot   string     `yaml:"http_root"`
	IPXEScript string     `yaml:"ipxe_script"`
	DHCP       DHCPConfig `yaml:"dhcp"`
}

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

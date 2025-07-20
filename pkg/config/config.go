package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// ISO represents a bootable ISO definition
type ISO struct {
	Name    string `yaml:"name"`
	Kernel  string `yaml:"kernel"`
	Initrd  string `yaml:"initrd"`
	Cmdline string `yaml:"cmdline"`
}

// Config represents the server configuration
type Config struct {
	ISOs []ISO `yaml:"isos"`
}

// Load reads the configuration from a yaml file
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

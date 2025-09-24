package main

import (
	"github.com/BurntSushi/toml"
)

// Config represents the configuration structure
type Config struct {
	Default DefaultConfig `toml:"default"`
}

// DefaultConfig represents the default section of the configuration
type DefaultConfig struct {
	EndpointURL     string `toml:"endpoint_url"`
	Region          string `toml:"region"`
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	Bucket          string `toml:"bucket"`
	ImgURLPrefix    string `toml:"img_url_prefix"`
	Directory       string `toml:"directory"`
	RenameFile      bool   `toml:"rename_file"`
}

// LoadConfigFromFile loads configuration from a TOML file
func LoadConfigFromFile(filename string) (*Config, error) {
	config := &Config{}
	_, err := toml.DecodeFile(filename, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

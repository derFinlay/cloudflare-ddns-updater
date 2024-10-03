package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ApiKey         string   `yaml:"api_key"`
	DDNSComment    string   `yaml:"ddns_comment"`
	Zones          []string `yaml:"zones"`
	UpdateInterval int      `yaml:"update_interval"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	yamlData, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(yamlData, &config)

	return config, nil
}

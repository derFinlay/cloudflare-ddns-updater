package config

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ApiKey         string   `yaml:"api_key"`
	DDNSComment    string   `yaml:"ddns_comment"`
	Zones          []string `yaml:"zones"`
	UpdateInterval int      `yaml:"update_interval"`
}

func LoadConfig() (*Config, error) {

	apiKey := os.Getenv("api_key")
	ddns_comment := os.Getenv("ddns_comment")
	zones := os.Getenv("zones")
	updateInterval := os.Getenv("update_interval")

	zonesArray := strings.Split(zones, ",")

	x, _ := strconv.Atoi(updateInterval)

	if apiKey != "" && ddns_comment != "" && updateInterval != "" && zones != "" {
		return &Config{
			ApiKey:         apiKey,
			DDNSComment:    ddns_comment,
			Zones:          zonesArray,
			UpdateInterval: x,
		}, nil
	}

	config := &Config{}
	yamlData, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(yamlData, &config)

	return config, nil
}

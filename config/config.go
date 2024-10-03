package config

type Config struct {
	ApiKey string `yaml:"api_key"`
	IgnoreComment string `yaml:"ignore_comment"`
	zones ZoneConfig[] `yaml:"zones"`
}

type ZoneConfig struct {
	Records string[]
}

func loadConfig() (*Config, error) {
	config := &Config{}
	yamlData, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}
 	yaml.Unmarshal(yamlData, &config)

	return config, nil
}
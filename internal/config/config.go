package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen ListenConfig `yaml:"listen"`
}

type ListenConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func DefaultConfig() Config {
	return Config{
		Listen: ListenConfig{
			Address: "127.0.0.1",
			Port:    8081,
		},
	}
}

const DefaultConfigPath = "/etc/system-services-ui/system-services-ui.yaml"

func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

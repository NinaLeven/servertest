package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BasePath string `yaml:"base_path"`
	HttpAddr string `yaml:"http_addr"`
}

func parseConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}

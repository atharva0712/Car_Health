package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Range struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
	Min2 float64 `yaml:"min2,omitempty"`  //for tire pressure warning range
	Max2 float64 `yaml:"max2,omitempty"`  //for tire pressure warning range  
}

type MetricConfig struct {
	Unit string `yaml:"unit"`
	Normal Range `yaml:"normal"`
	Warning Range `yaml:"warning"`
	Critical Range `yaml:"critical"`
}

type Config struct { //represents the main configuration structure for the application, including settings for metric generation, metric definitions, and backend configuration.
	Generation struct {
		IntervalSeconds int `yaml:"interval_seconds"`
		Distribution struct {
			Normal int `yaml:"normal"`
			Warning int `yaml:"warning"`
			Critical int `yaml:"critical"`
		} `yaml:"distribution"`
	} `yaml:"generation"`

	Metrics map[string]MetricConfig `yaml:"metrics"`

	Backend struct {
		URL string `yaml:"url"`
	} `yaml:"backend"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

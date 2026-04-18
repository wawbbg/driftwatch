package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceEntry defines a single service to watch for drift.
type ServiceEntry struct {
	Name       string            `yaml:"name"`
	SourcePath string            `yaml:"source_path"`
	DeployedAt string            `yaml:"deployed_at"`
	Labels     map[string]string `yaml:"labels,omitempty"`
}

// Config holds the top-level driftwatch configuration.
type Config struct {
	ConfigFile string
	Services   []ServiceEntry `yaml:"services"`
	OutputFmt  string         `yaml:"output_format"`
}

// Load reads the config from the default or env-specified path.
func Load() (*Config, error) {
	path := os.Getenv("DRIFTWATCH_CONFIG")
	if path == "" {
		path = "driftwatch.yaml"
	}
	return LoadFromFile(path)
}

// LoadFromFile parses a YAML config file at the given path.
func LoadFromFile(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	cfg.ConfigFile = path

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service[%d]: name is required", i)
		}
		if svc.SourcePath == "" {
			return fmt.Errorf("service %q: source_path is required", svc.Name)
		}
		if svc.DeployedAt == "" {
			return fmt.Errorf("service %q: deployed_at is required", svc.Name)
		}
	}
	return nil
}

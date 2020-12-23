package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is config for the application
type Config struct {
	Server struct {
		Host          string `yaml:"host"`
		Port          string `yaml:"port"`
		Env           string `yaml:"env"`
		JWTKey        string `yaml:"jwtKey"`
		ActivationURL string `yaml:"activationUrl"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Database string `yaml:"database"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

// New loads the config from the config file
func New(cfgFile string) (*Config, error) {
	cfg := &Config{}
	cfg.Server.JWTKey = "secret"
	if len(cfgFile) == 0 {
		return cfg, fmt.Errorf("invalid config file %s", cfgFile)
	}

	extension := filepath.Ext(cfgFile)
	if extension == "" || extension != ".yml" {
		return cfg, fmt.Errorf("invalid file extension for file %s extension %s", cfgFile, extension)
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return cfg, fmt.Errorf("file error: %v", err)
	}

	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return cfg, fmt.Errorf("yaml decoder error : %v", err)
	}

	if cfg.Server.ActivationURL == "" {
		cfg.Server.ActivationURL = fmt.Sprintf("http://%s:%s/api/v1/user/activate", cfg.Server.Host, cfg.Server.Port)
	}

	return cfg, nil
}

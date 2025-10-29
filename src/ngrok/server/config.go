package server

import (
	"os"

	"gopkg.in/yaml.v2"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Driver          string `yaml:"driver"`
	DSN             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

// Config 应用总配置
type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{
		// 默认配置
		Database: DatabaseConfig{
			Enabled:         false,
			Driver:          "mysql",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300,
		},
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

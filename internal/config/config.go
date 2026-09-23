package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	yaml "go.yaml.in/yaml/v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST"`
	Port string `yaml:"port" env:"SERVER_PORT"`
}

type DBConfig struct {
	Host            string        `yaml:"host" env:"DB_HOST"`
	Port            string        `yaml:"port" env:"DB_PORT"`
	User            string        `yaml:"user" env:"DB_USER"`
	Password        string        `yaml:"password" env:"DB_PASSWORD"`
	Name            string        `yaml:"name" env:"DB_NAME"`
	SSLMode         string        `yaml:"sslmode" env:"DB_SSLMODE"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load(".env")
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	applyEnvOverrides(&cfg)

	return &cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	setIfEnv(&cfg.DB.Host, "DB_HOST")
	setIfEnv(&cfg.DB.Port, "DB_PORT")
	setIfEnv(&cfg.DB.User, "DB_USER")
	setIfEnv(&cfg.DB.Name, "DB_NAME")
	setIfEnv(&cfg.DB.Password, "DB_PASSWORD")
	setIfEnv(&cfg.DB.SSLMode, "DB_SSLMODE")
	setIfEnv(&cfg.Server.Host, "SERVER_HOST")
	setIfEnv(&cfg.Server.Port, "SERVER_PORT")
}

func setIfEnv(p *string, key string) {
	if v := os.Getenv(key); v != "" {
		*p = v
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// url: "postgres://%s:%s@%s:%s/%s?sslmode=%s",
// 			c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,

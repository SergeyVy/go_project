package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	Storage    string `yaml:"storage" env:"STORAGE_PATH"`
	HTTPServer `yaml:"http_server"` // ← без имени поля (встроенное)
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD"`
}

// MustLoad: читает YAML, если файл указан и существует, затем накрывает значениями из env.
// Если YAML не найден, стартуем только на env-переменных (подходит для Railway).
func MustLoad() *Config {
	var cfg Config

	// 1) Если задан CONFIG_PATH и файл существует — читаем YAML
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
				log.Fatalf("Error reading config file: %s", err)
			}
		}
	}

	// 2) Поверх подхватываем env-переменные (по тегам в структурах)
	_ = cleanenv.ReadEnv(&cfg)

	// 3) Значения по умолчанию/поддержка Railway
	// DSN: если не пришёл ни из YAML, ни из STORAGE_PATH — берём DATABASE_URL (Railway Postgres)
	if cfg.Storage == "" {
		cfg.Storage = os.Getenv("DATABASE_URL")
	}

	// Порт: Railway передаёт через PORT
	if p := os.Getenv("PORT"); p != "" {
		cfg.Address = ":" + p
	}
	if cfg.Address == "" {
		cfg.Address = ":8080"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 4 * time.Second
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 60 * time.Second
	}

	return &cfg
}

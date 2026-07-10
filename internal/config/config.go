package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB     DBConfig
	Redis  RedisConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type RedisConfig struct {
	Addr string
}

type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := &Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
		Redis: RedisConfig{
			Addr: os.Getenv("REDIS_ADDR"),
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
	}

	if cfg.DB.Host == "" {
		return nil, fmt.Errorf("DB_HOST env var not set")
	}
	if cfg.DB.Port == "" {
		return nil, fmt.Errorf("DB_PORT env var not set")
	}
	if cfg.DB.Username == "" {
		return nil, fmt.Errorf("DB_USER env var not set")
	}
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD env var not set")
	}

	return cfg, nil

}

func (c DBConfig) String() string {
	return fmt.Sprintf("Host: %s, Port: %s, User: %s, Password: %s", c.Host, c.Port, c.Username, "****")
}

package config

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string `json:"mongo_uri"`
	Database       string `json:"database"`
	Collection     string `json:"collection"`
	ServerAddress  string `json:"server_address"`
	NotifySchedule string `json:"notify_schedule"`
}

func Load(path string) (Config, error) {
	_ = godotenv.Load("learn_web_service/.env")

	data, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	pass := os.Getenv("gocluster_pass")
	cfg.MongoURI = strings.ReplaceAll(cfg.MongoURI, "<gocluster_pass>", pass)

	return cfg, nil
}

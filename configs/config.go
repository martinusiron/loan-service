package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB_URL string `yaml:"db_url"`
}

func Load() Config {
	f, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var c Config
	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatal("invalid config:", err)
	}

	return c
}

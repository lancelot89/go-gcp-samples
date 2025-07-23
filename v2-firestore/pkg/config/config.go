package config

import (
	"os"
)

type Config struct {
	ProjectID         string
	FirestoreEmulatorHost string
	Port              string
}

func LoadConfig() *Config {
	return &Config{
		ProjectID:         os.Getenv("PROJECT_ID"),
		FirestoreEmulatorHost: os.Getenv("FIRESTORE_EMULATOR_HOST"),
		Port:              os.Getenv("PORT"), // Cloud Run typically uses PORT env var
	}
}

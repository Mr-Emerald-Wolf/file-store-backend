package config

import (
	"log"
	"os"
)

func CheckEnv() {
	envProps := []string{
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"CLIENT_ORIGIN",
		"PORT",
		"ACCESS_SECRET_KEY",
		"REFRESH_SECRET_KEY",
		"REDIS_HOST",
		"REDIS_PASSWORD",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
	}

	for _, k := range envProps {
		if os.Getenv(k) == "" {
			log.Fatalf("Environment variable %s not defined. Terminating application...", k)
		}
	}
}

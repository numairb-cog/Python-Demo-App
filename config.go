package main

import "os"

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresDB       string
	Port             string
}

func LoadConfig() *Config {
	return &Config{
		PostgresUser:     getEnv("DEMO_PGSQL_USER", "test"),
		PostgresPassword: getEnv("DEMO_PGSQL_PASSWORD", "test"),
		PostgresHost:     getEnv("DEMO_PGSQL_HOST", "127.0.0.1"),
		PostgresDB:       getEnv("DEMO_PGSQL_DB", "test"),
		Port:             getEnv("PORT", "9000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

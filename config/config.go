package config

import (
	"fmt"
	"os"

	"go.uber.org/fx"
)

// Config holds application configuration parameters.
type Config struct {
	GRPCServerAddress string
	GRPCClientTarget  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDSN      string
}

// Module exports the Config provider for FX.
var Module = fx.Options(
	fx.Provide(NewConfig),
)

// NewConfig creates a new configuration object, loading values from environment variables with defaults.
func NewConfig() (*Config, error) {
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "1433")
	dbName := getEnv("DB_NAME", "taskdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	return &Config{
		GRPCServerAddress: ":50051",
		GRPCClientTarget:  "localhost:50051",
		DBHost:            dbHost,
		DBPort:            dbPort,
		DBUser:            dbUser,
		DBPassword:        dbPassword,
		DBName:            dbName,
		DBDSN:             dsn,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

package repositories

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(host, user, pass, dbName, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pass, dbName, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func NewGormDBFromEnv() (*gorm.DB, error) {
	host := getEnvWithDefault("PG_HOST", "127.0.0.1")
	user := getEnvWithDefault("PG_USER", "note-taking")
	pass := getEnvWithDefault("PG_PASSWORD", "note-taking")
	dbName := getEnvWithDefault("PG_DATABASE", "postgres")
	port := getEnvWithDefault("PG_PORT", "5432")
	return NewGormDB(host, user, pass, dbName, port)
}

func getEnvWithDefault(env string, defaultValue string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

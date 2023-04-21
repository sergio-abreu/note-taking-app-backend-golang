package repositories

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure"
)

func NewGormDB(host, user, pass, dbName, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pass, dbName, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func NewGormDBFromEnv() (*gorm.DB, error) {
	host := infrastructure.GetEnvWithDefault("PG_HOST", "127.0.0.1")
	user := infrastructure.GetEnvWithDefault("PG_USER", "note-taking")
	pass := infrastructure.GetEnvWithDefault("PG_PASSWORD", "note-taking")
	dbName := infrastructure.GetEnvWithDefault("PG_DATABASE", "note-taking-app")
	port := infrastructure.GetEnvWithDefault("PG_PORT", "5432")
	return NewGormDB(host, user, pass, dbName, port)
}

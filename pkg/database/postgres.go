package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens a postgres GORM connection.
func Connect(host, port, user, password, dbName, sslMode string, _ bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

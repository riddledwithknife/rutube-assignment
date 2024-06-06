package config

import (
	"fmt"
	"os"
)

func GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
}

func GetJWTSecretKey() string {
	return os.Getenv("JWT_SECRET_KEY")
}

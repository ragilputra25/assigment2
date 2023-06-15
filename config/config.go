package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	Database   Database
	Auth       Auth
	Host       string
	ApiSecret  string
}

type Database struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type Auth struct {
	Username string
	Password string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config Config
	config.ServerPort = os.Getenv("PORT")
	config.Database.DBHost = os.Getenv("DB_HOST")
	config.Database.DBPort = os.Getenv("DB_PORT")
	config.Database.DBUser = os.Getenv("DB_USER")
	config.Database.DBPassword = os.Getenv("DB_PASSWORD")
	config.Database.DBName = os.Getenv("DB_NAME")
	config.Auth.Username = os.Getenv("AUTHUSER")
	config.Auth.Password = os.Getenv("AUTHPASS")
	config.Host = os.Getenv("APIHOST")
	config.ApiSecret = os.Getenv("API_SECRET")
	return config
}

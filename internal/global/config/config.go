package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string

	DBAddress  string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string

	WhitelistUrl string
}

var config Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR Error while read environment variable file")
	}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	config.DBUsername = os.Getenv("DB_USERNAME")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBAddress = os.Getenv("DB_ADDRESS")
	config.DBPort = os.Getenv("DB_PORT")
	config.DBName = os.Getenv("DB_NAME")
	config.WhitelistUrl = os.Getenv("WHITELISTED_URLS")
}

func GetConfig() *Config {
	return &config
}

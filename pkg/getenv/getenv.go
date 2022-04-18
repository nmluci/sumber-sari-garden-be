package getenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvironmentVariable() (map[string]string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR Error while read environment variable file")
	}
	
	// Add the local variable key and value below
	envVariable := make(map[string]string)

	envVariable["SERVER_ADDRESS"] = os.Getenv("SERVER_ADDRESS")
	envVariable["DB_USERNAME"] = os.Getenv("DB_USERNAME")
	envVariable["DB_PASSWORD"] = os.Getenv("DB_PASSWORD")
	envVariable["DB_ADDRESS"] = os.Getenv("DB_ADDRESS")
	envVariable["DB_PORT"] = os.Getenv("DB_PORT")
	envVariable["DB_NAME"] = os.Getenv("DB_NAME")
	envVariable["WHITELISTED_URLS"] = os.Getenv("WHITELISTED_URLS")

	return envVariable
}
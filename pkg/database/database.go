package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabase(dbUsername, dbPassword, dbAddress, dbName string) *sql.DB {
	log.Println("INFO GetDatabase database connection: starting database connection process")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
																dbUsername, dbPassword, dbAddress, dbName)

	var db *sql.DB

	for {
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			log.Printf("ERROR GetDatabase sql open connection fatal error: %v\n", err)
			log.Println("INFO GetDatabase re-attempting to reconnect to databse")
			time.Sleep(1 * time.Second)
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("ERROR GetDatabase sql open connection fatal error: %v\n", err)
			log.Println("INFO GetDatabase re-attempting to reconnect to databse")
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	log.Printf("INFO GetDatabase database connection: established succesfully with %s\n", dataSourceName)
	return db
}

func InitDatabaseFromEnvVariable(envVariable map[string]string) *sql.DB {
	return GetDatabase(
		envVariable["DB_USERNAME"],
		envVariable["DB_PASSWORD"],
		envVariable["DB_ADDRESS"],
		envVariable["DB_NAME"],
	)
}
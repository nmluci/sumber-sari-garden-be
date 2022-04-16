package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabase(dbUsername, dbPassword, dbAddress, dbPort, dbName string) *sql.DB {
	log.Println("INFO GetDatabase database connection: starting database connection process")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
																dbUsername, dbPassword, dbAddress, dbPort, dbName)

	var db *sql.DB

	for {
		db, err := sql.Open("mysql", dataSourceName)
		if pingErr := db.Ping(); err != nil || pingErr != nil {
			if err != nil {
				log.Printf("ERROR GetDatabase sql open connection fatal error: %v\n", err)
			} else if pingErr != nil {
				log.Printf("ERROR GetDatabase sql open connection fatal error: %v\n", pingErr)
			}
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
		envVariable["DB_PORT"],
		envVariable["DB_NAME"],
	)
}
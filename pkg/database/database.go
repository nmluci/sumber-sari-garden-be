package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nmluci/sumber-sari-garden/internal/global/config"
)

type DatabaseClient struct {
	*sql.DB
}

func Init() *DatabaseClient {
	log.Printf("INFO GetDatabase database connection: starting database connection process")

	conf := config.GetConfig()

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		conf.DBUsername, conf.DBPassword, conf.DBAddress, conf.DBName)

	var (
		db  *sql.DB
		err error
	)

	for {
		db, err = sql.Open("mysql", dataSourceName)
		if pingErr := db.Ping(); err != nil || pingErr != nil {
			if err != nil {
				log.Printf("ERROR GetDatabase sql open connection fatal error: %v\n", err)
			} else if pingErr != nil {
				log.Printf("ERROR GetDatabase db ping fatal error: %v\n", pingErr)
			}
			log.Println("INFO GetDatabase re-attempting to reconnect to database...")
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	log.Printf("INFO GetDatabase database connection: established successfully with %s\n", dataSourceName)
	return &DatabaseClient{DB: db}
}

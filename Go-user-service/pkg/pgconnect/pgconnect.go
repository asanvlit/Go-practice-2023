package pgconnect

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type ConnectionConfigData struct {
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Host         string
}

func buildConnectionString(configData ConnectionConfigData) string {
	connectionTemplate := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	return fmt.Sprintf(connectionTemplate, configData.Host, configData.Port, configData.Username, configData.Password,
		configData.DatabaseName)
}

func ConnectDatabase(configData ConnectionConfigData) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", buildConnectionString(configData))

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	for i := 0; i < 5; i++ {
		err := db.Ping()

		if err != nil {
			log.Printf("Unsuccessful ping: " + err.Error())
		} else {
			log.Printf("Successfully connected")
			return db, nil
		}
		time.Sleep(10 * time.Second)

		if i == 4 {
			return nil, err
		}
	}

	return db, nil
}

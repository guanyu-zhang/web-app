package initializers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := "host=0.0.0.0 user=nimble password=nimble dbname=nimble port=5432 sslmode=disable"
	//DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB, err = openWithRetry(dsn, 10, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	//if err != nil {
	//	panic("Failed to connect to DB")
	//}
}

func openWithRetry(dsn string, maxRetries int, backoff time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		log.Printf("Attempt %d: Failed to connect to the database. Retrying in %s...\n", i+1, backoff)
		time.Sleep(backoff)

		backoff *= 2
	}

	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}

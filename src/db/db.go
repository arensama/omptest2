package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Db *gorm.DB
}

func Init() *DB {

	// Connect to the database
	dsn := "host=localhost user=postgres password=postgres dbname=testapi2 port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return &DB{
		Db: db,
	}
}
func (db *DB) Migrate(model interface{}) error {
	err := db.Db.AutoMigrate(model)
	if err != nil {
		log.Fatalf("cant update model %v", err)
	}
	return err
}

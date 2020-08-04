package config

import (
	"log"
	"os"

	"github.com/amonoyflow/mongodb-go/controllers"
	"github.com/go-pg/pg/v9"
)

// Connect function
func Connect() *pg.DB {
	opts := &pg.Options{
		Database: "<dbName>",
		User:     "<user>",
		Password: "<password>",
		Addr:     "localhost:5432",
	}

	var db *pg.DB = pg.Connect(opts)

	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}

	log.Printf("Connected to db")
	controllers.CreateTodoTable(db)
	controllers.InitiateDB(db)
	return db
}

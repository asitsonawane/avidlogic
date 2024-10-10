package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var DB *pgx.Conn

// ConnectDB initializes the connection to PostgreSQL
func ConnectDB() {
	var err error
	DB, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to the database successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	DB.Close(context.Background())
}

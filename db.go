package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var DB *pgx.Conn

func ConnectDB() {
	var err error
	DB, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to the database successfully")
}

func CloseDB() {
	DB.Close(context.Background())
}

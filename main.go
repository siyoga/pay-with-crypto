package main

import (
	"flag"
	"log"
	"os"
	"pay-with-crypto/app"
	db "pay-with-crypto/app/datastore"

	"github.com/joho/godotenv"
)

func main() {
	prod := flag.Bool("p", false, "Flag for production run")
	flag.Parse()

	if *prod {
		if err := godotenv.Load("prod.env"); err != nil {
			log.Fatalf("Error loading prod.env file")
		}
	} else {
		if err := godotenv.Load("dev.env"); err != nil {
			log.Fatalf("Error loading dev.env file")
		}
	}

	config := db.DatabaseConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
		Host:     os.Getenv("DATABASE_HOST"),
	}

	server := app.Start(config)

	if err := server.Listen(":8081"); err != nil {
		log.Panic(err)
	}
}

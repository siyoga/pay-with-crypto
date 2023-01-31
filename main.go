package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pay-with-crypto/app"
	"pay-with-crypto/app/datastore"

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

	config := DatabaseConfig{
		User: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
		Host: os.Getenv("DATABASE_HOST"),
	}

	datastore.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	app.Start("8081")
}
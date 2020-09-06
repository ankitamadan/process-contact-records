package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

//GetDatabaseString Get the database string
func GetDatabaseString() string {

	var err error

	currentPath, _ := os.Getwd()
	dotEnvFile := fmt.Sprintf("%s/.env", currentPath)
	err = godotenv.Load(dotEnvFile)
	if err != nil {
		log.Printf("skipped loading .env file, file not present at location: %s", dotEnvFile)
	} else {
		log.Printf("loading .env file at location: %s", dotEnvFile)
	}

	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")

	log.Printf("connecting to postgres database name: %s at host: %s with user: %s", dbName, dbHost, dbUser)

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("connected to postgres database: %s", psqlInfo)
	return psqlInfo
}

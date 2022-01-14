package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/exec"
)
var count string

func main() {
	log.Println("Starting scooter Initiation")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", "order_micro", err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT COUNT(id) FROM scooters ")
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("There are %v scooters into Database\n", count)

	scootersInit()
}

func scootersInit() {
	fmt.Println("scooter containers creating...")
	command := exec.Command("docker-compose", "up", "-d", "--scale", "scooter_service=" + count)
	fmt.Println(command.String())
	err := command.Run()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("scooters were created")
}

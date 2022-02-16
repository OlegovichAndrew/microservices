package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
)

var count string

func main() {
	log.Println("Starting scooter Initiation")
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		POSTGRES_USER,
		POSTGRES_PASSWORD,
		PG_HOST,
		PG_PORT,
		POSTGRES_DB)

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

	fmt.Printf("There are %v scooters into the Database\n", count)

	scootersInit()
}

func scootersInit() {
	fmt.Println("scooter containers creating...")
	command := exec.Command("docker-compose", "up", "-d", "--scale", "scooter_client="+count)
	fmt.Println(command.String())
	err := command.Run()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("scooters were created")
}

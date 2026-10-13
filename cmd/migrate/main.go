package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"olx_api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage : make migration < up | down >")
	}
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal("migration.New() Error : ", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal("migration.Up() Error : ", err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal("migration.Down() Error : ", err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: make migrate-force version=<version>")
		}
		var version int
		if _, err := fmt.Sscanf(os.Args[2], "%d", &version); err != nil {
			log.Fatal("invalid version: ", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatal("migration.Force() Error : ", err)
		}
		fmt.Printf("forced version to %d successfully\n", version)
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatal("migration.Drop() Error : ", err)
		}
		fmt.Println("migrations dropped successfully")
	default:
		log.Fatalf("Some other command there")
	}

}

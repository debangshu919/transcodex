package main

import (
	"log"
	"os"

	"github.com/debangshu919/transcodex/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	switch os.Args[1] {
	case "up":
		log.Println("migrating up...")
		if err := m.Up(); err != nil {
			log.Fatalf("failed to migrate up: %v", err)
		}
		log.Println("migration up completed")

	case "down":
		log.Println("migrating down...")
		if err := m.Down(); err != nil {
			log.Fatalf("failed to migrate down: %v", err)
		}
		log.Println("migration down completed")

	default:
		log.Fatal("usage: migrate <up|down>")
	}
}

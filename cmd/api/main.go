package main

import (
	"log"
	"net/http"
	"time"

	"github.com/debangshu919/transcodex/internal/config"
	"github.com/debangshu919/transcodex/internal/db"
	router "github.com/debangshu919/transcodex/internal/http"
	"github.com/debangshu919/transcodex/internal/storage"
	logger "github.com/debangshu919/transcodex/internal/utils"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.NewLogger(cfg, "logs")

	log.Println("Starting server in", "env", cfg.Env, "mode")

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
	}
	defer database.Close()
	log.Println("Connected to database")

	store := storage.NewS3Storage(cfg)
	log.Println("S3 storage initialized")

	handler := router.HttpHandler(database, store, logger)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Server running at http://localhost" + srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

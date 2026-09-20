package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	mw "github.com/debangshu919/transcodex/internal/api/middlewares"
	"github.com/debangshu919/transcodex/internal/config"
	"github.com/debangshu919/transcodex/internal/db"
	router "github.com/debangshu919/transcodex/internal/http"
	"github.com/debangshu919/transcodex/internal/storage"
	logger "github.com/debangshu919/transcodex/internal/utils"
	"golang.org/x/net/http2"
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

	// Load certificate
	cert := "cert.pem"
	key := "key.pem"

	// Configure TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hdlr := mw.Cors(
		mw.SecurityHeaders(
			mw.Compression(handler),
		),
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      hdlr,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
		TLSConfig:    tlsConfig,
	}

	// Enable HTTP2
	http2.ConfigureServer(srv, &http2.Server{})

	log.Println("Server running at https://localhost" + srv.Addr)
	if err := srv.ListenAndServeTLS(cert, key); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

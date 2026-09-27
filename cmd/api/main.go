package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	mw "github.com/debangshu919/transcodex/internal/api/middlewares"
	router "github.com/debangshu919/transcodex/internal/api/router"
	"github.com/debangshu919/transcodex/internal/config"
	"github.com/debangshu919/transcodex/internal/db"
	"github.com/debangshu919/transcodex/internal/storage"
	"github.com/debangshu919/transcodex/pkg/logger"
	"github.com/debangshu919/transcodex/pkg/utils"
	"golang.org/x/net/http2"
)

func main() {
	var env string

	if len(os.Args) < 2 {
		log.Fatal("Environment not specified! Defaulting to development")
		env = "development"
	} else if os.Args[1] != "production" {
		env = "development"
	} else {
		env = "production"
	}

	cfg := config.MustLoad(env)
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

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"allowedParam"},
	}
	hpp := mw.Hpp(hppOptions)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      utils.ApplyMiddleware(handler, hpp, mw.Compression, mw.SecurityHeaders, rl.Middleware, mw.Cors),
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

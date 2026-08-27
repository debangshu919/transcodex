package main

import (
	"log"
	"net/http"
	"time"

	"github.com/debangshu919/transcodex/internal/config"
	router "github.com/debangshu919/transcodex/internal/http"
)

func main() {
	cfg := config.MustLoad()

	handler := router.HttpHandler()

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Println("Starting server in http://localhost" + srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

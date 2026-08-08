package main

import (
	"errors"
	"log"
	"net/http"

	"learn-network-go/learn_web_service/internal/amf/server"
	"learn-network-go/learn_web_service/internal/config"
)

func main() {
	cfg, err := config.Load("learn_web_service/config.amf.json")
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewAMFServer(cfg.ServerAddress)

	log.Printf("AMF server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

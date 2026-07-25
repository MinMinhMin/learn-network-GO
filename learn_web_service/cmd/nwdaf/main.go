package main

import (
	"errors"
	"log"
	"net/http"

	"learn-network-go/learn_web_service/internal/server"
)

func main() {
	srv := server.NewNWDAFServer(":8080")

	log.Printf("NWDAF API listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

package server

import (
	"net/http"
	"time"

	"learn-network-go/learn_web_service/internal/amf/handler"
)

func NewAMFServer(addr string) *http.Server {
	notifyHandler := handler.NewNotifyHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		notifyHandler.Notify(w, r)
	})

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

package server

import (
	"net/http"
	"time"

	"learn-network-go/learn_web_service/internal/handler"
	"learn-network-go/learn_web_service/internal/repository"
	"learn-network-go/learn_web_service/internal/service"
)

func NewNWDAFServer(addr string) *http.Server {
	subscriptionRepo := repository.NewSubscriptionRepository()
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	mux := http.NewServeMux()

	mux.Handle("/subscriptions", handler.Methods{
		http.MethodGet:  http.HandlerFunc(subscriptionHandler.ListSubscriptions),
		http.MethodPost: http.HandlerFunc(subscriptionHandler.Subscribe),
	})

	mux.Handle("/subscriptions/", handler.Methods{
		http.MethodDelete: http.HandlerFunc(subscriptionHandler.Unsubscribe),
	})

	mux.Handle("/analytics", handler.Methods{
		http.MethodGet: http.HandlerFunc(subscriptionHandler.GetAnalytics),
	})

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

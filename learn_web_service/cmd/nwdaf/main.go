package main

import (
	"context"
	"errors"
	"learn-network-go/learn_web_service/internal/config"
	"learn-network-go/learn_web_service/internal/nwdaf/repository"
	"learn-network-go/learn_web_service/internal/nwdaf/server"
	"learn-network-go/learn_web_service/internal/nwdaf/service"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {

	// load config
	cfg, err := config.Load("learn_web_service/config.nwdaf.json")
	if err != nil {
		log.Fatal(err)
	}

	// set limit time for ping mongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))

	if err != nil {
		log.Fatal(err)
	}

	if client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	defer func() {
		// old context is for ping the client at start (maybe is already met the deadline) so we need the newone for shutdown the connection (db)
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(disconnectCtx)
	}()

	// create repo
	collection := client.Database(cfg.Database).Collection(cfg.Collection)

	subscriptionRepo := repository.NewSubscriptionRepository(collection)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// notify
	notifyScheduler := service.NewNotifyScheduler(subscriptionService)
	if err := notifyScheduler.Start(cfg.NotifySchedule); err != nil {
		log.Fatal(err)
	}
	defer notifyScheduler.Stop()

	// start server
	srv := server.NewNWDAFServer(cfg.ServerAddress, subscriptionService)

	log.Printf("NWDAF API listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

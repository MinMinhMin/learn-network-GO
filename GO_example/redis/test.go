package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"errors"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	MongoURI      string `json:"mongo_uri"`
	Database      string `json:"database"`
	Collection    string `json:"collection"`
	ServerAddress string `json:"server_address"`
}

type Sample struct {
	Name string `bson:"name" json:"name"`
	Info string `bson:"info" json:"info"`
}

func main() {
	_ = godotenv.Load("GO_example/redis/.env")

	data, err := os.ReadFile("GO_example/redis/config.json")

	if err != nil {
		log.Fatal(err)
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}
	pass := os.Getenv("gocluster_pass")
	cfg.MongoURI = strings.ReplaceAll(cfg.MongoURI, "<gocluster_pass>", pass)

	ctx := context.Background()
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))

	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	collection := mongoClient.Database(cfg.Database).Collection(cfg.Collection)

	log.Println("MongoDb Connected")

	// make a sample in db if it dont exist

	test_sample := Sample{
		Name: "minmin",
		Info: "VNU student ",
	}

	_, err = collection.InsertOne(ctx, test_sample)

	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Redis connected")

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		if name == "" {
			log.Fatal("name is empty")
		}

		log.Printf("REDIS GET %s,", name)

		cached, err := rdb.Get(r.Context(), name).Result()

		if err == nil {
			log.Println("CACHE HIT")

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Source", "redis")
			fmt.Fprint(w, cached)
			return
		}

		if !errors.Is(err, redis.Nil) {
			log.Printf("Redis error: %v", err)
		}

		log.Println("CACHE MISS")

		log.Printf("MONGODB FIND name=%s", name)

		var sample Sample
		err = collection.FindOne(r.Context(), bson.M{"name": name}).Decode(&sample)

		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Cannot find sample", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(sample)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("REDIS SET %s TLS=30s,", name)

		if err := rdb.Set(r.Context(), name, result, 30*time.Second).Err(); err != nil {
			log.Printf("REDIS SET error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Source", "mongodb")
		w.Write(result)

	})

	log.Printf("Server run at http://localhost%s", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatal(err)
	}
}

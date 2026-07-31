package ch9buildinghttpservices

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	MongoURI   string `json:"mongo_uri"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type User struct {
	ID    bson.ObjectID `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Email string        `bson:"email"`
}

func loadConfig(t *testing.T, filename string) Config {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Cannot read config file: %v", err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Config JSON invalid: %v", err)
	}

	if config.MongoURI == "" ||
		config.Database == "" ||
		config.Collection == "" {
		t.Fatal("Config's mongo_uri, database or collection key is missing")
	}

	return config
}

func TestMongoDB(t *testing.T) {
	config := loadConfig(t, "./db_config/config.json")

	client, err := mongo.Connect(
		options.Client().ApplyURI(config.MongoURI),
	)
	if err != nil {
		t.Fatalf("Cannot create MongoDB client: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			t.Logf("Disconnection error MongoDB: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("Cant connect to MongoDB: %v", err)
	}

	t.Logf("connect to MongoDB succeed")

	collection := client.
		Database(config.Database).
		Collection(config.Collection)

	testEmail := "test-" + time.Now().Format("20060102150405") + "@example.com"

	user := User{
		Name:  "Test User",
		Email: testEmail,
	}

	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		t.Fatalf("Cannot add user: %v", err)
	}

	t.Logf("User added, ID: %v", insertResult.InsertedID)

	// Delete added data when test close.
	// defer func() {
	// 	_, err := collection.DeleteOne(
	// 		context.Background(),
	// 		bson.M{"email": testEmail},
	// 	)

	// 	if err != nil {
	// 		t.Logf("Cannot delete test data": %v", err)
	// 	}
	// }()

	var savedUser User

	err = collection.FindOne(
		ctx,
		bson.M{"email": testEmail},
	).Decode(&savedUser)

	if err != nil {
		t.Fatalf("Cannot read user from MongoDB: %v", err)
	}

	t.Logf(
		"Read User succeed: ID=%s, Name=%s, Email=%s",
		savedUser.ID.Hex(),
		savedUser.Name,
		savedUser.Email,
	)

	if savedUser.Name != user.Name {
		t.Fatalf(
			"Incorrect name: expect %q, actual %q",
			user.Name,
			savedUser.Name,
		)
	}

	if savedUser.Email != user.Email {
		t.Fatalf(
			"Incorrect email: expect %q, actual %q",
			user.Email,
			savedUser.Email,
		)
	}
}

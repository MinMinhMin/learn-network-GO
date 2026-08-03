package repository

import (
	"context"
	"errors"
	"learn-network-go/learn_web_service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// mongoDB already have method to prevent race condition. To limit the time for each db interaction, we will use context with timeout

var ErrSubscriptionNotFound = errors.New("subscription ID not found")

type SubscriptionRepository struct {
	collection *mongo.Collection
}

func NewSubscriptionRepository(collection *mongo.Collection) *SubscriptionRepository {
	return &SubscriptionRepository{
		collection: collection,
	}
}

func (r *SubscriptionRepository) Save(sub model.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert a subscription using InsertOne
	_, err := r.collection.InsertOne(ctx, sub)
	return err
}

func (r *SubscriptionRepository) FindByID(subscriptionID string) (model.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sub model.Subscription

	// find a subscription by its id with FindOne using bson.M (map)
	err := r.collection.FindOne(ctx, bson.M{
		"_id": subscriptionID,
	}).Decode(&sub)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Subscription{}, ErrSubscriptionNotFound
	}

	return sub, err
}

func (r *SubscriptionRepository) ListActive() ([]model.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// using mongo.Collection.Find will return a cursor (iter) to iterate all the docs that match the filter
	// mongo.Cursor have .Next(), .Close(), .All()
	cursor, err := r.collection.Find(ctx, bson.M{
		"active": true,
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var result []model.Subscription

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SubscriptionRepository) Deactivate(subscriptionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// $set is set method
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": subscriptionID},
		bson.M{"$set": bson.M{"active": false}},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}

package repository

import (
	"errors"
	"learn-network-go/learn_web_service/internal/model"
	"sync"
)

var ErrSubscriptionNotFound = errors.New("subscription ID not found")

type SubscriptionRepository struct {
	mu            sync.RWMutex
	subscriptions map[string]model.Subscription
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		subscriptions: make(map[string]model.Subscription),
	}
}

func (r *SubscriptionRepository) Save(sub model.Subscription) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subscriptions[sub.SubscriptionID] = sub
}

func (r *SubscriptionRepository) FindByID(subscriptionID string) (model.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sub, ok := r.subscriptions[subscriptionID]

	if !ok {
		return model.Subscription{}, ErrSubscriptionNotFound
	}

	return sub, nil
}

func (r *SubscriptionRepository) ListActive() []model.Subscription {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Subscription

	for _, sub := range r.subscriptions {
		if sub.Active {
			result = append(result, sub)
		}
	}

	return result
}

func (r *SubscriptionRepository) Deactivate(subscriptionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sub, ok := r.subscriptions[subscriptionID]
	if !ok {
		return ErrSubscriptionNotFound
	}

	sub.Active = false
	r.subscriptions[subscriptionID] = sub

	return nil
}

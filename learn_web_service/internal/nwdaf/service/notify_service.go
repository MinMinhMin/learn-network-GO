package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"learn-network-go/learn_web_service/internal/model"

	"github.com/robfig/cron/v3"
)

type NotifyScheduler struct {
	cron                *cron.Cron
	subscriptionService *SubscriptionService
	client              *http.Client
}

func NewNotifyScheduler(subscriptionService *SubscriptionService) *NotifyScheduler {
	return &NotifyScheduler{
		cron:                cron.New(cron.WithSeconds()),
		subscriptionService: subscriptionService,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *NotifyScheduler) Start(schedule string) error {
	_, err := s.cron.AddFunc(schedule, s.NotifyActiveSubscriptions)
	if err != nil {
		return err
	}

	s.cron.Start()
	return nil
}

func (s *NotifyScheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}

func (s *NotifyScheduler) NotifyActiveSubscriptions() {
	subscriptions, err := s.subscriptionService.ListActiveSubscriptions()
	if err != nil {
		log.Printf("list active subscriptions failed: %v", err)
		return
	}

	for _, sub := range subscriptions {
		if err := s.sendNotify(sub); err != nil {
			log.Printf("send notify failed: subscriptionId=%s err=%v", sub.SubscriptionID, err)
		}
	}
}

func (s *NotifyScheduler) sendNotify(sub model.Subscription) error {
	report := s.subscriptionService.BuildAnalyticsReport(sub)

	notifyReq := model.NotifyRequest{
		SubscriptionID: sub.SubscriptionID,
		Report:         report,
	}

	body, err := json.Marshal(notifyReq)
	if err != nil {
		return err
	}

	resp, err := s.client.Post(sub.CallbackURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Printf("notify returned non-2xx: subscriptionId=%s status=%s", sub.SubscriptionID, resp.Status)
	}

	return nil
}

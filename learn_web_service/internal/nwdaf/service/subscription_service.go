package service

import (
	"errors"
	"fmt"
	"learn-network-go/learn_web_service/internal/model"
	"learn-network-go/learn_web_service/internal/nwdaf/repository"
	"strings"
	"sync/atomic"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
	// using atomic to make subscriptionID - prevent from data race when calculate ID
	nextID atomic.Uint64
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func isValidConsumerType(consumerType model.NFType) bool {
	switch consumerType {
	case model.NFTypeAMF, model.NFTypePCF, model.NFTypeSMF:
		return true
	default:
		return false
	}
}

func isAllowedAnalyticsID(consumerType model.NFType, analyticsID string) bool {
	switch consumerType {
	case model.NFTypeAMF:
		return analyticsID == "UE Mobility"
	case model.NFTypePCF:
		return analyticsID == "NF load information"
	case model.NFTypeSMF:
		return analyticsID == "PDU Session traffic"
	default:
		return false
	}
}

var invalidAnalyticsRequest = errors.New("Invalid analytics request")

func (s *SubscriptionService) Subscribe(req model.SubscribeRequest) (model.SubscribeResponse, error) {
	if strings.TrimSpace(req.ConsumerID) == "" || strings.TrimSpace(req.AnalyticsID) == "" || strings.TrimSpace(req.CallbackURL) == "" {
		return model.SubscribeResponse{}, errors.New("Invalid subscribe request")
	}

	if !isValidConsumerType(req.ConsumerType) {
		return model.SubscribeResponse{}, errors.New("Invalid subscribe request")
	}

	if !isAllowedAnalyticsID(req.ConsumerType, req.AnalyticsID) {
		return model.SubscribeResponse{}, errors.New("Invalid subscribe request")
	}

	subscriptionID := fmt.Sprintf("sub-%d", s.nextID.Add(1))

	sub := model.Subscription{
		SubscribeRequest: req,
		SubscriptionID:   subscriptionID,
		Active:           true,
	}

	if err := s.repo.Save(sub); err != nil {
		return model.SubscribeResponse{}, err
	}

	return model.SubscribeResponse{
		SubscriptionID: subscriptionID,
		Message:        "subscription created",
	}, nil
}

func (s *SubscriptionService) Unsubscribe(subscriptionID string) error {
	if strings.TrimSpace(subscriptionID) == "" {
		return errors.New("Invalid subscribe request")
	}

	return s.repo.Deactivate(subscriptionID)
}

func (s *SubscriptionService) ListActiveSubscriptions() ([]model.Subscription, error) {
	return s.repo.ListActive()
}

func (s *SubscriptionService) GetAnalytics(consumerType model.NFType, analyticsID string) (model.AnalyticsReport, error) {
	if !isValidConsumerType(consumerType) {
		return model.AnalyticsReport{}, invalidAnalyticsRequest
	}

	if strings.TrimSpace(analyticsID) == "" || !isAllowedAnalyticsID(consumerType, analyticsID) {
		return model.AnalyticsReport{}, invalidAnalyticsRequest
	}

	return s.BuildAnalyticsReport(model.Subscription{
		SubscribeRequest: model.SubscribeRequest{
			ConsumerType: consumerType,
			AnalyticsID:  analyticsID,
		},
		Active: true,
	}), nil
}

func (s *SubscriptionService) BuildAnalyticsReport(sub model.Subscription) model.AnalyticsReport {
	analyticsID := sub.AnalyticsID

	switch sub.ConsumerType {
	case model.NFTypeAMF:
		return model.AnalyticsReport{
			AnalyticsID: analyticsID,
			Value:       "AMF received UE mobility analytics",
		}
	case model.NFTypePCF:
		return model.AnalyticsReport{
			AnalyticsID: analyticsID,
			Value:       "PCF received policy analytics",
		}
	case model.NFTypeSMF:
		return model.AnalyticsReport{
			AnalyticsID: analyticsID,
			Value:       "SMF received session analytics",
		}
	default:
		return model.AnalyticsReport{
			AnalyticsID: analyticsID,
			Value:       "Unknown consumer analytics",
		}
	}
}

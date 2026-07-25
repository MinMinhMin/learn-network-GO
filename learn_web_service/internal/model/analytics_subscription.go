package model

type NFType string

const (
	NFTypeAMF NFType = "AMF"
	NFTypePCF NFType = "PCF"
	NFTypeSMF NFType = "SMF"
)

type SubscribeRequest struct {
	ConsumerID   string `json:"consumerId"`
	ConsumerType NFType `json:"consumerType"`
	AnalyticsID  string `json:"analyticsId"`
}

type SubscribeResponse struct {
	SubscriptionID string `json:"subscriptionId"`
	Message        string `json:"message"`
}

type Subscription struct {
	SubscribeRequest
	SubscriptionID string `json:"subscriptionId"`
	Active         bool   `json:"active"`
}

type AnalyticsReport struct {
	AnalyticsID string `json:"analyticsId"`
	Value       string `json:"value"`
}

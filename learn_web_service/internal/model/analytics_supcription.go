package model

type SubscribeRequest struct {
	ConsumerID      string   `json:"consumerId"`
	AnalyticsIDs    []string `json:"analyticsIds"`
	CallbackURL     string   `json:"callbackUrl"`
	IntervalSeconds int      `json:"intervalSeconds"`
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

type NotifyRequest struct {
	SubscriptionID string          `json:"subscriptionId"`
	Report         AnalyticsReport `json:"report"`
}

type AnalyticsReport struct {
	AnalyticsID string `json:"analyticsId"`
	Value       string `json:"value"`
}

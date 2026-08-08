package model

type NFType string

const (
	NFTypeAMF NFType = "AMF"
	NFTypePCF NFType = "PCF"
	NFTypeSMF NFType = "SMF"
)

// add bson tag for mongo db
type SubscribeRequest struct {
	ConsumerID   string `json:"consumerId" bson:"consumerId" `
	ConsumerType NFType `json:"consumerType" bson:"consumerType"`
	AnalyticsID  string `json:"analyticsId" bson:"analyticsId"`
	CallbackURL  string `json:"callbackUrl" bson:"callbackUrl"`
}

type SubscribeResponse struct {
	SubscriptionID string `json:"subscriptionId"`
	Message        string `json:"message"`
}

// _id is for primary key , we use inline for flatten fields instead of nested loop fields
type Subscription struct {
	SubscribeRequest `bson:",inline"`
	SubscriptionID   string `json:"subscriptionId" bson:"_id"`
	Active           bool   `json:"active" bson:"active"`
}

type AnalyticsReport struct {
	AnalyticsID string `json:"analyticsId"`
	Value       string `json:"value"`
}

type NotifyRequest struct {
	SubscriptionID string          `json:"subscriptionId"`
	Report         AnalyticsReport `json:"report"`
}

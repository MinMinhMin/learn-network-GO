package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"learn-network-go/learn_web_service/internal/model"
)

type NotifyHandler struct{}

func NewNotifyHandler() *NotifyHandler {
	return &NotifyHandler{}
}

func (h *NotifyHandler) Notify(w http.ResponseWriter, r *http.Request) {
	var req model.NotifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid json body",
		})
		return
	}

	log.Printf("AMF received notify: subscriptionId=%s analyticsId=%s value=%s", req.SubscriptionID, req.Report.AnalyticsID, req.Report.Value)

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "notify received",
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

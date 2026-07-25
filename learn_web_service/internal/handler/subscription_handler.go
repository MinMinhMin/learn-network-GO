package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"

	"learn-network-go/learn_web_service/internal/model"
	"learn-network-go/learn_web_service/internal/repository"
	"learn-network-go/learn_web_service/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

// Method routing
type Methods map[string]http.Handler

func (h Methods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(body io.ReadCloser) {
		_, _ = io.Copy(io.Discard, body)
		_ = body.Close()
	}(r.Body)

	if handler, ok := h[r.Method]; ok {
		if handler == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Allow", h.allowMethods())

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (h Methods) allowMethods() string {
	methods := make([]string, 0, len(h))

	for method := range h {
		methods = append(methods, method)
	}

	sort.Strings(methods)
	return strings.Join(methods, ", ")
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req model.SubscribeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid json body",
		})
		return
	}

	resp, err := h.service.Subscribe(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.service.ListActiveSubscriptions())
}

func (h *SubscriptionHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	consumerType := model.NFType(r.URL.Query().Get("consumerType"))
	analyticsID := r.URL.Query().Get("analyticsId")

	report, err := h.service.GetAnalytics(consumerType, analyticsID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, report)
}

func (h *SubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	// Get sub ID from path, if subscriptionID same as path URL -> TrimPrefix failure
	subscriptionID := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	if strings.TrimSpace(subscriptionID) == "" || subscriptionID == r.URL.Path {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": "missing subscription id",
		})
		return
	}

	err := h.service.Unsubscribe(subscriptionID)
	if errors.Is(err, repository.ErrSubscriptionNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "subscription deactivated",
	})
}

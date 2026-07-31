package ch9buildinghttpservices

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// authentication middleware

type contextKey string

const userKey contextKey = "authenticatedUser"

func writeJsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeJsonError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeJsonError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		token := parts[1]

		if token != "secret-token" {
			writeJsonError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userKey, "user-123")
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userKey).(string)
	if !ok {
		writeJsonError(w, http.StatusInternalServerError, "user not found in context")
		return
	}

	w.Header().Set("Content-Type", "application")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "athentication successful",
		"user_id": userID,
	})
}

func TestAuthenticationMiddleware(t *testing.T) {
	serveMux := http.NewServeMux()

	serveMux.Handle(
		"/profile",
		AuthenticationMiddleware(http.HandlerFunc(profileHandler)),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"http://test/profile",
		nil,
	)

	// add token to the request
	req.Header.Set("Authorization", "Bearer secret-token")

	recorder := httptest.NewRecorder()

	serveMux.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status code %d; actual %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}
}

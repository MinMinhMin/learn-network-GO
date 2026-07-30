package ch9buildinghttpservices

// middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeOutMiddleware(t *testing.T) {

	//http.TimeoutHandler is a middleware that accept a handler and return a handler
	handler := http.TimeoutHandler(

		// simulate a handler that take time purposefully to pretend reading client exceed the timeout
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
			time.Sleep(time.Minute)
		}),
		time.Second,
		"Time out while reading response",
	)

	r := httptest.NewRequest(http.MethodGet, "http://test/", nil)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status code: %q", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	_ = resp.Body.Close()
	if actual := string(b); actual != "Timed out while reading response" {
		t.Logf("unexpected body: %q", actual)
	}

}

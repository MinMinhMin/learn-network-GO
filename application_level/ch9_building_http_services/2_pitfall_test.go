package ch9buildinghttpservices

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWriteHeader_test(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		// in net/http when u call write before writeheader -> status default will be 200 OK
		_, _ = w.Write([]byte("Bad Request"))
		w.WriteHeader(http.StatusBadRequest)
	}

	r := httptest.NewRequest(http.MethodGet, "http://test", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	t.Logf("Response status: %q", w.Result().Status)

	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))

	}

	r = httptest.NewRequest(http.MethodGet, "http://test", nil)
	w = httptest.NewRecorder()
	handler(w, r)

	t.Logf("Response status: %q", w.Result().Status)

}

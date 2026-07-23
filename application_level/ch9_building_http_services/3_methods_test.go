package ch9buildinghttpservices

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

type Methods map[string]http.Handler

func (h Methods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	if handler, ok := h[r.Method]; ok {
		if handler == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			handler.ServeHTTP(w, r)
		}
		return
	}

	w.Header().Add("Allow", h.allowMethods())
	if r.Method != http.MethodOptions {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
	}

}

func (h Methods) allowMethods() string {
	a := make([]string, 0, len(h))

	for k := range h {
		a = append(a, k)
	}
	sort.Strings(a)
	return strings.Join(a, ", ")
}

// HTTP method router (a bit of multiplexer)
func DefaultMethodHandler() http.Handler {
	return Methods{
		http.MethodGet: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("Hello, friend!"))
			},
		),

		http.MethodPost: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				_, _ = fmt.Fprintf(w, "Hello, %s!", html.EscapeString(string(b)))
			},
		),
	}
}

func TestMethod(t *testing.T) {
	srv := http.Server{
		Addr:              "127.0.0.1:8081",
		Handler:           http.TimeoutHandler(DefaultMethodHandler(), 2*time.Minute, ""),
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	l, err := net.Listen("tcp", srv.Addr)

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		// this method serve request from the listener
		err := srv.Serve(l)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	testCases := []struct {
		method   string
		body     io.Reader
		code     int
		response string
	}{
		{http.MethodGet, nil, http.StatusOK, "Hello, friend!"},
		{http.MethodPost, bytes.NewBufferString("<world>"), http.StatusOK, "Hello, &lt;world&gt;!"},
		{http.MethodHead, nil, http.StatusMethodNotAllowed, ""},
	}

	client := new(http.Client)
	path := fmt.Sprintf("http://%s/", srv.Addr)

	for i, c := range testCases {
		r, err := http.NewRequest(c.method, path, c.body)

		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		resp, err := client.Do(r)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		if resp.StatusCode != c.code {
			t.Errorf("%d : unexpected status code: %q", i, resp.Status)
		}

		// respond is byte
		b, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("%d :%v", i, err)
			continue
		}
		_ = resp.Body.Close()

		// convert to string for comparing
		if c.response != string(b) {
			t.Errorf("%d: expected %q; actual %q", i, c.response, b)
		}

	}
	if err := srv.Close(); err != nil {
		t.Fatal(err)
	}
}

package ch9buildinghttpservices

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	// handler is provided by the auther, the logic of GET, POST, HEAD are set by author
	"github.com/awoodbeck/gnp/ch09/handlers"
)

// This make a simulate server service, in real scenario, the tcp could be hiden
func TestSimpleHTTPServer(t *testing.T) {
	srv := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: http.TimeoutHandler(handlers.DefaultHandler(), 2*time.Minute, ""),
		// time limit to keep alive a tcp session, if no request provide from client for a time of idletimeout
		IdleTimeout: 5 * time.Minute,
		// time limit of a client to send the request header
		ReadHeaderTimeout: time.Minute,
	}
	// create a listener bound to the server address
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

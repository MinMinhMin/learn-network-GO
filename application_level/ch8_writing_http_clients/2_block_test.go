package ch8writinghttpclients

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
	select {}
}

// commenting cuz this test take really long time

// func TestBlockIndefinitely(t *testing.T) {

// 	// make a test server that have http handler that block Indefinitely
// 	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

// 	// http.Get(ts.URL) will on queue overtime because we dont get any respond
// 	_, _ = http.Get(ts.URL)

// 	// this one will never trigger
// 	t.Fatal("client did not indefinitely block")
// }

func TestBlockIndefinitelyWithTimeout_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

	// make a ctx with timeout, the timer will trigger after initiate
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// make a request with nil body (payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
		return
	}

	// In real case, the context can expire while the response body is being read.

	_ = resp.Body.Close()
}

func TestBlockIndefinitelyWithTimeout_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

	// make a ctx with timeout, the timer will trigger after initiate
	ctx, cancel := context.WithCancel(context.Background())

	timer := time.AfterFunc(5*time.Second, cancel)

	// make a request with nil body (payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
		return
	}

	//We can add more 5 second to read our body in real case
	timer.Reset(5 * time.Second)

	// close body becase we dont use it in this test
	_ = resp.Body.Close()
}

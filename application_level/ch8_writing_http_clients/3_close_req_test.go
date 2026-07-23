package ch8writinghttpclients

import (
	"context"
	"net/http"
	"testing"
)

func TestCloseReq(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://www.time.gov/", nil)

	if err != nil {
		t.Fatal(err)
	}

	//Do not reuse the TCP connection after this request.
	req.Close = true

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	date := resp.Header.Get("Date")

	t.Logf("Date %s", date)

	_ = resp.Body.Close()

}

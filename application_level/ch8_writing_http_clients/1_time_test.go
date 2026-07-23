package ch8writinghttpclients

import (
	"net/http"
	"testing"
	"time"
)

func TestHeadTime(t *testing.T) {

	// Date is a header so we will HEAD the website
	resp, err := http.Head("https://www.time.gov/")

	if err != nil {
		t.Fatal(err)
	}

	// remember to close respond body in every case
	_ = resp.Body.Close()

	now := time.Now().Round(time.Second)

	// get the Date header
	date := resp.Header.Get("Date")

	if date == "" {
		t.Fatal("no Date header received from time.gov")
	}

	dt, err := time.Parse(time.RFC1123, date)

	if err != nil {
		t.Fatal(err)
	}

	// dif from local date to time.gov date
	t.Logf("time.gov: %s (skew %s)", dt, now.Sub(dt))
}

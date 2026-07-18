package ch3reliabletcpdatastream

import (
	"context"
	"io"
	"testing"
	"time"
)

// this simulate ping pong but on I/O Go's write and read (which we can imagine server and client)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration
	select {
	case <-ctx.Done():
		return
	case interval = <-reset:
	default:
	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)

	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case newInterval := <-reset:
			if !timer.Stop() {
				<-timer.C
			}

			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}

		}

		_ = timer.Reset(interval)

	}
}

func TestExamplePinger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r, w := io.Pipe()
	done := make(chan struct{})

	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second

	go func() {
		Pinger(ctx, w, resetTimer)
		close(done) // equivalent to done <- struct{}{}
	}()

	receivePing := func(d time.Duration, r io.Reader) {

		if d >= 0 {
			t.Logf("resetting timer (%s)\n", d)
			resetTimer <- d
		}

		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf)

		if err != nil {
			t.Fatal(err)
		}

		t.Logf("received %q (%s)\n", buf[:n], time.Since(now).Round(100*time.Millisecond))

	}
	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} {
		t.Logf("Run %d:\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}

	cancel()
	<-done
	// ensures the pinger exits after canceling the context
}

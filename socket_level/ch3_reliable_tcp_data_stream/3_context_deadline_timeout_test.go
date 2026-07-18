package ch3reliabletcpdatastream

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	dl := time.Now().Add(1 * time.Second)

	//context.WithDeadline is limit the time of a or mutiple goroutine (same context)

	ctx, cancel := context.WithDeadline(context.Background(), dl)

	defer cancel()

	var d net.Dialer

	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// Sleep enough to meet the context's deadline
		time.Sleep(time.Second + time.Millisecond)
		return nil // make sure our error is timeout
	}

	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")

	if err == nil {
		conn.Close()
		t.Fatal("connection did not timeout")
	}

	nErr, ok := err.(net.Error)

	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a time out: %v", err)
		}
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}

}

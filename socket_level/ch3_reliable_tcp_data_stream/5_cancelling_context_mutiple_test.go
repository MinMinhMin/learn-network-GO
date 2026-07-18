package ch3reliabletcpdatastream

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)

	listener, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	go func() {
		conn, err := listener.Accept()

		if err == nil {
			conn.Close()
		}
	}()

	// use wg *sync.WaitGroup to make sure we can test cancel when all 10 goroutine has been active
	dial := func(ctx context.Context, address string, respond chan int, id int, wg *sync.WaitGroup) {
		defer wg.Done()

		var d net.Dialer

		c, err := d.DialContext(ctx, "tcp", address)

		if err != nil {
			return
		}

		c.Close()

		select {
		// trigger when id not been choosen
		case <-ctx.Done():
		// trigger when id is the first pick
		case respond <- id:
		}
	}

	res := make(chan int)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res
	//use cancel() to send cancel to all gorountine with same context
	cancel()

	//use this to make sure all 10 go routine are canceled
	wg.Wait()

	// this run after wg.Wait() to make sure we know what id is strigger first, to prevent from panic (send on close channel)
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s", ctx.Err())
	}

	t.Logf("dialer %d retrieve the resource", response)

}

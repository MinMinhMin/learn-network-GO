package ch3reliabletcpdatastream

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		defer func() {
			conn.Close()
			close(sync)
		}()

		// setDeadline has been enable in conn (the change not only assign value for error but also set deadline method for conn)

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))

		if err != nil {
			t.Log(err)
			return
		}

		buf := make([]byte, 1)

		//normally Read or Write is call blocking - should using inside a goroutine, but in this we dont use it to simulate the deadline
		_, err = conn.Read(buf)

		nErr, ok := err.(net.Error)

		if !ok || !nErr.Timeout() {
			t.Errorf("expected timeout error; actual: %v", err)
		}

		// send done to sync, after this, our read should work normally
		sync <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))

		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)

		if err != nil {
			t.Error(err)
		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String())

	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	<-sync
	// work only after receiver from sync
	_, err = conn.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)

	_, err = conn.Read(buf)

	// cuz our server wont respond anything so our read from client didnt get anything
	if err != io.EOF {
		t.Errorf("expected server temination; actual %v", err)
	}

}

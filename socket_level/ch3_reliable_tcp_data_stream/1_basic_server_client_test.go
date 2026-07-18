package ch3reliabletcpdatastream

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // basic server

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} /* empty struct, usually use cuz it nearly didnt cost*/ }()

		//this loop of get connection is reason for this goroutine
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024) // allocate slice of byte with 1024 byte (each index is 1 byte)

				// this loop of reading message from connect is the reason for this goroutine
				for {
					n, err := c.Read(buf) // if the message is "hello - 5 letter -> 5 byte -> n=5"

					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n]) // we use :n to make sure that is real data
				}

			}(conn)

		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String()) // basic client

	if err != nil {
		t.Fatal(err)
	}

	message := "hello server"

	_, err = conn.Write([]byte(message))

	if err != nil {
		t.Fatal(err)
	}

	// close from inside then to outside
	conn.Close()
	<-done

	listener.Close()
	<-done

}

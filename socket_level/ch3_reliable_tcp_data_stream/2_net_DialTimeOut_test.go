package ch3reliabletcpdatastream

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, controlAddress string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        controlAddress,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},

		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {

	// simlate net.DialTimeout to make sure it gonna return timeout error
	// if error is not timeout we will know what it is
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)

	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}

	nErr, ok := err.(net.Error) // check if error is implemented of net.Error
	if !ok {
		t.Fatal(err)
	}

	if !nErr.Timeout() { // true if it's a time out
		t.Fatal("error is not a timeout")
	}

}

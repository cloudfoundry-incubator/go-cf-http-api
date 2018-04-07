package testhelpers

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"
)

func GetOpenPort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "0", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "0", err
	}
	defer l.Close()

	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}

func PollForUp(port string) error {
	for {
		select {
		case <-time.After(1 * time.Second):
			return errors.New("timed out waiting for eva api to come up")
		default:
			_, err := http.Get("http://localhost:" + port)
			if err == nil {
				return nil
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

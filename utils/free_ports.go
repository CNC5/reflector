package utils

import (
	"fmt"
	"net"
)

func IsPortBindable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func FindFreePorts(n int) ([]int, error) {
	ports := make([]int, 0, n)

	for i := 0; i < n; i++ {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, err
		}
		defer l.Close()

		addr := l.Addr().(*net.TCPAddr)
		ports = append(ports, addr.Port)
	}
	return ports, nil
}

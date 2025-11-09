package utils

import (
	"net"
)

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

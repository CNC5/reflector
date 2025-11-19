package utils

import (
	"context"
	"net"
	"reflector/log"
)

func IsLocalIPAddress(address net.IPAddr) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.GetDefaultLogger().
			Error().
			Update("err", err.Error()).
			Msg("failed to get interface addresses")
		return false
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			if ipNet.IP.Equal(address.IP) {
				return true
			}
		}
	}
	return false
}

func NSLookup(fqdn string) ([]net.IPAddr, error) {
	return net.DefaultResolver.LookupIPAddr(context.Background(), fqdn)
}

func IsDomainPointingToThisHost(fqdn string) bool {
	addresses, err := NSLookup(fqdn)
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("fqdn", fqdn).
			Update("err", err.Error()).
			Msg("failed to lookup domain name")
		return false
	}
	for _, address := range addresses {
		if IsLocalIPAddress(address) {
			return true
		}
	}
	return false
}

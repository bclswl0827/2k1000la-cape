package main

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

func getIPv4Addrs() (map[string]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ipv4Addresses := make(map[string]string)
	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addrs {
			ipNet, ok := address.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil && len(ipNet.IP.String()) > 0 && !ipNet.IP.IsLoopback() {
				ipv4Addresses[inter.Name] = ipNet.IP.String()
			}
		}
	}

	if len(ipv4Addresses) == 0 {
		return nil, errors.New("no available IPv4 addresses found")
	}

	return ipv4Addresses, nil
}

func getInterfaceByPattern(pattern string, fuzzy bool) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if fuzzy {
			r := regexp.MustCompile(pattern)
			if r.MatchString(iface.Name) {
				return iface.Name, nil
			}
		} else {
			if iface.Name == pattern {
				return iface.Name, nil
			}
		}
	}

	return "", fmt.Errorf("interface %s not found", pattern)
}

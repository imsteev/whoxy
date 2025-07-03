package main

import (
	"fmt"
	"net"
)

func getGlobalIPv6Addresses() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ip := ipNet.IP
				// Check if it's IPv6 and not link-local
				if ip.To4() == nil && ip.To16() != nil && !ip.IsLinkLocalUnicast() {
					fmt.Println(ip.String())
				}
			}
		}
	}
}

func main() {
	getGlobalIPv6Addresses()
}

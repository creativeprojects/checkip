package main

import (
	"net/netip"
	"strings"
)

func parseIP(ip string) string {
	var addr netip.Addr
	switch {
	case strings.Count(ip, ":") == 1 || strings.Count(ip, "]:") == 1:
		addrPort, _ := netip.ParseAddrPort(ip)
		addr = addrPort.Addr()
	default:
		addr, _ = netip.ParseAddr(ip)
	}
	return addr.String()
}

package helper

import (
	"net"
	"runtime"
	"strings"
)

func GetIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipaddr, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if ipaddr.IsLoopback() {
			continue
		}
		if ipaddr.To4() != nil {
			if runtime.GOOS == "darwin" {
				if !strings.HasPrefix(ipaddr.String(), "192") {
					continue
				}
			}
			return ipaddr.String()
		}
	}
	return ""
}

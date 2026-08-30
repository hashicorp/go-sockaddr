//go:build !linux
// +build !linux

package sockaddr

import (
	"net"
)

func NetInterfaces() ([]net.Interface, error) {
	return net.Interfaces()
}

func NewIfAddr(addr IPAddr, intf net.Interface) IfAddr {
	return IfAddr{addr, intf}
}

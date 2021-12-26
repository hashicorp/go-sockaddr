package sockaddr

import (
	"fmt"
	"net"
)

type OSNetProvider struct{}

// GetAllInterfaces iterates over all available network interfaces and finds all
// available IP addresses on each interface and converts them to
// sockaddr.IPAddrs, and returning the result as an array of IfAddr.
func (n *OSNetProvider) GetAllInterfaces() (IfAddrs, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ifAddrs := make(IfAddrs, 0, len(ifs))
	for _, intf := range ifs {
		addrs, err := intf.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ipAddr IPAddr
			ipAddr, err = NewIPAddr(addr.String())
			if err != nil {
				return IfAddrs{}, fmt.Errorf("unable to create an IP address from %q", addr.String())
			}

			ifAddr := IfAddr{
				SockAddr:  ipAddr,
				Interface: intf,
			}
			ifAddrs = append(ifAddrs, ifAddr)
		}
	}

	return ifAddrs, nil
}

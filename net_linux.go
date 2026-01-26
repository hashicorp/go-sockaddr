package sockaddr

import (
	"net"
	"os"
	"syscall"
	"unsafe"
)

// LinuxNetInterface contains a net.Interface along with its associated []net.Addr.
type LinuxNetInterface struct {
	net.Interface
	addrs []net.Addr
}

func (ifi *LinuxNetInterface) Addrs() ([]net.Addr, error) {
	return ifi.addrs, nil
}

// NetInterfaces returns all net.Interfaces along with their associated []net.Addr.
// This Linux optimization avoids a separate Netlink dump of addresses for each individual interface,
// which is prohibitively slow on servers with large numbers of interfaces:
// https://github.com/golang/go/issues/53660
func NetInterfaces() ([]LinuxNetInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	netInterfaces := make([]LinuxNetInterface, 0, len(interfaces))
	for _, ifi := range interfaces {
		netInterfaces = append(netInterfaces, LinuxNetInterface{ifi, make([]net.Addr, 0)})
	}
	ifMap := make(map[int]*LinuxNetInterface, len(netInterfaces))
	for i, ifi := range netInterfaces {
		ifMap[ifi.Index] = &netInterfaces[i]
	}

	tab, err := syscall.NetlinkRIB(syscall.RTM_GETADDR, syscall.AF_UNSPEC)
	if err != nil {
		return nil, os.NewSyscallError("NetlinkRIB", err)
	}
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, os.NewSyscallError("ParseNetLinkMessage", err)
	}

	for _, m := range msgs {
		if m.Header.Type == syscall.RTM_NEWADDR {
			ifam := (*syscall.IfAddrmsg)(unsafe.Pointer(&m.Data[0]))
			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
			if err != nil {
				return nil, os.NewSyscallError("ParseNetLinkRouteAttr", err)
			}
			if ifi, ok := ifMap[int(ifam.Index)]; ok {
				ifi.addrs = append(ifi.addrs, newAddr(ifam, attrs))
			}
		}
	}

	return netInterfaces, err
}

// Vendored unexported function:
// https://github.com/golang/go/blob/8bcc490667d4dd44c633c536dd463bbec0a3838f/src/net/interface_linux.go#L178-L203
func newAddr(ifam *syscall.IfAddrmsg, attrs []syscall.NetlinkRouteAttr) net.Addr {
	var ipPointToPoint bool
	for _, a := range attrs {
		if a.Attr.Type == syscall.IFA_LOCAL {
			ipPointToPoint = true
			break
		}
	}
	for _, a := range attrs {
		if ipPointToPoint && a.Attr.Type == syscall.IFA_ADDRESS {
			continue
		}
		switch ifam.Family {
		case syscall.AF_INET:
			return &net.IPNet{IP: net.IPv4(a.Value[0], a.Value[1], a.Value[2], a.Value[3]), Mask: net.CIDRMask(int(ifam.Prefixlen), 8*net.IPv4len)}
		case syscall.AF_INET6:
			ifa := &net.IPNet{IP: make(net.IP, net.IPv6len), Mask: net.CIDRMask(int(ifam.Prefixlen), 8*net.IPv6len)}
			copy(ifa.IP, a.Value[:])
			return ifa
		}
	}
	return nil
}

func NewIfAddr(addr IPAddr, intf LinuxNetInterface) IfAddr {
	return IfAddr{addr, intf.Interface}
}

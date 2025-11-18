// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package sockaddr_test

import (
	"fmt"
	"slices"
	"sort"
	"testing"

	"github.com/hashicorp/go-sockaddr"
)

func TestIPAddrs(t *testing.T) {
	goodTestInputs := map[string][]string{
		"sockAddrs": []string{
			"10.0.0.0/8",
			"172.16.1.3/12",
			"192.168.0.0/16",
			"128.95.120.1/32",
			"192.168.1.10/24",
			"240.0.0.1/4",
		},
		"sortedBySpecificMasklen": []string{
			"128.95.120.1/32",
			"192.168.1.10/24",
			"192.168.0.0/16",
			"172.16.1.3/12",
			"10.0.0.0/8",
			"240.0.0.1/4",
		},
		"sortedByBroadMasklen": []string{
			"240.0.0.1/4",
			"10.0.0.0/8",
			"172.16.1.3/12",
			"192.168.0.0/16",
			"192.168.1.10/24",
			"128.95.120.1/32",
		},
		"sortedByNetwork": []string{
			"10.0.0.0/8",
			"128.95.120.1/32",
			"172.16.1.3/12",
			"192.168.0.0/16",
			"192.168.1.10/24",
			"240.0.0.1/4",
		},
	}
	for name, addrs := range goodTestInputs {
		t.Run(name, func(t *testing.T) {
			for _, addr := range addrs {
				_, err := sockaddr.NewSockAddr(addr)
				if err != nil {
					t.Errorf("Expected valid network: %s", addr)
				}
			}
		})
	}
}

func TestSockAddr_IPAddrs_BySpecificMaskLen(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"128.95.120.1/32",
				"192.168.1.10/24",
				"192.168.0.0/16",
				"172.16.1.3/12",
				"10.0.0.0/8",
				"240.0.0.1/4",
			},
		},
	}

	for idx, test := range testInputs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			inputAddrs := convertToSockAddrs(t, test.inputAddrs)
			sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)
			sockaddrs := slices.Clone(inputAddrs)
			filteredAddrs, _ := sockaddrs.FilterByType(sockaddr.TypeIPv4)
			for _, x := range filteredAddrs {
				switch x.(type) {
				case sockaddr.IPv4Addr:
					// pass
				default:
					t.Fatalf("invalid type")
				}
			}

			ipAddrs := make([]sockaddr.IPAddr, 0, len(filteredAddrs))
			for _, x := range filteredAddrs {
				ipAddr, ok := x.(sockaddr.IPAddr)
				if !ok {
					t.Fatalf("Unable to typecast to IPAddr")
				}
				ipAddrs = append(ipAddrs, ipAddr)
			}
			sort.Sort(sockaddr.SortIPAddrsBySpecificMaskLen{ipAddrs})

			lastLen := 32
			for i, netaddr := range ipAddrs {
				maskLen := netaddr.Maskbits()
				if lastLen < maskLen {
					t.Fatalf("Sort by specific mask length failed")
				}
				lastLen = maskLen

				if sortedAddrs[i] != netaddr {
					t.Errorf("Expected %s, received %s in iteration %d", sortedAddrs[i], netaddr, i)
				}
			}
		})
	}
}

func TestSockAddr_IPAddrs_ByBroadMaskLen(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"240.0.0.1/4",
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"128.95.120.1/32",
			},
		},
	}

	for idx, test := range testInputs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			inputAddrs := convertToSockAddrs(t, test.inputAddrs)
			sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)
			sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
			filteredAddrs, _ := sockaddrs.FilterByType(sockaddr.TypeIP)
			ipAddrs := make([]sockaddr.IPAddr, 0, len(filteredAddrs))
			for _, x := range filteredAddrs {
				ipAddr, ok := x.(sockaddr.IPAddr)
				if !ok {
					t.Fatalf("Unable to typecast to IPAddr")
				}
				ipAddrs = append(ipAddrs, ipAddr)
			}
			sort.Sort(sockaddr.SortIPAddrsByBroadMaskLen{ipAddrs})

			var lastLen int
			for i, netaddr := range ipAddrs {
				maskLen := netaddr.Maskbits()
				if lastLen > maskLen {
					t.Fatalf("Sort by specific mask length failed")
				}
				lastLen = maskLen

				if sortedAddrs[i] != netaddr {
					t.Errorf("Expected %s, received %s in iteration %d", sortedAddrs[i], netaddr, i)
				}
			}
		})
	}
}

func TestSockAddr_IPAddrs_IPAddrsByNetwork(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"10.0.0.0/8",
				"128.95.120.1/32",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
			},
		},
	}

	for idx, test := range testInputs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			inputAddrs := convertToSockAddrs(t, test.inputAddrs)
			sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)
			sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
			ipaddrs, _ := sockaddrs.FilterByType(sockaddr.TypeIP)
			sockaddr.OrderedAddrBy(sockaddr.AscAddress).Sort(ipaddrs)

			var lastIpUint sockaddr.IPv4Address
			for i, sa := range ipaddrs {
				ipv4 := *sockaddr.ToIPv4Addr(sa)
				if lastIpUint > ipv4.Address {
					t.Fatalf("Sort by network failed")
				}
				lastIpUint = ipv4.Address

				if !ipv4.Equal(sortedAddrs[i]) {
					t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, sortedAddrs[i], ipv4)
				}
			}
		})
	}
}

func TestSockAddr_IPAddrs_IPAddrsByNetworkSize(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"128.95.120.2:53",
				"128.95.120.2/32",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"128.95.120.2:8600",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"128.95.120.1/32",
				"128.95.120.2:53",
				"128.95.120.2:8600",
				"128.95.120.2/32",
				"192.168.1.10/24",
				"192.168.0.0/16",
				"172.16.1.3/12",
				"10.0.0.0/8",
				"240.0.0.1/4",
			},
		},
	}

	for idx, test := range testInputs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			inputAddrs := convertToSockAddrs(t, test.inputAddrs)
			sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)

			sockaddrs := append(sockaddr.SockAddrs(nil), inputAddrs...)
			filteredAddrs, _ := sockaddrs.FilterByType(sockaddr.TypeIP)
			ipAddrs := make([]sockaddr.IPAddr, 0, len(filteredAddrs))
			for _, x := range filteredAddrs {
				ipAddr, ok := x.(sockaddr.IPAddr)
				if !ok {
					t.Fatalf("Unable to typecast to IPAddr")
				}
				ipAddrs = append(ipAddrs, ipAddr)
			}
			sort.Sort(sockaddr.SortIPAddrsByNetworkSize{ipAddrs})

			// var prevAddr sockaddr.IPAddr
			for i, ipAddr := range ipAddrs {
				// if i == 0 {
				// 	prevAddr = ipAddr
				// 	continue
				// }

				// if prevAddr.Cmp(ipAddr) > 0 {
				// 	t.Logf("[%d] Prev:\t%v", i, prevAddr)
				// 	t.Logf("[%d] ipAddr:\t%v", i, ipAddr)
				// 	t.Fatalf("Sort by network failed")
				// }
				// prevAddr = ipAddr

				if !ipAddr.Equal(sortedAddrs[i]) {
					t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, sortedAddrs[i], ipAddr)
				}
			}
		})
	}
}

// func TestSockAddr_IPAddrs_IPAddrsByCmp(t *testing.T) {
// 	testInputs := testIPAddrsInputs{
// 		{
// 			sockAddrs: []string{
// 				"10.0.0.0/8",
// 				"172.16.1.3/12",
// 				"128.95.120.2:53",
// 				"128.95.120.2/32",
// 				"192.168.0.0/16",
// 				"128.95.120.1/32",
// 				"192.168.1.10/24",
// 				"128.95.120.2:8600",
// 				"240.0.0.1/4",
// 			},
// 			sortedSockAddrs: []string{
// 				"128.95.120.1/32",
// 				"128.95.120.2:53",
// 				"128.95.120.2:8600",
// 				"128.95.120.2/32",
// 				"192.168.1.10/24",
// 				"192.168.0.0/16",
// 				"172.16.1.3/12",
// 				"10.0.0.0/8",
// 				"240.0.0.1/4",
// 			},
// 		},
// 	}

// 	for _, test := range makeTestsFromInput(t, testInputs) {
// 		sockaddrs := append(sockaddr.SockAddrs(nil), test.sockAddrs...)
// 		ipAddrs := sockaddrs.FilterByTypeIPAddr()
// 		sort.Sort(sockaddr.SortIPAddrsByCmp{ipAddrs})
// 		t.Logf("Here: %+v", ipAddrs)

// 		var prevAddr sockaddr.IPAddr
// 		for i, ipAddr := range ipAddrs {
// 			if i == 0 {
// 				prevAddr = ipAddr
// 				continue
// 			}

// 			if prevAddr.Cmp(ipAddr) > 0 {
// 				t.Logf("[%d] Prev:\t%v", i, prevAddr)
// 				t.Logf("[%d] ipAddr:\t%v", i, ipAddr)
// 				t.Fatalf("Sort by network failed")
// 			}
// 			prevAddr = ipAddr

// 			if !ipAddr.Equal(test.sortedSockAddrs[i]) {
// 				t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, test.sortedSockAddrs[i], ipAddr)
// 			}
// 		}
// 	}
// }

func TestSockAddr_IPAddrs_IPAddrsByCmp(t *testing.T) {
	testInputs := sockAddrStringInputs{
		{
			inputAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"128.95.120.2:53",
				"128.95.120.2:53",
				"128.95.120.2/32",
				"192.168.0.0/16",
				"128.95.120.1/32",
				"192.168.1.10/24",
				"128.95.120.2:8600",
				"0:0:0:0:0:0:0:0",
				"0:0:0:0:0:0:0:1",
				"2607:f0d0:1002:0051:0000:0000:0000:0004",
				"2607:f0d0:1002:0051:0000:0000:0000:0003",
				"2607:f0d0:1002:0051:0000:0000:0000:0005",
				"[2607:f0d0:1002:0051:0000:0000:0000:0004]:8600",
				"240.0.0.1/4",
			},
			sortedAddrs: []string{
				"10.0.0.0/8",
				"172.16.1.3/12",
				"192.168.0.0/16",
				"192.168.1.10/24",
				"240.0.0.1/4",
				"128.95.120.1/32",
				"128.95.120.2/32",
				"128.95.120.2:53",
				"128.95.120.2:53",
				"128.95.120.2:8600",
				"0:0:0:0:0:0:0:0",
				"0:0:0:0:0:0:0:1",
				"2607:f0d0:1002:0051:0000:0000:0000:0003",
				"2607:f0d0:1002:0051:0000:0000:0000:0004",
				"[2607:f0d0:1002:0051:0000:0000:0000:0004]:8600",
				"2607:f0d0:1002:0051:0000:0000:0000:0005",
			},
		},
	}

	for idx, test := range testInputs {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			shuffleStrings(test.inputAddrs)

			inputAddrs := convertToSockAddrs(t, test.inputAddrs)
			sortedAddrs := convertToSockAddrs(t, test.sortedAddrs)

			sockaddr.OrderedAddrBy(sockaddr.AscType, sockaddr.AscPrivate, sockaddr.AscAddress, sockaddr.AscPort).Sort(inputAddrs)

			for i, sockAddr := range inputAddrs {
				if !sockAddr.Equal(sortedAddrs[i]) {
					t.Errorf("[%d] Sort equality failed: expected %s, received %s", i, sortedAddrs[i], sockAddr)
				}
			}
		})
	}
}

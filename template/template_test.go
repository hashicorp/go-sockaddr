package template_test

import (
	"net"
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
	socktmpl "github.com/hashicorp/go-sockaddr/template"
)

func TestSockAddr_Parse(t *testing.T) {

	interfaceList := []struct {
		SockAddress string
		Interface   net.Interface
	}{
		{
			SockAddress: "127.0.0.1/8",
			Interface: net.Interface{
				Index:        1,
				MTU:          16384,
				Name:         "lo0",
				HardwareAddr: []byte{},
				Flags:        net.FlagUp | net.FlagLoopback | net.FlagMulticast,
			},
		},
		{
			SockAddress: "::1/128",
			Interface: net.Interface{
				Index:        1,
				MTU:          16384,
				Name:         "lo0",
				HardwareAddr: []byte{},
				Flags:        net.FlagUp | net.FlagLoopback | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::1/64",
			Interface: net.Interface{
				Index:        1,
				MTU:          16384,
				Name:         "lo0",
				HardwareAddr: []byte{},
				Flags:        net.FlagUp | net.FlagLoopback | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::603e:5fff:fe48:75ff/64",
			Interface: net.Interface{
				Index:        14,
				MTU:          1500,
				Name:         "ap1",
				HardwareAddr: []byte{0x62, 0x3e, 0x5f, 0x48, 0x75, 0xff},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::2b:112f:ce21:7b6f/64",
			Interface: net.Interface{
				Index:        15,
				MTU:          1500,
				Name:         "en0",
				HardwareAddr: []byte{0x60, 0x3e, 0x5f, 0x48, 0x75, 0xff},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "2406:7400:63:ef5:1415:8bc3:fa5e:2578/64",
			Interface: net.Interface{
				Index:        15,
				MTU:          1500,
				Name:         "en0",
				HardwareAddr: []byte{0x60, 0x3e, 0x5f, 0x48, 0x75, 0xff},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "2406:7400:63:ef5:3c37:71a6:e3b4:b565/64",
			Interface: net.Interface{
				Index:        15,
				MTU:          1500,
				Name:         "en0",
				HardwareAddr: []byte{0x60, 0x3e, 0x5f, 0x48, 0x75, 0xff},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "192.168.0.102/24",
			Interface: net.Interface{
				Index:        15,
				MTU:          1500,
				Name:         "en0",
				HardwareAddr: []byte{0x60, 0x3e, 0x5f, 0x48, 0x75, 0xff},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::3871:85ff:fed8:aadc/64",
			Interface: net.Interface{
				Index:        16,
				MTU:          1500,
				Name:         "awdl0",
				HardwareAddr: []byte{0x3a, 0x71, 0x85, 0xd8, 0xaa, 0xdc},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::3871:85ff:fed8:aadc/64",
			Interface: net.Interface{
				Index:        17,
				MTU:          1500,
				Name:         "llw0",
				HardwareAddr: []byte{0x3a, 0x71, 0x85, 0xd8, 0xaa, 0xdc},
				Flags:        net.FlagUp | net.FlagBroadcast | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::83f9:d7fb:f204:cee5/64",
			Interface: net.Interface{
				Index:        18,
				MTU:          1500,
				Name:         "utun0",
				HardwareAddr: nil,
				Flags:        net.FlagUp | net.FlagPointToPoint | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::2023:8d30:551e:18d/64",
			Interface: net.Interface{
				Index:        19,
				MTU:          1380,
				Name:         "utun1",
				HardwareAddr: nil,
				Flags:        net.FlagUp | net.FlagPointToPoint | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::f9c1:463:96e6:fa8a/64",
			Interface: net.Interface{
				Index:        20,
				MTU:          2000,
				Name:         "utun2",
				HardwareAddr: nil,
				Flags:        net.FlagUp | net.FlagPointToPoint | net.FlagMulticast,
			},
		},
		{
			SockAddress: "fe80::ce81:b1c:bd2c:69e/64",
			Interface: net.Interface{
				Index:        21,
				MTU:          1000,
				Name:         "utun3",
				HardwareAddr: nil,
				Flags:        net.FlagUp | net.FlagPointToPoint | net.FlagMulticast,
			},
		},
	}
	inputList := []sockaddr.IfAddr{}
	for i := range interfaceList {
		sockAddr, _ := sockaddr.NewIPAddr(interfaceList[i].SockAddress)
		inputList = append(inputList, sockaddr.IfAddr{
			SockAddr:  sockAddr,
			Interface: interfaceList[i].Interface,
		})
	}

	tests := []struct {
		name          string
		input         string
		output        string
		fail          bool
		requireOnline bool
	}{
		{
			name:   `basic include "name"`,
			input:  `{{. | include "name" "lo0" | printf "%v"}}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast} ::1 {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "invalid input",
			input:  `{{`,
			output: ``,
			fail:   true,
		},
		{
			name:   `include "name" regexp`,
			input:  `{{. | include "name" "^(en|lo)0$" | exclude "name" "^en0$" | sort "type" | sort "address" | join "address" " " }}`,
			output: `127.0.0.1 ::1 fe80::1`,
		},
		{
			name:   `exclude "name"`,
			input:  `{{. | include "name" "^(en|lo)0$" | exclude "name" "^en0$" | sort "type" | sort "address" | join "address" " " }}`,
			output: `127.0.0.1 ::1 fe80::1`,
		},
		{
			name:   `"dot" pipeline, IPv4 type`,
			input:  `{{. | include "type" "IPv4" | include "name" "^lo0$" | sort "type" | sort "address" }}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   `include "type" "IPv6`,
			input:  `{{. | include "type" "IPv6" | include "name" "^lo0$" | sort "address" }}`,
			output: `[::1 {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "better example for IP types",
			input:  `{{. | include "type" "IPv4|IPv6" | include "name" "^lo0$" | sort "type" | sort "address" }}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast} ::1 {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "ifAddrs1",
			input:  `{{. | include "type" "IPv4" | include "name" "^lo0$"}}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   "ifAddrs2",
			input:  `{{. | include "type" "IP" | include "name" "^lo0$" | sort "type" | sort "address" }}`,
			output: `[127.0.0.1/8 {1 16384 lo0  up|loopback|multicast} ::1 {1 16384 lo0  up|loopback|multicast} fe80::1/64 {1 16384 lo0  up|loopback|multicast}]`,
		},
		{
			name:   `range "dot" example`,
			input:  `{{range . | include "type" "IP" | include "name" "^lo0$"}}{{.Name}} {{.SockAddr}} {{end}}`,
			output: `lo0 127.0.0.1/8 lo0 ::1 lo0 fe80::1/64 `,
		},
		{
			name:   `exclude "type"`,
			input:  `{{. | exclude "type" "IPv4" | include "name" "^lo0$" | sort "address" | unique "name" | join "name" " "}} {{range . | exclude "type" "IPv4" | include "name" "^lo0$"}}{{.SockAddr}} {{end}}`,
			output: `lo0 ::1 fe80::1/64 `,
		},
		{
			name:   "with variable pipeline",
			input:  `{{with $ifSet := include "type" "IPv4" . | include "name" "^lo0$"}}{{range $ifSet }}{{.Name}} {{end}}{{range $ifSet}}{{.SockAddr}} {{end}}{{end}}`,
			output: `lo0 127.0.0.1/8 `,
		},
		{
			name:   "range sample on lo0",
			input:  `{{with $ifAddrs := . | exclude "rfc" "1918" | include "name" "lo0" | sort "type,address" }}{{range $ifAddrs }}{{.Name}}/{{.SockAddr.NetIP}} {{end}}{{end}}`,
			output: `lo0/127.0.0.1 lo0/::1 lo0/fe80::1 `,
		},
		{
			name:   "non-RFC1918 on on lo0",
			input:  `{{. | exclude "rfc" "1918" | include "name" "lo0" | sort "type,address" | len | eq 3}}`,
			output: `true`,
		},
		{
			// NOTE(sean@): Difficult to reliably test includeByRFC.
			// In this case, we ass-u-me that the host running the
			// test has at least one RFC1918 address on their host.
			name:          `include "rfc"`,
			input:         `{{(. | include "rfc" "1918" | attr "name")}}`,
			output:        `en0`,
			requireOnline: true,
		},
		{
			name:   "test for non-empty array",
			input:  `{{. | include "type" "IPv4" | include "rfc" "1918" | print | len | eq (len "[]")}}`,
			output: `false`,
		},
		{
			// NOTE(sean@): This will fail if there is a public IPv4 address on loopback.
			name:   "non-IPv4 RFC1918",
			input:  `{{. | include "name" "lo0" | exclude "type" "IPv4" | include "rfc" "1918" | len | eq 0}}`,
			output: `true`,
		},
		{
			// NOTE(sean@): There are no RFC6598 addresses on most testing hosts so this should be empty.
			name:   "rfc6598",
			input:  `{{. | include "type" "IPv4" | include "rfc" "6598" | print | len | eq (len "[]")}}`,
			output: `true`,
		},
		{
			name:   "invalid RFC",
			input:  `{{. | include "type" "IPv4" | include "rfc" "99999999999" | print | len | eq (len "[]")}}`,
			output: `true`,
			fail:   true,
		},
		{
			name:   `sort asc address`,
			input:  `{{ . | include "name" "lo0" | sort "type,address" | join "address" " " }}`,
			output: `127.0.0.1 ::1 fe80::1`,
		},
		{
			name:   `sort asc address old`,
			input:  `{{with $ifSet := include "name" "lo0" . }}{{ range include "type" "IPv4" $ifSet | sort "address"}}{{ .SockAddr }} {{end}}{{ range include "type" "IPv6" $ifSet | sort "address"}}{{ .SockAddr }} {{end}}{{end}}`,
			output: `127.0.0.1/8 ::1 fe80::1/64 `,
		},
		{
			name:   `sort desc address`,
			input:  `{{ . | include "name" "lo0" | sort "type,-address" | join "address" " " }}`,
			output: `127.0.0.1 fe80::1 ::1`,
		},
		{
			name:   `sort desc address`,
			input:  `{{ . | include "name" "lo0" | include "type" "IPv6" | sort "type,-address" | join "address" " " }}`,
			output: `fe80::1 ::1`,
		},
		{
			name:   `sort asc address`,
			input:  `{{with $ifSet := include "name" "lo0" . }}{{ range include "type" "IPv6" $ifSet | sort "address"}}{{ .SockAddr }} {{end}}{{end}}`,
			output: `::1 fe80::1/64 `,
		},
		{
			name:   "lo0 limit 1",
			input:  `{{. | include "name" "lo0" | include "type" "IPv6" | sort "address" | limit 1 | len}}`,
			output: `1`,
		},
		{
			name:   "join address",
			input:  `{{. | include "name" "lo0" | include "type" "IPv6" | sort "address" | join "address" " " }}`,
			output: `::1 fe80::1`,
		},
		{
			name:   "join name",
			input:  `{{. | include "name" "lo0" | include "type" "IPv6" | sort "address" | join "name" " " }}`,
			output: `lo0 lo0`,
		},
		{
			name:   "lo0 flags up and limit 1",
			input:  `{{. | include "name" "lo0" | include "flag" "up" | sort "-type,+address" | attr "address" }}`,
			output: `::1`,
		},
		{
			name:   "math address +",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "address" "+2" | sort "+type,+address" | join "address" " " }}`,
			output: `127.0.0.3 ::3 fe80::3`,
		},
		{
			name: "math address + overflow",
			input: `|{{- with $ifAddrs := . | include "name" "^lo0$" | include "type" "IP" | math "address" "+16777217" | sort "+type,+address" -}}
		  {{- range $ifAddrs -}}
		    {{- attr "address" . }} -- {{ attr "network" . }}/{{ attr "size" . }}|{{ end -}}
		{{- end -}}`,
			output: `|128.0.0.2 -- 128.0.0.0/16777216|::100:2 -- ::100:2/1|fe80::100:2 -- fe80::/18446744073709551616|`,
		},
		{
			name:   "math address + overflow+wrap",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "address" "+4294967294" | sort "+type,+address" | join "address" " " }}`,
			output: `126.255.255.255 ::ffff:ffff fe80::ffff:ffff`,
		},
		{
			name:   "math address -",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "address" "-256" | sort "+type,+address" | join "address" " " }}`,
			output: `126.255.255.1 fe7f:ffff:ffff:ffff:ffff:ffff:ffff:ff01 ffff:ffff:ffff:ffff:ffff:ffff:ffff:ff01`,
		},
		{
			name:   "math address - underflow",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "address" "-4278190082" | sort "+type,+address" | join "address" " " }}`,
			output: `127.255.255.255 fe7f:ffff:ffff:ffff:ffff:ffff:ff:ffff ffff:ffff:ffff:ffff:ffff:ffff:ff:ffff`,
		},
		{
			// Note to readers: lo0's link-local address (::1) address has a mask of
			// /128 which means its value never changes and this is expected.  lo0's
			// site-local address has a /64 address and is expected to change.
			name:   "math network",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "network" "+2" | sort "+type,+address" | join "address" " " }}`,
			output: `127.0.0.2 ::1 fe80::2`,
		},
		{
			// Assume an IPv4 input of 127.0.0.1.  With a value of 0xff00ff01, we wrap once on purpose.
			name:   "math network + wrap",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "network" "+4278255368" | sort "+type,+address" | join "address" " " }}`,
			output: `127.0.255.8 ::1 fe80::ff00:ff08`,
		},
		{
			name:   "math network -",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | math "network" "-2" | sort "+type,+address" | join "address" " " }}`,
			output: `127.255.255.254 ::1 fe80::ffff:ffff:ffff:fffe`,
		},
		{
			// Assume an IPv4 input of 127.0.0.1.  With a value of 0xff000008 it
			// should wrap and underflow by 8.  Assume an IPv6 input of ::1.  With a
			// value of -0xff000008 the value underflows and wraps.
			name:   "math network - underflow+wrap",
			input:  `{{. | include "name" "^lo0$" | include "type" "IP" | sort "+type,+address" | math "network" "-4278190088" | join "address" " " }}`,
			output: `127.255.255.248 ::1 fe80::ffff:ffff:ff:fff8`,
		},
	}

	for i, test := range tests {
		test := test // capture range variable
		if test.name == "" {
			t.Fatalf("test number %d has an empty test name", i)
		}
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			out, err := socktmpl.ParseIfAddrs(test.input, inputList)
			if err != nil && !test.fail {
				t.Fatalf("%q: bad: %v", test.name, err)
			}

			if out != test.output && !test.fail {
				t.Fatalf("%q: Expected %+q, received %+q\n%+q", test.name, test.output, out, test.input)
			}
		})
	}
}

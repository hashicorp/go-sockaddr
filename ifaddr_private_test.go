package sockaddr

import (
	"net"
	"testing"
)

//Test_getInterfaceIPByFlags is here to facilitate testing of the private function getInterfaceIPByFlags
func Test_getInterfaceIPByFlags(t *testing.T) {

	ifAddrs := IfAddrs{
		{
			SockAddr: MustIPv4Addr("127.0.0.0/8"),
			Interface: net.Interface{
				Index: 1,
				MTU:   65536,
				Name:  "lo",
				Flags: net.FlagUp | net.FlagLoopback,
			},
		},
		{
			SockAddr: MustIPv4Addr("172.16.0.0/12"),
			Interface: net.Interface{
				Index: 2,
				MTU:   1500,
				Name:  "eth0",
				Flags: net.FlagUp,
			},
		},
		{
			SockAddr: MustIPv4Addr("169.254.0.0/16"),
			Interface: net.Interface{
				Index: 3,
				MTU:   1500,
				Name:  "dummy",
				Flags: net.FlagBroadcast,
			},
		},
		{
			SockAddr: MustIPv6Addr("fe80::/10"),
			Interface: net.Interface{
				Index: 3,
				MTU:   1500,
				Name:  "dummyv6",
				Flags: net.FlagBroadcast,
			},
		},
	}

	type args struct {
		namedIfRE string
		flags     []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "loopback: no flags provided => 127.0.0.0",
			args:    args{namedIfRE: "lo", flags: []string{}},
			want:    "127.0.0.0",
			wantErr: false,
		},
		{
			name:    "loopback: `forwardable` flag provided => empty string",
			args:    args{namedIfRE: "lo", flags: []string{"forwardable"}},
			want:    "",
			wantErr: false,
		},
		{
			name:    "private (RFC1918): no flags provided => 172.16.0.0",
			args:    args{namedIfRE: "eth0", flags: []string{}},
			want:    "172.16.0.0",
			wantErr: false,
		},
		{
			name:    "private (RFC1918): `broadcast` flag provided but iface is not broadcast => empty string",
			args:    args{namedIfRE: "eth0", flags: []string{"broadcast"}},
			want:    "",
			wantErr: false,
		},
		{
			name:    "dummy (RFC3927): no flags provided => 169.254.0.0",
			args:    args{namedIfRE: "dummy", flags: []string{}},
			want:    "169.254.0.0",
			wantErr: false,
		},
		{
			name:    "dummy (RFC3927): `forwardable` flag provided => empty string",
			args:    args{namedIfRE: "dummy", flags: []string{"forwardable"}},
			want:    "",
			wantErr: false,
		},
		{
			name:    "dummyv6 (RFC4291) IPv6: no flags provided => fe80::",
			args:    args{namedIfRE: "dummyv6", flags: []string{}},
			want:    "fe80::",
			wantErr: false,
		},
		{
			name:    "dummyv6 (RFC4291) IPv6: `forwardable` flag provided => empty string",
			args:    args{namedIfRE: "dummyv6", flags: []string{"forwardable"}},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInterfaceIP(tt.args.namedIfRE, tt.args.flags, ifAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInterfaceIPByFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getInterfaceIPByFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

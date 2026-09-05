// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package sockaddr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultInterfaceName_IPv6Only(t *testing.T) {
	installFakeIP(t, "", "default via fe80::1 dev ens192 metric 1 pref medium\n")

	ri, err := NewRouteInfo()
	if err != nil {
		t.Fatalf("NewRouteInfo: %v", err)
	}

	got, err := ri.GetDefaultInterfaceName()
	if err != nil {
		t.Fatalf("GetDefaultInterfaceName: %v", err)
	}
	if got != "ens192" {
		t.Errorf("got %+q; want %+q", got, "ens192")
	}
}

func TestGetDefaultInterfaceName_IPv4Preferred(t *testing.T) {
	installFakeIP(t,
		"default via 10.1.2.1 dev eth0 \n10.1.2.0/24 dev eth0  proto kernel  scope link  src 10.1.2.5 \n",
		"default via fe80::1 dev ens192 metric 1 pref medium\n",
	)

	ri, err := NewRouteInfo()
	if err != nil {
		t.Fatalf("NewRouteInfo: %v", err)
	}

	got, err := ri.GetDefaultInterfaceName()
	if err != nil {
		t.Fatalf("GetDefaultInterfaceName: %v", err)
	}
	if got != "eth0" {
		t.Errorf("got %+q; want %+q", got, "eth0")
	}
}

func installFakeIP(t *testing.T, ipv4Out, ipv6Out string) {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ipv4.out"), []byte(ipv4Out), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ipv6.out"), []byte(ipv6Out), 0o644); err != nil {
		t.Fatal(err)
	}

	script := `#!/bin/sh
dir=$(dirname "$0")
if [ "$1" = "route" ]; then
	cat "$dir/ipv4.out"
	exit 0
fi
if [ "$1" = "-6" ] && [ "$2" = "route" ]; then
	cat "$dir/ipv6.out"
	exit 0
fi
exit 1
`
	ipPath := filepath.Join(dir, "ip")
	if err := os.WriteFile(ipPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(ipPath, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

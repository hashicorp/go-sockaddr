// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !android

package sockaddr

import (
	"errors"
	"os/exec"
)

// NewRouteInfo returns a Linux-specific implementation of the RouteInfo
// interface.
func NewRouteInfo() (routeInfo, error) {
	// CoreOS Container Linux moved ip to /usr/bin/ip, so look it up on
	// $PATH and fallback to /sbin/ip on error.
	path, _ := exec.LookPath("ip")
	if path == "" {
		path = "/sbin/ip"
	}

	return routeInfo{
		cmds: map[string][]string{
			"ip":  {path, "route"},
			"ip6": {path, "-6", "route"},
		},
	}, nil
}

// GetDefaultInterfaceName returns the interface name attached to the default
// route on the default interface. IPv4 is preferred; IPv6 is used when no IPv4
// default route is present.
func (ri routeInfo) GetDefaultInterfaceName() (string, error) {
	for _, name := range []string{"ip", "ip6"} {
		cmd, ok := ri.cmds[name]
		if !ok || len(cmd) == 0 {
			continue
		}
		out, err := exec.Command(cmd[0], cmd[1:]...).Output()
		if err != nil {
			continue
		}
		ifName, err := parseDefaultIfNameFromIPCmd(string(out))
		if err != nil {
			continue
		}
		return ifName, nil
	}
	return "", errors.New("no default interface found")
}

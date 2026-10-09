// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !android

package sockaddr

import (
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
		cmds: map[string][]string{"ip": {path, "route"}},
	}, nil
}

// GetDefaultInterfaceName returns the interface name attached to the default
// route on the default interface.
func (ri routeInfo) GetDefaultInterfaceName() (string, error) {
	ip := ri.cmds["ip"]
	return defaultInterfaceFromIPRoute(func(args ...string) ([]byte, error) {
		return exec.Command(ip[0], args...).Output()
	}, ip[1:])
}

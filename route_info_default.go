// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build nacl || plan9 || js || wasip1
// +build nacl plan9 js wasip1

package sockaddr

// getDefaultIfName is the default interface function for unsupported platforms.
func getDefaultIfName() (string, error) {
	return "", ErrNoInterface
}

func NewRouteInfo() (routeInfo, error) {
	return routeInfo{}, ErrNoRoute
}

// GetDefaultInterfaceName returns the interface name attached to the default
// route on the default interface.
func (ri routeInfo) GetDefaultInterfaceName() (string, error) {
	return "", ErrNoInterface
}

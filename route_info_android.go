//go:build android

package sockaddr

import (
	"errors"
	"os/exec"
	"strings"
)

// NewRouteInfo returns a Android-specific implementation of the RouteInfo
// interface.
func NewRouteInfo() (routeInfo, error) {
	return routeInfo{
		cmds: map[string][]string{"ip": {"/system/bin/ip", "route", "get", "8.8.8.8"}},
	}, nil
}

// GetDefaultInterfaceName returns the interface name attached to the default
// route on the default interface.
func (ri routeInfo) GetDefaultInterfaceName() (string, error) {
	out, err := exec.Command(ri.cmds["ip"][0], ri.cmds["ip"][1:]...).Output()
	if err != nil {
		return "", err
	}

	var ifName string
	if ifName, err = parseDefaultIfNameFromIPCmdAndroid(string(out)); err != nil {
		return "", errors.New("No default interface found")
	}
	return ifName, nil
}

// parseDefaultIfNameFromIPCmdAndroid parses the default interface from ip(8) for
// Android.
func parseDefaultIfNameFromIPCmdAndroid(routeOut string) (string, error) {
	parsedLines := parseIfNameFromIPCmd(routeOut)
	if len(parsedLines) > 0 {
		ifName := strings.TrimSpace(parsedLines[0][4])
		return ifName, nil
	}

	return "", errors.New("No default interface found")
}

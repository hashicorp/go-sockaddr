// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package sockaddr

import (
	"errors"
	"slices"
)

var (
	ErrNoInterface = errors.New("no default interface found (unsupported platform)")
	ErrNoRoute     = errors.New("no route info found (unsupported platform)")
)

// RouteInterface specifies an interface for obtaining memoized route table and
// network information from a given OS.
type RouteInterface interface {
	// GetDefaultInterfaceName returns the name of the interface that has a
	// default route or an error and an empty string if a problem was
	// encountered.
	GetDefaultInterfaceName() (string, error)
}

type routeInfo struct {
	cmds map[string][]string
}

// defaultInterfaceFromIPRoute returns the interface on the IPv4 default route.
// ip route only shows that table. When it has no default route, the IPv6
// table is checked with the same arguments prefixed by -6. An IPv4 default
// route is returned without querying IPv6.
func defaultInterfaceFromIPRoute(run func(args ...string) ([]byte, error), args []string) (string, error) {
	out, err := run(args...)
	if err == nil {
		if name, perr := parseDefaultIfNameFromIPCmd(string(out)); perr == nil {
			return name, nil
		}
	}

	v6 := make([]string, 0, len(args)+1)
	v6 = append(v6, "-6")
	v6 = append(v6, args...)
	out6, err6 := run(v6...)
	if err6 != nil {
		if err != nil {
			return "", err
		}
		return "", err6
	}
	name, perr := parseDefaultIfNameFromIPCmd(string(out6))
	if perr != nil {
		return "", errors.New("no default interface found")
	}
	return name, nil
}

// VisitCommands visits each command used by the platform-specific RouteInfo
// implementation.
func (ri routeInfo) VisitCommands(fn func(name string, cmd []string)) {
	for k, v := range ri.cmds {
		cmds := slices.Clone(v)
		fn(k, cmds)
	}
}

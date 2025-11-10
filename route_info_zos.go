//go:build zos

package sockaddr

import (
   "errors"
   "os/exec"
   "regexp"
   "strings"
)

var defaultRouteRE *regexp.Regexp = regexp.MustCompile(`^Default +([0-9\.\:]+) +([^ ]+) +([0-9]+) +([^ ]+)`)

func NewRouteInfo() (routeInfo, error) {
	return routeInfo{}, nil
}

// zosGetDefaultInterfaceName executes the onetstat command and returns its output
func zosGetDefaultInterfaceName() (string, error) {
	out, err := exec.Command("/bin/onetstat", "-r").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// zosProcessOneStatOutput processes the output of onetstat -r and returns the default interface name
func zosProcessOnetstatOutput(output string) (string, error) {
	linesout := strings.Split(output, "\n")
	for _, line := range linesout {
		result := defaultRouteRE.FindStringSubmatch(line)
		if result != nil {
			return result[4], nil
		}
	}
	return "", errors.New("no default interface found")
}

// GetDefaultInterfaceName returns the interface name attached to the default route
func (ri routeInfo) GetDefaultInterfaceName() (string, error) {
	output, err := zosGetDefaultInterfaceName()
	if err != nil {
		return "", err
	}
	return zosProcessOnetstatOutput(output)
}

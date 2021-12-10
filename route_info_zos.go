// +build zos
package sockaddr

import (
   "errors"
   "os/exec"
   "regexp"
   "strings"
)

type routeInfo struct {
   cmds map[string][]string
}

var defaultRouteRE *regexp.Regexp = regexp.MustCompile(`^Default +([0-9\.\:]+) +([^ ]+) +([0-9]+) +([^ ]+)`)

// NewRouteInfo returns a ZOS-specific implementation of the RouteInfo
// interface.
func NewRouteInfo() (routeInfo, error) {
   return routeInfo{
       cmds: map[string][]string{"ip": {"/bin/onetstat", "-r"}},
   }, nil
}

// GetDefaultInterfaceName returns the interface name attached to the default
// route on the default interface.
func (ri routeInfo) GetDefaultInterfaceName() (string, error) {
   out, err := exec.Command(ri.cmds["ip"][0], ri.cmds["ip"][1:]...).Output()
   if err != nil {
       return "", err
   }
   linesout := strings.Split(string(out), "\n")
   for _, line := range linesout {
       result := defaultRouteRE.FindStringSubmatch(line)
       if result != nil {
           return result[4], nil
       }
   }
   return "", errors.New("No default interface found")
}

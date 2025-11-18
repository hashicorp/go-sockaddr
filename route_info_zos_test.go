// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package sockaddr

import "testing"

func Test_zosProcessOnetstatOutput(t *testing.T) {
	dummyOutput := `MVS TCP/IP NETSTAT CS V2R5       TCPIP Name: TCPIP           12:52:35
IPv4 Destinations
Destination        Gateway         Flags    Refcnt     Interface
-----------        -------         -----    ------     ---------
Default            0.0.0.0         UGS      0000000000 TCPIPLINK
0.0.0.0            0.0.0.0         US       0000000000 TCPIPLINK
0.0.0.0            0.0.0.0         UH       0000000000 TCPIPLINK
0.0.0.0            0.0.0.0         UH       0000000000 LOOPBACK
0.0.0.0            0.0.0.0         H        0000000000 EZAZCX
0.0.0.0            0.0.0.0         H        0000000000 EZASAMEMVS
IPv6 Destinations
DestIP:   Default
  Gw:     ::
  Intf:   TCPIPLINK6        Refcnt:  0000000000
  Flgs:   UGS               MTU:     2000
DestIP:   ::
  Gw:     ::
  Intf:   LOOPBACK6         Refcnt:  0000000000
  Flgs:   UH                MTU:     65535
DestIP:   ::
  Gw:     ::
  Intf:   TCPIPLINK6        Refcnt:  0000000000
  Flgs:   US                MTU:     2000
DestIP:   ::
  Gw:     ::
  Intf:   TCPIPLINK6        Refcnt:  0000000000
  Flgs:   UH                MTU:     9008
DestIP:   ::
  Gw:     ::
  Intf:   TCPIPLINK6        Refcnt:  0000000000
  Flgs:   UH                MTU:     9008`

	if _, err := zosProcessOnetstatOutput(dummyOutput); err != nil {
		t.Errorf("err != nil for zosProcessOnetstatOutput - %v", err)
	}

	if result, err := zosProcessOnetstatOutput(""); err == nil {
		t.Errorf("got %s; want \"\"", result)
	}
}

// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package sockaddr

import (
	"errors"
	"os/exec"
	"testing"
)

func Test_hasPowershell(t *testing.T) {
	origLookPath := execLookPath
	defer func() {
		execLookPath = origLookPath
	}()

	execLookPath = func(file string) (string, error) {
		if file == "powershell" {
			return `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`, nil
		}
		return "", exec.ErrNotFound
	}
	if !hasPowershell() {
		t.Errorf("hasPowershell() = false, want true when powershell exists")
	}

	execLookPath = func(file string) (string, error) {
		return "", exec.ErrNotFound
	}
	if hasPowershell() {
		t.Errorf("hasPowershell() = true, want false when powershell not found")
	}

	execLookPath = func(file string) (string, error) {
		return "", errors.New("lookup error")
	}
	if hasPowershell() {
		t.Errorf("hasPowershell() = true, want false on error")
	}
}

func Test_parseWindowsDefaultIfName_new_vs_old(t *testing.T) {
	if !hasPowershell() {
		t.Skip("this test requires powershell.")
		return
	}
	ri, err := NewRouteInfo()
	if err != nil {
		t.Fatalf("bad: %v", err)
	}
	psVer, err1 := ri.GetDefaultInterfaceName()
	legacyVer, err2 := ri.GetDefaultInterfaceNameLegacy()
	if err1 != nil {
		t.Errorf("err != nil for GetDefaultInterfaceName - %v", err1)
	}
	if err2 != nil {
		t.Errorf("err != nil for GetDefaultInterfaceNameLegacy - %v", err2)
	}
	if psVer != legacyVer {
		t.Errorf("got %s; want %s", psVer, legacyVer)
	}
}

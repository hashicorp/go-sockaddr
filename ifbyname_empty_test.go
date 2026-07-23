package sockaddr

import "testing"

func TestIfByNameEmptyPattern(t *testing.T) {
	_, _, err := IfByName("  ", IfAddrs{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

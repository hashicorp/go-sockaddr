package sockaddr

import "testing"

func Test_zosProcessOnetstatOutput(t *testing.T) {
	dummyOutput := "Default            127.0.0.1:9999   ABC      0000000000 TCPIPLINK"
	if _, err := zosProcessOnetstatOutput(dummyOutput); err != nil {
		t.Errorf("err != nil for zosProcessOnetstatOutput - %v", err)
	}

	if result, err := zosProcessOnetstatOutput(""); err == nil {
		t.Errorf("got %s; want \"\"", result)
	}
}

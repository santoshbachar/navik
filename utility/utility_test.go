package utility

import "testing"

func TestGetNumberAndProtocolFromPort(t *testing.T) {
	port := "8080"
	wantPort, wantProtocol := 8080, ""
	gotPort, gotProtocol := GetNumberAndProtocolFromPort(port)

	if wantPort != gotPort && wantProtocol != gotProtocol {
		t.Fatalf("failed")
	}
}

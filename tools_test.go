package http

import (
	"net"
	"testing"
)

func TestHttpTools(t *testing.T) {
	ip, err := GetServerOuterIpAddr()
	if err != nil {
		t.Fatal(ip, err)
	}

	address := net.ParseIP(ip)
	if address == nil {
		t.Fatal(ip)
	}
}

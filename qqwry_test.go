package qqwry

import (
	"testing"
)

func init() {
	if err := LoadFile("assets/qqwry.dat"); err != nil {
		panic(err)
	}
}

func TestQueryIP(t *testing.T) {
	queryIp := "1.1.1.1"
	city, isp, err := QueryIP(queryIp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("城市：%s，运营商：%s", city, isp)
}

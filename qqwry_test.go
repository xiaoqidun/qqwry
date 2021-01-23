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
	city, area, err := QueryIP(queryIp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("城市：%s，区域：%s", city, area)
}

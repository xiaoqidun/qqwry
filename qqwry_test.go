package qqwry

import (
	"testing"
)

func init() {
	if err := LoadFile("qqwry.dat"); err != nil {
		panic(err)
	}
}

func TestQueryIP(t *testing.T) {
	queryIp := "1.1.1.1"
	country, area, err := QueryIP(queryIp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(country, area)
}

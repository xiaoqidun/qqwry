package qqwry

import (
	"testing"
)

func init() {
	if err := LoadFile("assets/qqwry.ipdb"); err != nil {
		panic(err)
	}
}

func TestQueryIP(t *testing.T) {
	queryIp := "119.29.29.29"
	location, err := QueryIP(queryIp)
	if err != nil {
		t.Fatal(err)
	}
	emptyVal := func(val string) string {
		if val != "" {
			return val
		}
		return "未知"
	}
	t.Logf("国家：%s，省份：%s，城市：%s，区县：%s，运营商：%s",
		emptyVal(location.Country),
		emptyVal(location.Province),
		emptyVal(location.City),
		emptyVal(location.District),
		emptyVal(location.ISP),
	)
}

package main

import (
	"fmt"
	"github.com/xiaoqidun/qqwry"
	"github.com/xiaoqidun/qqwry/assets"
	"os"
)

func init() {
	qqwry.LoadData(assets.QQWryIpdb)
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	queryIp := os.Args[1]
	location, err := qqwry.QueryIP(queryIp)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
		return
	}
	emptyVal := func(val string) string {
		if val != "" {
			return val
		}
		return "未知"
	}
	fmt.Printf("国家：%s，省份：%s，城市：%s，区县：%s，运营商：%s\n",
		emptyVal(location.Country),
		emptyVal(location.Province),
		emptyVal(location.City),
		emptyVal(location.District),
		emptyVal(location.ISP),
	)
}

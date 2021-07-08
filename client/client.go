package main

import (
	"fmt"
	"github.com/xiaoqidun/qqwry"
	"github.com/xiaoqidun/qqwry/assets"
	"os"
)

func init() {
	qqwry.LoadData(assets.QQWryDat)
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	queryIp := os.Args[1]
	city, isp, err := qqwry.QueryIP(queryIp)
	if err != nil {
		fmt.Printf("错误：%v\n", err)
		return
	}
	fmt.Printf("城市：%s，运营商：%s\n", city, isp)
}

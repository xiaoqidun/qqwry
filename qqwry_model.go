package qqwry

import (
	"strings"
)

// Location IP位置信息
// 字段: Country 国家, Province 省份, City 城市, District 区县, ISP 运营商, IP IP地址
type Location struct {
	Country  string // 国家
	Province string // 省份
	City     string // 城市
	District string // 区县
	ISP      string // 运营商
	IP       string // IP地址
}

// Clone 克隆Location对象
// 返回: newLocation 克隆后的对象
func (l *Location) Clone() *Location {
	return &Location{
		Country:  l.Country,
		Province: l.Province,
		City:     l.City,
		District: l.District,
		ISP:      l.ISP,
		IP:       l.IP,
	}
}

// SplitResult 按照调整后的纯真社区版IP库地理位置格式返回结果
// 入参: addr 地址信息, isp 运营商信息, ipv4 IP地址
// 返回: location 位置信息
func SplitResult(addr string, isp string, ipv4 string) (location *Location) {
	location = &Location{ISP: isp, IP: ipv4}
	splitList := strings.Split(addr, "–")
	for i := 0; i < len(splitList); i++ {
		switch i {
		case 0:
			location.Country = splitList[i]
		case 1:
			location.Province = splitList[i]
		case 2:
			location.City = splitList[i]
		case 3:
			location.District = splitList[i]
		}
	}
	if location.Country == "局域网" {
		location.ISP = location.Country
	}
	return
}

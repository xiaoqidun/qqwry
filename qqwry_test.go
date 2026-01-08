package qqwry

import (
	"fmt"
	"testing"
)

// TestClient_QueryIP 测试实例IP查询功能
func TestClient_QueryIP(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		ipAddrList []string
	}{
		{
			name:     "DAT数据库",
			filePath: "assets/qqwry.dat",
			ipAddrList: []string{
				"119.29.29.29",
				"8.8.8.8",
			},
		},
		{
			name:     "IPDB数据库",
			filePath: "assets/qqwry.ipdb",
			ipAddrList: []string{
				"119.29.29.29",
				"8.8.8.8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.filePath)
			if err != nil {
				t.Fatal(err)
			}
			for _, ip := range tt.ipAddrList {
				location, err := client.QueryIP(ip)
				if err != nil {
					t.Error(err)
					continue
				}
				fmt.Printf("国家：%s，省份：%s，城市：%s，区县：%s，运营商：%s\n",
					location.Country,
					location.Province,
					location.City,
					location.District,
					location.ISP,
				)
			}
		})
	}
}

// TestGlobal_QueryIP 测试全局IP查询功能
func TestGlobal_QueryIP(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		ipAddrList []string
	}{
		{
			name:     "兼容性-DAT",
			filePath: "assets/qqwry.dat",
			ipAddrList: []string{
				"119.29.29.29",
				"8.8.8.8",
			},
		},
		{
			name:     "兼容性-IPDB",
			filePath: "assets/qqwry.ipdb",
			ipAddrList: []string{
				"119.29.29.29",
				"8.8.8.8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadFile(tt.filePath); err != nil {
				t.Fatal(err)
			}
			for _, ip := range tt.ipAddrList {
				location, err := QueryIP(ip)
				if err != nil {
					t.Error(err)
					continue
				}
				fmt.Printf("国家：%s，省份：%s，城市：%s，区县：%s，运营商：%s\n",
					location.Country,
					location.Province,
					location.City,
					location.District,
					location.ISP,
				)
			}
		})
	}
}

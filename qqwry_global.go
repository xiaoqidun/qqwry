package qqwry

import (
	"os"
	"sync"
)

// defaultClient 默认客户端，用于向后兼容
var (
	clientLock    sync.RWMutex
	defaultClient = &Client{dataType: dataTypeDat}
)

// LoadData 从内存加载IP数据库
// 入参: database DAT数据库或IPDB数据库
func LoadData(database []byte) {
	c, err := NewClientFromData(database)
	if err != nil {
		panic(err)
	}
	clientLock.Lock()
	defaultClient = c
	clientLock.Unlock()
}

// LoadFile 从文件加载IP数据库
// 入参: filepath 文件路径
// 返回: err 错误信息
func LoadFile(filepath string) (err error) {
	body, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	LoadData(body)
	return
}

// QueryIP 从内存或缓存查询IP
// 入参: ip IP地址
// 返回: location 位置信息, err 错误信息
func QueryIP(ip string) (location *Location, err error) {
	clientLock.RLock()
	c := defaultClient
	clientLock.RUnlock()
	return c.QueryIP(ip)
}

// QueryIPDat 从DAT数据库查询IP，仅加载DAT格式数据库时使用
// 入参: ipv4 IPv4地址
// 返回: location 位置信息, err 错误信息
func QueryIPDat(ipv4 string) (location *Location, err error) {
	clientLock.RLock()
	c := defaultClient
	clientLock.RUnlock()
	return c.queryIPDat(ipv4)
}

// QueryIPIpdb 从IPDB数据库查询IP，仅加载IPDB格式数据库时使用
// 入参: ip IP地址
// 返回: location 位置信息, err 错误信息
func QueryIPIpdb(ip string) (location *Location, err error) {
	clientLock.RLock()
	c := defaultClient
	clientLock.RUnlock()
	return c.queryIPIpdb(ip)
}

package qqwry

import (
	"errors"
	"os"

	"github.com/ipipdotnet/ipdb-go"
)

const (
	dataTypeDat  = 0
	dataTypeIpdb = 1
)

// Client IP查询客户端
// 字段: data DAT数据库, dataLen DAT数据库长度, ipdbCity IPDB数据库, dataType 数据类型, cache 结果缓存
type Client struct {
	data     []byte
	dataLen  uint32
	ipdbCity *ipdb.City
	dataType int
	cache    *Cache
}

// NewClient 创建新的IP查询客户端
// 入参: filePath 文件路径
// 返回: c 客户端实例, err 错误信息
func NewClient(filePath string) (c *Client, err error) {
	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewClientFromData(body)
}

// NewClientFromData 从数据创建新的IP查询客户端
// 入参: body DAT数据库或IPDB数据库
// 返回: c 客户端实例, err 错误信息
func NewClientFromData(body []byte) (c *Client, err error) {
	c = &Client{}
	c.cache = NewCache(10000)
	if len(body) > 11 && string(body[6:11]) == "build" {
		c.dataType = dataTypeIpdb
		c.ipdbCity, err = ipdb.NewCityFromBytes(body)
		if err != nil {
			return nil, err
		}
	} else {
		c.dataType = dataTypeDat
		c.data = body
		c.dataLen = uint32(len(c.data))
	}
	return c, nil
}

// QueryIP 查询IP
// 入参: ip IP地址
// 返回: location 位置信息, err 错误信息
func (c *Client) QueryIP(ip string) (location *Location, err error) {
	if v, ok := c.cache.Get(ip); ok {
		return v.Clone(), nil
	}
	switch c.dataType {
	case dataTypeDat:
		location, err = c.queryIPDat(ip)
	case dataTypeIpdb:
		location, err = c.queryIPIpdb(ip)
	default:
		return nil, errors.New("data type not support")
	}
	if err == nil && location != nil {
		c.cache.Add(ip, location)
	}
	return location, err
}

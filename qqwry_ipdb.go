package qqwry

// queryIPIpdb 从IPDB数据库查询IP
// 入参: ip IP地址
// 返回: location 位置信息, err 错误信息
func (c *Client) queryIPIpdb(ip string) (location *Location, err error) {
	ret, err := c.ipdbCity.Find(ip, "CN")
	if err != nil {
		return
	}
	location = SplitResult(ret[0], ret[1], ip)
	return location, nil
}

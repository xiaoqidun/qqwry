package qqwry

import (
	"encoding/binary"
	"errors"
	"net"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

// byte3ToUInt32 将3字节切片转换为uint32
// 入参: data 3字节切片
// 返回: uint32转换后的值
func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

// gb18030Decode 将GB18030编码解码为UTF-8
// 入参: src GB18030编码的字节切片
// 返回: string UTF-8编码的字符串
func gb18030Decode(src []byte) string {
	d, _, _ := transform.Bytes(simplifiedchinese.GB18030.NewDecoder(), src)
	return string(d)
}

// queryIPDat 从DAT数据库查询IP
// 入参: ipv4 IPv4地址
// 返回: location 位置信息, err 错误信息
func (c *Client) queryIPDat(ipv4 string) (location *Location, err error) {
	ip := net.ParseIP(ipv4).To4()
	if ip == nil {
		return nil, errors.New("ip is not ipv4")
	}
	ip32 := binary.BigEndian.Uint32(ip)
	posA := binary.LittleEndian.Uint32(c.data[:4])
	posZ := binary.LittleEndian.Uint32(c.data[4:8])
	var offset uint32 = 0
	for {
		mid := posA + (((posZ-posA)/indexLen)>>1)*indexLen
		buf := c.data[mid : mid+indexLen]
		_ip := binary.LittleEndian.Uint32(buf[:4])
		if posZ-posA == indexLen {
			offset = byte3ToUInt32(buf[4:])
			buf = c.data[mid+indexLen : mid+indexLen+indexLen]
			if ip32 < binary.LittleEndian.Uint32(buf[:4]) {
				break
			} else {
				offset = 0
				break
			}
		}
		if _ip > ip32 {
			posZ = mid
		} else if _ip < ip32 {
			posA = mid
		} else if _ip == ip32 {
			offset = byte3ToUInt32(buf[4:])
			break
		}
	}
	if offset <= 0 {
		return nil, errors.New("ip not found")
	}
	posM := offset + 4
	mode := c.data[posM]
	var ispPos uint32
	var addr, isp string
	switch mode {
	case redirectMode1:
		posC := byte3ToUInt32(c.data[posM+1 : posM+4])
		mode = c.data[posC]
		posCA := posC
		if mode == redirectMode2 {
			posCA = byte3ToUInt32(c.data[posC+1 : posC+4])
			posC += 4
		}
		for i := posCA; i < c.dataLen; i++ {
			if c.data[i] == 0 {
				addr = string(c.data[posCA:i])
				break
			}
		}
		if mode != redirectMode2 {
			posC += uint32(len(addr) + 1)
		}
		ispPos = posC
	case redirectMode2:
		posCA := byte3ToUInt32(c.data[posM+1 : posM+4])
		for i := posCA; i < c.dataLen; i++ {
			if c.data[i] == 0 {
				addr = string(c.data[posCA:i])
				break
			}
		}
		ispPos = offset + 8
	default:
		posCA := offset + 4
		for i := posCA; i < c.dataLen; i++ {
			if c.data[i] == 0 {
				addr = string(c.data[posCA:i])
				break
			}
		}
		ispPos = offset + uint32(5+len(addr))
	}
	if addr != "" {
		addr = strings.TrimSpace(gb18030Decode([]byte(addr)))
	}
	ispMode := c.data[ispPos]
	if ispMode == redirectMode1 || ispMode == redirectMode2 {
		ispPos = byte3ToUInt32(c.data[ispPos+1 : ispPos+4])
	}
	if ispPos > 0 {
		for i := ispPos; i < c.dataLen; i++ {
			if c.data[i] == 0 {
				isp = string(c.data[ispPos:i])
				if isp != "" {
					if strings.Contains(isp, "CZ88.NET") {
						isp = ""
					} else {
						isp = strings.TrimSpace(gb18030Decode([]byte(isp)))
					}
				}
				break
			}
		}
	}
	location = SplitResult(addr, isp, ipv4)
	return location, nil
}

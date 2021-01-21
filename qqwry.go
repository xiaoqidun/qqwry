package qqwry

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"sync"
)

var (
	data    []byte
	dataLen uint32
	ipCache = &sync.Map{}
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

type cache struct {
	Country string
	Area    string
}

func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

func gb18030Decode(src []byte) string {
	in := bytes.NewReader(src)
	out := transform.NewReader(in, simplifiedchinese.GB18030.NewDecoder())
	d, _ := ioutil.ReadAll(out)
	return string(d)
}

// QueryIP 从内存或缓存查询IP
func QueryIP(queryIp string) (country string, area string, err error) {
	if v, ok := ipCache.Load(queryIp); ok {
		country = v.(cache).Country
		area = v.(cache).Area
		return
	}
	ip := net.ParseIP(queryIp).To4()
	if ip == nil {
		err = errors.New("ip is not ipv4")
		return
	}
	ip32 := binary.BigEndian.Uint32(ip)
	posA := binary.LittleEndian.Uint32(data[:4])
	posZ := binary.LittleEndian.Uint32(data[4:8])
	var offset uint32 = 0
	for {
		mid := posA + (((posZ-posA)/indexLen)>>1)*indexLen
		buf := data[mid : mid+indexLen]
		_ip := binary.LittleEndian.Uint32(buf[:4])
		if posZ-posA == indexLen {
			offset = byte3ToUInt32(buf[4:])
			buf = data[mid+indexLen : mid+indexLen+indexLen]
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
		return
	}
	posM := offset + 4
	mode := data[posM]
	var areaPos uint32
	switch mode {
	case redirectMode1:
		posC := byte3ToUInt32(data[posM+1 : posM+4])
		mode = data[posC]
		var cA uint32 = 0
		if mode == redirectMode2 {
			cA = byte3ToUInt32(data[posC+1 : posC+4])
			posC += 4
		}
		for i := cA; i < dataLen; i++ {
			if data[i] == 0 {
				country = string(data[cA:i])
				break
			}
		}
		if mode != redirectMode2 {
			posC += uint32(len(country) + 1)
		}
		areaPos = posC
	case redirectMode2:
		cA := byte3ToUInt32(data[posM+1 : posM+4])
		for i := cA; i < dataLen; i++ {
			if data[i] == 0 {
				country = string(data[cA:i])
				break
			}
		}
		areaPos = offset + 8
	default:
		cA := offset + 4
		for i := cA; i < dataLen; i++ {
			if data[i] == 0 {
				country = string(data[cA:i])
				break
			}
		}
		areaPos = offset + uint32(5+len(country))
	}
	areaMode := data[areaPos]
	if areaMode == redirectMode1 || areaMode == redirectMode2 {
		areaPos = byte3ToUInt32(data[areaPos+1 : areaPos+4])
	}
	if areaPos > 0 {
		for i := areaPos; i < dataLen; i++ {
			if data[i] == 0 {
				area = string(data[areaPos:i])
				break
			}
		}
	}
	country = gb18030Decode([]byte(country))
	area = gb18030Decode([]byte(area))
	ipCache.Store(queryIp, cache{Country: country, Area: area})
	return
}

// LoadData 从内存加载IP数据库
func LoadData(database []byte) {
	data = database
	dataLen = uint32(len(data))
}

// LoadFile 从文件加载IP数据库
func LoadFile(filepath string) (err error) {
	data, err = ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	dataLen = uint32(len(data))
	return
}
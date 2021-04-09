package qqwry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"strings"
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
	City string
	Area string
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
func QueryIP(queryIp string) (city string, area string, err error) {
	if v, ok := ipCache.Load(queryIp); ok {
		city = v.(cache).City
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
		err = errors.New("ip not found")
		return
	}
	posM := offset + 4
	mode := data[posM]
	var areaPos uint32
	switch mode {
	case redirectMode1:
		posC := byte3ToUInt32(data[posM+1 : posM+4])
		mode = data[posC]
		posCA := posC
		if mode == redirectMode2 {
			posCA = byte3ToUInt32(data[posC+1 : posC+4])
			posC += 4
		}
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				city = string(data[posCA:i])
				break
			}
		}
		if mode != redirectMode2 {
			posC += uint32(len(city) + 1)
		}
		areaPos = posC
	case redirectMode2:
		posCA := byte3ToUInt32(data[posM+1 : posM+4])
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				city = string(data[posCA:i])
				break
			}
		}
		areaPos = offset + 8
	default:
		posCA := offset + 4
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				city = string(data[posCA:i])
				break
			}
		}
		areaPos = offset + uint32(5+len(city))
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
	city = strings.TrimSpace(gb18030Decode([]byte(city)))
	area = strings.TrimSpace(gb18030Decode([]byte(area)))
	ipCache.Store(queryIp, cache{City: city, Area: area})
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

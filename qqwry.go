package qqwry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ipipdotnet/ipdb-go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	data          []byte
	dataLen       uint32
	ipdbCity      *ipdb.City
	dataType      = dataTypeDat
	locationCache = &sync.Map{}
)

const (
	dataTypeDat  = 0
	dataTypeIpdb = 1
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

type Location struct {
	Country  string // 国家
	Province string // 省份
	City     string // 城市
	District string // 区县
	ISP      string // 运营商
	IP       string // IP地址
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
	d, _ := io.ReadAll(out)
	return string(d)
}

// QueryIP 从内存或缓存查询IP
func QueryIP(ip string) (location *Location, err error) {
	if v, ok := locationCache.Load(ip); ok {
		return v.(*Location), nil
	}
	switch dataType {
	case dataTypeDat:
		return QueryIPDat(ip)
	case dataTypeIpdb:
		return QueryIPIpdb(ip)
	default:
		return nil, errors.New("data type not support")
	}
}

// QueryIPDat 从dat查询IP，仅加载dat格式数据库时使用
func QueryIPDat(ipv4 string) (location *Location, err error) {
	ip := net.ParseIP(ipv4).To4()
	if ip == nil {
		return nil, errors.New("ip is not ipv4")
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
		return nil, errors.New("ip not found")
	}
	posM := offset + 4
	mode := data[posM]
	var ispPos uint32
	var addr, isp string
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
				addr = string(data[posCA:i])
				break
			}
		}
		if mode != redirectMode2 {
			posC += uint32(len(addr) + 1)
		}
		ispPos = posC
	case redirectMode2:
		posCA := byte3ToUInt32(data[posM+1 : posM+4])
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				addr = string(data[posCA:i])
				break
			}
		}
		ispPos = offset + 8
	default:
		posCA := offset + 4
		for i := posCA; i < dataLen; i++ {
			if data[i] == 0 {
				addr = string(data[posCA:i])
				break
			}
		}
		ispPos = offset + uint32(5+len(addr))
	}
	if addr != "" {
		addr = strings.TrimSpace(gb18030Decode([]byte(addr)))
	}
	ispMode := data[ispPos]
	if ispMode == redirectMode1 || ispMode == redirectMode2 {
		ispPos = byte3ToUInt32(data[ispPos+1 : ispPos+4])
	}
	if ispPos > 0 {
		for i := ispPos; i < dataLen; i++ {
			if data[i] == 0 {
				isp = string(data[ispPos:i])
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
	locationCache.Store(ipv4, location)
	return location, nil
}

// QueryIPIpdb 从ipdb查询IP，仅加载ipdb格式数据库时使用
func QueryIPIpdb(ip string) (location *Location, err error) {
	ret, err := ipdbCity.Find(ip, "CN")
	if err != nil {
		return
	}
	location = SplitResult(ret[0], ret[1], ip)
	locationCache.Store(ip, location)
	return location, nil
}

// LoadData 从内存加载IP数据库
func LoadData(database []byte) {
	if string(database[6:11]) == "build" {
		dataType = dataTypeIpdb
		loadCity, err := ipdb.NewCityFromBytes(database)
		if err != nil {
			panic(err)
		}
		ipdbCity = loadCity
		return
	}
	data = database
	dataLen = uint32(len(data))
}

// LoadFile 从文件加载IP数据库
func LoadFile(filepath string) (err error) {
	body, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	LoadData(body)
	return
}

// SplitResult 按照调整后的纯真社区版IP库地理位置格式返回结果
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

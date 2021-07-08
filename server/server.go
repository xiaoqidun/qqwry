package main

import (
	"encoding/json"
	"flag"
	"github.com/xiaoqidun/qqwry"
	"github.com/xiaoqidun/qqwry/assets"
	"net"
	"net/http"
)

type resp struct {
	IP   string `json:"ip"`
	Err  string `json:"err"`
	City string `json:"city"`
	Isp  string `json:"isp"`
}

func init() {
	qqwry.LoadData(assets.QQWryDat)
}

func main() {
	listen := flag.String("listen", "127.0.0.1:80", "http server listen addr")
	flag.Parse()
	http.HandleFunc("/ip/", IpAPI)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		panic(err)
	}
}

func IpAPI(writer http.ResponseWriter, request *http.Request) {
	ip := request.URL.Path[4:]
	if ip == "" {
		ip, _, _ = net.SplitHostPort(request.RemoteAddr)
	}
	rw := &resp{IP: ip}
	city, isp, err := qqwry.QueryIP(ip)
	if err != nil {
		rw.Err = err.Error()
	} else {
		rw.City = city
		rw.Isp = isp
	}
	b, _ := json.MarshalIndent(rw, "", "    ")
	_, _ = writer.Write(b)
}

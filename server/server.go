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
	Data    *qqwry.Location `json:"data"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
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
	write := &resp{}
	location, err := qqwry.QueryIP(ip)
	if err != nil {
		write.Message = err.Error()
	} else {
		write.Data = location
		write.Success = true
	}
	b, _ := json.MarshalIndent(write, "", "    ")
	_, _ = writer.Write(b)
}

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
	qqwry.LoadData(assets.QQWryIpdb)
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
	response := &resp{}
	location, err := qqwry.QueryIP(ip)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Data = location
		response.Success = true
	}
	b, _ := json.MarshalIndent(response, "", "  ")
	_, _ = writer.Write(b)
}

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

const (
	port = ":9222"
)

func main() {
	server := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(handler),
	}

	fmt.Printf("HTTP server listening on http://127.0.0.1%s ...\n", port)
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("启动 web 服务失败, err: %s", err.Error()))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	address, err := getMacAddress()
	if err != nil {
		_, err = fmt.Fprintf(w, "Found an error: %s", err.Error())
		if err != nil {
			log.Println("Failed to write response: ", err.Error())
		}
		return
	}

	// 不同的 host 返回不同的数据
	var msg string
	if strings.Contains(r.Host, "8855") {
		msg = fmt.Sprintf("http B mac地址是：%s", address)
	} else if strings.Contains(r.Host, "8854") {
		msg = fmt.Sprintf("https B mac地址是：%s", address)
	} else {
		msg = address
	}

	_, err = fmt.Fprint(w, msg)
	if err != nil {
		log.Println("Failed to write response: ", err.Error())
	}
}

// getMacAddress 获取 mac 地址
func getMacAddress() (string, error) {
	var macAddr string
	ifaces, err := net.Interfaces()
	// 发生错误
	if err != nil {
		log.Println("Failed to get mac address", err.Error())
		return "", err
	}

	for _, iface := range ifaces {
		// 获取网卡地址：首先需要判断接口是否可用，并且忽略回环接口（不需要获取 mac 地址）
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			macAddr = iface.HardwareAddr.String()
			break
		}
	}

	return macAddr, nil
}

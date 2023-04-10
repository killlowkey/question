package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	debug     = true
	httpPort  = ":8855"
	httpsPort = ":8854"
)

var backendAddress string

func init() {
	if debug {
		backendAddress = "http://127.0.0.1:9222"
	} else {
		backendAddress = "http://10.18.13.101:9222"
	}
}

// main 设备 A 代码通过反向代理服务器实现
// 启动 http 与 https 服务器，将请求转发到设备B去
func main() {
	// 设备 B 的地址和端口
	backendURL, err := url.Parse(backendAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 创建反向代理服务器
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// 处理请求的函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		// 打印客户端地址和端口
		log.Printf("Client address: %s, Client port: %s", r.RemoteAddr, r.RemoteAddr)

		// 转发请求到设备 B
		proxy.ServeHTTP(w, r)
	}

	// 创建 http 服务器
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: http.HandlerFunc(handler),
	}
	// 创建 https 服务器
	httpsServer := &http.Server{
		Addr:    httpsPort,
		Handler: http.HandlerFunc(handler),
	}

	// 启动 http 服务器
	go func() {
		fmt.Printf("HTTP server listening on http://127.0.0.1%s ...\n", httpPort)
		log.Fatal(httpServer.ListenAndServe())
	}()

	// 启动 https 服务器
	go func() {
		fmt.Printf("HTTPS server listening on https://127.0.0.1%s ...\n", httpsPort)
		log.Fatal(httpsServer.ListenAndServeTLS(
			"./configs/server.crt",
			"./configs/server.key",
		))
	}()

	// 阻塞主 routine
	select {}
}

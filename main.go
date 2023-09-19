package main

import (
	"fmt"
	"net/http"
)

func main() {
	addr := ":8080"
	fmt.Printf("Server started at %s\n", addr)
	http.HandleFunc("/", retrieveIPInfo) // 设置路由和处理器函数
	http.ListenAndServe(addr, nil)       // 启动服务器
}

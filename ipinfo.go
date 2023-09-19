package main

import (
	"embed"
	"encoding/json"
	"github.com/oschwald/geoip2-golang"
	"io"
	"net"
	"net/http"
)

// IPInfo 结构体用于存储IP相关信息
type IPInfo struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Location struct {
		AccuracyRadius uint16  `json:"accuracy_radius"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		MetroCode      uint    `json:"metro_code"`
		TimeZone       string  `json:"time_zone"`
	} `json:"location"`
}

// db 是 geoip2 数据库的全局对象
var db *geoip2.Reader

// localNetworkNames 用于存储本地网络的名字，支持中文和英文
var localNetworkNames = map[string]string{
	"zh-CN": "局域网",
	"en":    "local network",
}

//go:embed GeoLite2-City.mmdb
var mmdb embed.FS

// init 函数用于初始化 geoip2 数据库
func init() {
	fs, err := mmdb.Open("GeoLite2-City.mmdb") // 打开嵌入的数据库文件
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(fs) // 从文件中读取数据库
	if err != nil {
		panic(err)
	}
	db, err = geoip2.FromBytes(data)
	if err != nil {
		panic(err)
	}
}

// isPublicIP 函数用于判断给定的IP是否是公网IP
func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10, ip4[0] == 192 && ip4[1] == 168, ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		default:
			return true
		}
	}

	return false
}

// retrieveIPInfo 处理IP查询请求
func retrieveIPInfo(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language") // 从请求参数中获取语言选项
	if language != "zh-CN" && language != "en" {
		language = "zh-CN"
	}

	ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}
	ip := net.ParseIP(ipStr)

	if !isPublicIP(ip) { // 如果IP地址是私网地址，返回预设好的本地网络信息
		renderResponse(w, http.StatusOK, IPInfo{IP: ip.String(), City: localNetworkNames[language]})
		return
	}

	city, err := db.City(ip) // 查询IP地址信息
	if err != nil {          // 如果查询出错，返回错误信息
		http.Error(w, ipStr+" is invalid ip", http.StatusBadRequest)
		return
	}

	// 根据查询结果创建IP信息对象
	ipInfo := IPInfo{
		IP:       ip.String(),
		Country:  city.Country.Names[language],
		Province: city.Subdivisions[0].Names[language],
		City:     city.City.Names[language],
		Location: struct {
			AccuracyRadius uint16  `json:"accuracy_radius"` //精度半径
			Latitude       float64 `json:"latitude"`        //纬度
			Longitude      float64 `json:"longitude"`       //经度
			MetroCode      uint    `json:"metro_code"`      //都市区编码
			TimeZone       string  `json:"time_zone"`       //时区
		}{
			AccuracyRadius: city.Location.AccuracyRadius,
			Latitude:       city.Location.Latitude,
			Longitude:      city.Location.Longitude,
			MetroCode:      city.Location.MetroCode,
			TimeZone:       city.Location.TimeZone,
		},
	}

	renderResponse(w, http.StatusOK, ipInfo) // 返回请求结果
}

// renderResponse 函数用于构造JSON格式的响应并发送
func renderResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)               //将payload内容编组为json
	w.Header().Set("Content-Type", "application/json") //设置Header内容，声明返回体的内容类型是JSON
	w.WriteHeader(code)                                //写入HTTP响应状态码
	w.Write(response)                                  //写入HTTP响应主体，即我们的payload
}

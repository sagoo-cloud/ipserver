# ipserver 获取IP地理位置服务

在应用端可以通过调用ipserver服务获取IP的地理位置信息，ipserver服务是一个http服务，可以通过http请求的方式获取IP的地理位置信息。

在go程序中的使用方式：

```go
// 代码示例 

resp, err := http.Get("http://localhost:8080/?language=cn&format=json")
if err != nil {
  // 处理错误
}

// 解析响应
var ipInfo IPInfo 
json.Unmarshal(resp.Body, &ipInfo)

// 使用ipInfo
fmt.Println(ipInfo.City)

```



这个IP地理信息库是基于MaxMind提供的 GeoLite2 和 GeoIP2 数据库。

需要去MAXMID官网注册账号，然后才能下载最新的IP地理位置库。

注册地址：https://www.maxmind.com/en/geolite2/signup
登录
下载最新的IP地理位置库
登陆后点击 Download Databases 进入下载选择页面
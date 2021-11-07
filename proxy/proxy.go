package proxy

type Proxy struct {
	Code int64 `json:"code"`
	Data []struct {
		Ip         string `json:"ip"`
		Port       int64  `json:"port"`
		ExpireTime string `json:"expire_time"`
	}
}

var proxyUrl string = "http://webapi.http.zhimacangku.com/getip?num=1&type=2&pro=&city=0&yys=0&port=1&pack=134598&ts=1&ys=0&cs=0&lb=1&sb=0&pb=4&mr=1&regions="

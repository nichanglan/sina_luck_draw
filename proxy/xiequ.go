package proxy

type XiequProxy struct {
	Code int64 `json:"code"`
	Data []struct {
		Ip   string `json:"IP"`
		Port int64  `json:"Port"`
	}
}

var XiequProxyUrl string = "http://api.xiequ.cn/VAD/GetIp.aspx?act=get&num=200&time=30&plat=0&re=1&type=0&so=1&ow=1&spl=1&addr=&db=1"

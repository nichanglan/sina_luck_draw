package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sina/config"
	notify "sina/draw_notify"
	"sina/function"
	"sina/model"
	"time"

	"github.com/nichanglan/weibo-notify/sender/email"
	"gopkg.in/yaml.v2"
)

// TODO: SUB 自动更新？？？

func main() {
	go Notify()
	for {
		run()
		notify.C.GlobalStartTime = notify.C.GlobalEndTime
		notify.C.GlobalStartTime += notify.C.Offset
		time.Sleep(time.Duration(c.NotiftPeriod * time.Second.Nanoseconds()))
	}
}

func init() {
	loadConfig()
	addDefaultConfig()
	model.InitDB()
	notify.C = notify.Config{
		Sub:             c.Sina.Sub,
		Email:           c.Email,
		Offset:          c.NotiftPeriod,
		GlobalStartTime: time.Now().Unix() - c.NotiftPeriod,
		GlobalEndTime:   time.Now().Unix(),
	}
	notify.InitConfig()
}

func Notify() {
	for {
		notify.DrawNotify()
		fmt.Println("waiting next notify period...............")
		time.Sleep(time.Duration(c.Period * time.Second.Nanoseconds()))
	}
}

func run() {
	function.GetLuckSearchApi()
	function.FollowSet()
	function.HuaTiZhuanFa()
	function.LikeSet()
	fmt.Println("waiting next scan period...............")
}

type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"db"`
	Sina struct {
		Sub string `yaml:"sub"`
		Uid string `yaml:"uid"`
	} `yaml:"sina"`
	Period       int64             `yaml:"period"`
	NotiftPeriod int64             `yaml:"notify_period"`
	Email        email.EmailConfig `yaml:"email"`
}

var c Config

func loadConfig() {
	configFilePath := flag.String("config", "config/config.yml", "config file")
	if configFilePath != nil {
		configFile, err := ioutil.ReadFile(*configFilePath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(configFile, &c)
		if err != nil {
			panic(err)
		}
	}
}

func addDefaultConfig() {
	config.UrlConfig = make(map[string]interface{})
	config.UrlConfig["SUB"] = c.Sina.Sub
	config.UrlConfig["UID"] = c.Sina.Uid

	config.UrlConfig["LUCKING"] = "https://weibo.com/p/100808557b69009a8ef6588f9124fe9c30d36c/super_index"
	config.UrlConfig["LUCKING_TIME"] = "https://m.weibo.cn/api/container/getIndex?jumpfrom=weibocom&containerid=100808557b69009a8ef6588f9124fe9c30d36c_-_sort_time"
	config.UrlConfig["LUCKING_SEARCH"] = "https://m.weibo.cn/api/container/getIndex?containerid=100103type%3D1%26q%3D@%E5%BE%AE%E5%8D%9A%E6%8A%BD%E5%A5%96%E5%B9%B3%E5%8F%B0&page_type=searchall"              //关键词1
	config.UrlConfig["LUCKING_SEARCH_ZHUANFA"] = "https://m.weibo.cn/api/container/getIndex?containerid=100103type%3D1%26q%3D%25E8%25BD%25AC%25E5%258F%2591%25E6%258A%25BD%25E5%25A5%2596&page_type=searchall" //关键词2
	config.UrlConfig["LUCKING_SEARCH_XIANGQING"] = "https://m.weibo.cn/api/container/getIndex?containerid=100103type%3D1%26q%3D%E6%8A%BD%E5%A5%96%E8%AF%A6%E6%83%85&page_type=searchall"                       //关键词3
	config.UrlConfig["LUCKING_STATUS"] = "https://lottery.media.weibo.com/lottery/h5/history/list?mid="                                                                                                        //查看是否存在页面

	config.UrlConfig["REFERER"] = "https://m.weibo.cn/p/100808557b69009a8ef6588f9124fe9c30d36c/super_index?jumpfrom=weibocom"
	config.UrlConfig["PDETAIL"] = "100808557b69009a8ef6588f9124fe9c30d36c"
	//config.UrlConfig["PAGE_ID"] = "page_100808_super_index"
	config.UrlConfig["LOCATION"] = "100808557b69009a8ef6588f9124fe9c30d36c"

	//自己的uid
	config.UrlConfig["COMMENT_URL"] = "https://weibo.com/aj/v6/comment/add"
	config.UrlConfig["FOLLOW_URL"] = "https://weibo.com/aj/f/followed"
	config.UrlConfig["LIKE_URL"] = "https://weibo.com/aj/v6/like/add"
	config.UrlConfig["PAGE_ID"] = "page_100505_home"
	config.UrlConfig["ZHUANFA_URL"] = "https://weibo.com/aj/v6/mblog/forward"
	config.UrlConfig["PAGES_ID"] = "page_100606_home"

	config.DBConfig = make(map[string]interface{})
	// 自定义配置
	config.DBConfig["DB_HOST"] = c.DB.Host
	config.DBConfig["DB_PORT"] = c.DB.Port
	config.DBConfig["DB_NAME"] = c.DB.Name
	config.DBConfig["DB_USER"] = c.DB.User
	config.DBConfig["DB_PWD"] = c.DB.Password
	config.DBConfig["DB_CHARSET"] = "utf8mb4"
	config.DBConfig["DB_PREFIX"] = "sl_"

	// 其他
	config.DBConfig["DB_MAX_OPEN_CONNS"] = "20"          // 连接池最大连接数
	config.DBConfig["DB_MAX_IDLE_CONNS"] = "10"          // 连接池最大空闲数
	config.DBConfig["DB_MAX_LIFETIME_CONNS"] = time.Hour // 连接池链接最长生命周期

	config.DBConfig["REDIS_HOST"] = "127.0.0.1" //之前想过用redis缓存代理ip，无视就好
	config.DBConfig["REDIS_PORT"] = 6379
	config.DBConfig["REDIS_PWD"] = ""
	config.DBConfig["REDIS_SELECT"] = 3
	config.DBConfig["MAX_IDLE"] = 512
	config.DBConfig["MAX_ACTIVE"] = 10
	config.DBConfig["MAX_IDLE_TIMEOUT"] = 200
	config.DBConfig["TIMEOUT"] = 200

}

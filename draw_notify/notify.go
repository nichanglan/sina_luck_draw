package notify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nichanglan/weibo-notify/sender"
	"github.com/nichanglan/weibo-notify/sender/email"
)

type Config struct {
	Sub             string            `yaml:"sub"`
	Email           email.EmailConfig `yaml:"email"`
	Offset          int64             `yaml:"offset"`
	GlobalStartTime int64
	GlobalEndTime   int64
}

var C Config
var globalStartTime, globalEndTime int64

var weiboSender sender.Sender

func InitConfig() {

	weiboSender = email.EmailSender{
		Conf: C.Email,
	}
}

func DrawNotify() {
	url := "https://weibo.com/ajax/statuses/mentions"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.AddCookie(&http.Cookie{Name: "SUB", Value: C.Sub})
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"94\", \"Google Chrome\";v=\"94\", \";Not A Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("get blof error", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error", err)
	}
	var detail Resp
	err = json.Unmarshal(body, &detail)
	if err != nil {
		fmt.Println("unmarshal body ", string(body), " error", err)
	}
	for _, v := range detail.Data.Statuses {
		sendAt := parseToTimestamp(v.CreatedAt)
		if sendAt < globalStartTime {
			break
		}
		sendHtml("中奖了", v.Text)
	}
}

func sendHtml(subject, body string) {
	err := weiboSender.Send(subject, body)
	if err != nil {
		fmt.Println("send error", err)
	}
}

func parseToTimestamp(s string) int64 {
	layout := time.RubyDate
	t, err := time.Parse(layout, s)
	if err != nil {
		fmt.Println("parse time error", err)
	}
	return t.Unix()
}

type Resp struct {
	Ok   int `json:"ok"`
	Data struct {
		PreviousCursor int     `json:"previous_cursor"`
		TotalNumber    int     `json:"total_number"`
		TipsShow       int     `json:"tips_show"`
		NextCursor     int     `json:"next_cursor"`
		Hasvisible     bool    `json:"hasvisible"`
		MissIds        []int64 `json:"miss_ids"`
		Statuses       []struct {
			CreatedAt string `json:"created_at"`
			ID        int64  `json:"id"`
			Idstr     string `json:"idstr"`
			Mid       string `json:"mid"`
			Mblogid   string `json:"mblogid"`
			TextRaw   string `json:"text_raw"`
			Text      string `json:"text"`
		} `json:"statuses"`
	} `json:"data"`
}

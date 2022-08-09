package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Results []Results `json:"results"`
}

type Results struct {
	Location Location `json:"location"`
	Hourly   []Hourly `json:"hourly"`
}

type Location struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Country        string `json:"country"`
	Path           string `json:"path"`
	Timezone       string `json:"timezone"`
	TimezoneOffset string `json:"timezone_offset"`
}

type Hourly struct {
	Time          string `json:"time"`
	Text          string `json:"text"`
	Code          string `json:"code"`
	Temperature   string `json:"temperature"`
	Humidity      string `json:"humidity"`
	WindDirection string `json:"wind_direction"`
	WindSpeed     string `json:"wind_speed"`
}

var (
	key = ""
	u   = ""
	p   = getMD5("#")
	m   = ""
	c   = "【Lee】今天下雨"
)

func main() {
	var base int64 = 24 * 60 * 60 * 1000
	for {
		t := time.UnixMilli(time.Now().UnixMilli() / base * base)
		t = t.Add(time.Duration(23) * time.Hour)

		sleepTime := t.UnixMilli() - time.Now().UnixMilli()

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		if err := doMain(); err == nil {
			log.Println("success")
		} else {
			log.Println(err)
		}
	}
}

func doMain() error {
	resp, err := http.Get("https://api.seniverse.com/v3/weather/hourly.json?key=" + key + "&location=shenzhen&language=zh-Hans&unit=c&start=0&hours=24")
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := Response{}
	if err = json.Unmarshal(data, &response); err != nil {
		return err
	}

	for _, result := range response.Results {
		for i, hourly := range result.Hourly {
			if i >= 6 && i <= 19 {
				if strings.Contains(hourly.Text, "雨") {
					send()
					break
				}
			}
		}
	}

	return nil
}

func send() {
	response, err := http.Get(fmt.Sprintf("https://api.smsbao.com/sms?u=%s&p=%s&m=%s&c=%s", u, p, m, c))
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(response.Body)

	fmt.Println(string(data))
}

func getMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

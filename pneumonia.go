package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type Config struct {
	City     []string `json:"city"`
	Mail     string   `json:"mail"`
	Password string   `json:"password"`
	Server   string   `json:"server"`
	Addr     string   `json:"addr"`
	To       []string `json:"to"`
	Nickname string   `json:"nickname"`
}

var config Config

func init() {
	file, err := os.Open("feiyan.json")
	if err != nil {
		log.Fatalf("missing configuration file")
	}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalf("error configuration")
	}
}

func condition() string {
	care := config.City
	s := make([]string, len(care))
	for c := range care {
		s[c] = fmt.Sprintf("div[city=%s]", care[c])
	}
	return strings.Join(s, ",")
}
func getNext() time.Time {
	t := time.Now().Add(24 * time.Hour)
	return time.Date(t.Year(), t.Month(), t.Day(), 10, 0, 0, 0, t.Location())
}
func main() {
	for {
		timer := time.NewTimer(getNext().Sub(time.Now()))
		<-timer.C
		send()
	}
}
func send() {
	content := ""
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var html string
	log.Printf("start request!\n")
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://news.qq.com/zt2020/page/feiyan.htm`),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		fmt.Printf("%s: occur error: %s", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	log.Printf("finish request!\n")
	//file,_:=os.Create("index.html")
	//file.WriteString(html)
	//file.Close()
	document, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
	if err != nil {
		fmt.Printf("%s: occur error: %s", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	document.Find("div#charts div.topdataWrap div.recentNumber>div").Each(func(i int, selection *goquery.Selection) {
		var text string
		var number string
		var add string
		text = selection.Find(".text").Text()
		number = selection.Find(".number").Text()
		add = selection.Find(".add").Text()
		content += fmt.Sprintf("%s：%s\t%s\n", text, number, add)
	})
	document.Find(condition()).Each(func(i int, selection *goquery.Selection) {
		city := selection.Find("h2").Text()
		number := make([]string, 4)
		selection.Find("div").Each(func(i int, selection *goquery.Selection) {
			if i < 4 {
				number[i] = selection.Text()
			}
		})
		content += fmt.Sprintf("%s 新增确诊：%s\t累计确诊：%s\t治愈：%s\t死亡：%s\n", city, number[0], number[1], number[2], number[3])
	})
	auth := smtp.PlainAuth("", config.Mail, config.Password, config.Server)
	to := config.To
	nickname := config.Nickname
	user := config.Mail
	subject := "每日疫情"
	contentType := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + content)
	err = smtp.SendMail(config.Addr, auth, user, to, msg)
	if err != nil {
		log.Printf("send mail error: %v", err)
	}
	log.Printf("send mail successfully\n")
}

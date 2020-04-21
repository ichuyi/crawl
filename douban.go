package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main() {
	file, err := os.Create("douban.csv")
	if err != nil {
		fmt.Println("create file error")
		return
	}
	file.WriteString("\xEF\xBB\xBF")
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.Write([]string{
		"name", "链接",
	})
	defer writer.Flush()
	c := colly.NewCollector()
	c.OnRequest(func(request *colly.Request) {
		log.Println("start request")
	})
	c.OnError(func(response *colly.Response, e error) {
		log.Printf("occur error: %s", e.Error())
	})
	c.OnResponse(func(response *colly.Response) {
		log.Printf("get response")
	})
	c.OnHTML("div#content div div.article div div table tbody tr td div a", func(element *colly.HTMLElement) {
		err = writer.Write([]string{
			element.Text, element.Attr("href"),
		})
	})
	err = c.Visit("https://movie.douban.com/chart")
	if err != nil {
		log.Println(err.Error())
	}
}

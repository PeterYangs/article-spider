package main

import (
	"context"
	"fmt"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.g74.com",
		Channel:      "/azyx/[PAGE]/",
		ListSelector: "body > div:nth-child(4) > ul > li",
		HrefSelector: " div > h4 > a",
		PageStart:    1,
		Length:       20,
		ListFields: map[string]articleSpider.Field{
			"title": {Types: articleSpider.Text, Selector: "  div > h4 > a"},
		},
		DetailFields: map[string]articleSpider.Field{
			"content": {Types: articleSpider.HtmlWithImage, Selector: "body > div.w1200.clearfix.jspage > div.left.w810.mr50 > div:nth-child(2) > div.nawContent"},
		},
		DetailCoroutineNumber: 1,
		FilterError:           true,
		HttpHeader: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal, context.Background())

	err := s.Start()

	if err != nil {

		fmt.Println(err)
	}

}

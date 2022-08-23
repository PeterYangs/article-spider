package main

import (
	"context"
	"fmt"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

//采集文章标题
func main() {

	f := articleSpider.Form{
		Host:         "https://www.g74.com",
		Channel:      "/azyx/[PAGE]/",
		ListSelector: "body > div:nth-child(4) > ul > li",
		PageStart:    1,
		Length:       2,
		ListFields: map[string]articleSpider.Field{
			"title": {Types: articleSpider.Text, Selector: "  div > h4 > a"},
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

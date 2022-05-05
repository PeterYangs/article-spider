package main

import (
	"context"
	"fmt"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.925g.com",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       10,
		ListFields: map[string]articleSpider.Field{

			"title": {ExcelHeader: "K", Types: articleSpider.Text, Selector: " a > div > span"},
		},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 2,
		ResultCallback: func(item map[string]string, form *articleSpider.Form) {

			for s2, s3 := range item {

				fmt.Println(s2, ":", s3)

			}

		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal, context.Background())

	s.Start()

}

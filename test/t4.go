package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host: "https://www.925g.com",
		ChannelFunc: func(form *articleSpider.Form) []string {

			return []string{
				"https://www.925g.com/zixun_page1.html/",
				"https://www.925g.com/zixun_page2.html/",
				"https://www.925g.com/zixun_page3.html/",
				"https://www.925g.com/zixun_page4.html/",
				"https://www.925g.com/zixun_page5.html/",
				"https://www.925g.com/zixun_page6.html/",
				"https://www.925g.com/zixun_page7.html/",
				"https://www.925g.com/zixun_page8.html/",
				"https://www.925g.com/zixun_page9.html/",
			}
		},
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       2,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "app", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "app"
			}},
			"content": {ExcelHeader: "I", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "/api"
			}},
		},
		ListFields:            map[string]articleSpider.Field{},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 5,
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal, context.Background())

	s.Start()

}

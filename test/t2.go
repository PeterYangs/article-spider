package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
	"strings"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.925g.com/",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       2,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1", ConversionFunc: func(data string, resList map[string]string) string {

				return strings.Replace(strings.Replace(data, "《", " ", -1), "》", " ", -1)
			}},
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

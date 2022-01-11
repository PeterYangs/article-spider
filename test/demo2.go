package main

import (
	"github.com/PeterYangs/article-spider/v3/fileTypes"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "https://www.925g.com/",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       10,
		DetailFields: map[string]form.Field{
			"title": {ExcelHeader: "J", Types: fileTypes.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: fileTypes.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "app", ImagePrefix: func(form *form.Form, path string) string {

				return "app"
			}},
			"content": {ExcelHeader: "I", Types: fileTypes.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *form.Form, path string) string {

				return "/api"
			}},
		},
		ListFields: map[string]form.Field{

			"desc": {ExcelHeader: "K", Types: fileTypes.Text, Selector: "  a > div > p"},
		},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 5,
		//HttpProxy: "127.0.0.1:4780",
		//HttpHeader:
	})

	s.Start()

}

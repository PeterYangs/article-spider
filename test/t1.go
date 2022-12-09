package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.925g.com/",
		Channel:      "/gonglue/list_[PAGE].html",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       20,
		//DetailFields: map[string]articleSpider.Field{
		//	//"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
		//	"img": {ExcelHeader: "H", Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "[date:md]/[random:1-100]", ImageResizePercent: 10},
		//	//"content": {ExcelHeader: "I", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {
		//	//
		//	//	return "/api"
		//	//}, ImageDir: "[date:md]/[random:1-100]", ImageResizePercent: 10},
		//},
		ListFields: map[string]articleSpider.Field{
			"title": {Types: articleSpider.Text, Selector: " a > div > span"},
			"image": {Types: articleSpider.Image, Selector: " a > img", ImageDir: "[date:Y-m-d]"},
		},
		//CustomExcelHeader:     true,
		DetailCoroutineNumber: 1,
		FilterError:           true,
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal, context.Background())

	//s.SetImageDir("")
	//
	//s.SetSavePath("D:/down")

	s.Start()

}

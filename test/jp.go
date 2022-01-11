package main

import (
	"github.com/PeterYangs/article-spider/v3/fileTypes"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/spider"
	"time"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "https://www.rakuei.jp",
		Channel:      "/shopbrand/ct293/page[PAGE]/order/",
		ListSelector: "#r_categoryList > ul > li",
		HrefSelector: " div > a",
		PageStart:    1,
		Length:       1,
		ListFields: map[string]form.Field{
			"title": {ExcelHeader: "A", Types: fileTypes.Text, Selector: "div > div.detail > p.name"},
			"price": {ExcelHeader: "B", Types: fileTypes.Text, Selector: "div > div.detail > p.price > em"},
			//"img": {ExcelHeader: "E", Types: fileTypes.Image, Selector: " div > a   img", ImageDir: "cos-onsen", ImagePrefix: func(form *form.Form, path string) string {
			//
			//	return "app"
			//}},
		},
		CustomExcelHeader: true,
		HttpTimeout:       5 * time.Second,
	})

	s.Start()

}

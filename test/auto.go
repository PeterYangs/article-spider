package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:               "https://www.925g.com",
		Channel:            "/zixun/",
		ListSelector:       "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector:       "  a",
		NextSelector:       "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > ul > li:nth-child(11) > a",
		ListWaitSelector:   "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li:nth-child(1)",
		DetailWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1",
		Length:             3,
		DetailFields: map[string]form.Field{
			"title": {ExcelHeader: "J", Types: fileTypes.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: fileTypes.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonRightDiv.uk-float-right > div.single-sidebar > div > ul > li > a > img", ImageDir: "app", ImagePrefix: func(form *form.Form, path string) string {

				return "app"
			}},
		},
		ListFields: map[string]form.Field{
			"desc": {Types: fileTypes.Text, Selector: " a > div > p"},
		},
	})

	s.StartAuto()

}

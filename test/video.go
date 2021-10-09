package main

import (
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:               "https://www.ahjingcheng.com",
		Channel:            "/show/dongzuo--------[PAGE]---/",
		ListSelector:       "body > div:nth-child(5) > div > div.col-lg-wide-75.col-xs-1.padding-0 > div:nth-child(2) > div > div.stui-pannel_bd > ul > li",
		HrefSelector:       " div > a",
		PageStart:          1,
		Length:             1,
		MiddleHrefSelector: []string{"body > div:nth-child(3) > div > div.col-lg-wide-75.col-xs-1.padding-0 > div:nth-child(1) > div > div:nth-child(2) > div.stui-content__thumb > a"},
		//DetailFields: map[string]form.Field{
		//	"title": {ExcelHeader: "J", Types: fileTypes.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
		//	"img": {ExcelHeader: "H", Types: fileTypes.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "app", ImagePrefix: func(form *form.Form, path string) string {
		//
		//		return "app"
		//	}},
		//	"content": {ExcelHeader: "I", Types: fileTypes.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *form.Form, path string) string {
		//
		//		return "/api"
		//	}},
		//},

		//CustomExcelHeader:     true,
		DetailCoroutineNumber: 1,
	})

	s.Start()

}

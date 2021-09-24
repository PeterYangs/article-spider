package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(&form.Form{
		Host:         "https://www.duote.com",
		Channel:      "/sort/50_0_wdow_0_[PAGE]_.html",
		ListSelector: "body > div.wrap > div.box > div.main-left-box > div > div.bd > div > div.soft-info-lists > div",
		HrefSelector: "  a",
		PageStart:    1,
		Length:       1,
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.Text, Selector: "body > div.wrap.mt_5 > div > div.main-left-box > div.down-box > div.soft-name > div > h1"},
			"img": {Types: fileTypes.Image, Selector: "body > div.wrap.mt_5 > div > div.main-left-box > div.down-box > div.soft-name > img", ImageDir: "[singleField:title]", ImagePrefix: func(form *form.Form, path string) string {

				return "app"
			}},
		},
		ListFields: map[string]form.Field{

			"desc": {Types: fileTypes.Text, Selector: " div.sub-title"},
		},
	})

	s.Start()

}

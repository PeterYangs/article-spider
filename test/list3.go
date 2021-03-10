package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//爬多图
	f := form.Form{

		Host:             "https://www.duote.com",
		Channel:          "/sort/50_0_wdow_0_[PAGE]_.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.box > div.main-left-box > div > div.bd > div > div.soft-info-lists > div",
		ListHrefSelector: " a",
		DetailFields: map[string]form.Field{
			"list_img": {Types: fileTypes.ListImages, Selector: ".print-box img"},
		},
		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}

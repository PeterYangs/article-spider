package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.duote.com",
		Channel:          "/sort/50_0_wdow_0_[PAGE]_.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.box > div.main-left-box > div > div.bd > div > div.soft-info-lists > div",
		ListHrefSelector: " a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "body > div.wrap.mt_5 > div > div.main-left-box > div.down-box > div.soft-name > div > h1"},
		},
		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}

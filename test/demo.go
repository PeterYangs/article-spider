package main

import (
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
	})

	s.Start()

}

package main

import (
	"article-spider/chromedpSpider"
	"article-spider/fileTypes"
	"article-spider/form"
)

func main() {

	f := form.Form{

		Host:                "https://www.522gg.com",
		Channel:             "/game",
		Limit:               1,
		WaitForListSelector: "body > div:nth-child(5) > div > div.row.fn_mgsx10 > div",
		ListPath:            "/html/body/div[5]/div/div[2]/div",
		ListClickPath:       "/div/div/a",
		DetailFields:        map[string]form.Field{"title": {Types: fileTypes.SingleField, Selector: "body > div:nth-child(5) > div > div > div.col-xs12.col-sm12.col-md8.col-lg8 > div:nth-child(1) > div > div > div.info.w160 > div.l > h1"}},
		NextSelector:        "body > div:nth-child(5) > div > div:nth-child(3) > div > ul > li:nth-child(13) > a",
	}

	//chromedpSpider.GetList(f)

	chromedpSpider.Start(f)

}

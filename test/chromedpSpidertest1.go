package main

import (
	"article-spider/chromedpForm"
	"article-spider/chromedpSpider"
	"article-spider/fileTypes"
)

func main() {

	f := chromedpForm.Form{

		Host:            "https://www.522gg.com",
		Channel:         "/game",
		Limit:           5,
		WaitForListPath: "/html/body/div[5]/div/div[2]/div",
		ListPath:        "/html/body/div[5]/div/div[2]/div",
		ListClickPath:   "/div/div/a",
		DetailFields:    map[string]chromedpForm.Field{"title": {Types: fileTypes.SingleField, Path: "/html/body/div[5]/div/div/div[1]/div[1]/div/div/div[2]/div[1]/h1"}},
		NextPath:        "/html/body/div[5]/div/div[3]/div/ul/li[13]/a",
	}

	chromedpSpider.GetList(f)
}

package main

import (
	"article-spider/chromedpForm"
	"article-spider/chromedpSpider"
)

func main() {

	f := chromedpForm.Form{

		Host:                "https://down.gamersky.com",
		Channel:             "/Special/bigpc/",
		Limit:               5,
		WaitForListSelector: "body > div.Mid > div.Mid2 > ul > li",
		ListSelector:        "body > div.Mid > div.Mid2 > ul > li",
		ListClickSelector:   "div.tit > a",
	}

	chromedpSpider.GetList(f)
}

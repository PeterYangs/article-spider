package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.ccaaka.com",
		Channel:          "/cate/32.html?page=[PAGE]",
		Limit:            5,
		PageStart:        2,
		ListSelector:     "body > section > div > div > div.col-md-7.col-xs-12.article-container > div",
		ListHrefSelector: "div.col-md-8.col-xs-8 > a",
		DetailFields:     map[string]form.Field{"title": {Types: fileTypes.SingleField, SingleSelector: "body > section > div > div > div.col-md-7 > div:nth-child(2) > div > h3"}},
	}

	spider.Start(f)

}

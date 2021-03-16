package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//只爬列表
	f := form.Form{

		Host:             "https://store.shopping.yahoo.co.jp",
		Channel:          "/cossky/a5b3a5b9a5.html?page=[PAGE]#CentSrchFilter1",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "#itmlst > ul > li",
		ListHrefSelector: "",
		ListFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: " div:nth-child(2) > div:nth-child(1) > div > a > span"},
			"price": {Types: fileTypes.SingleField, Selector: " div:nth-child(2) > div:nth-child(3) > div > div:nth-child(1) > span"},
		},
		DetailMaxCoroutine: 2,
	}

	spider.Start(f)

}

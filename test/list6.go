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
		Channel:          "/sakuranokoi/5bb3a2a955a.html?page=[PAGE]#CentSrchFilter1",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "#itmlst > ul > li",
		ListHrefSelector: " div:nth-child(1) > div > div > a",
		//ListFields: map[string]form.Field{
		//	"title": {Types: fileTypes.SingleField, Selector: " div:nth-child(2) > div:nth-child(1) > div > a > span"},
		//	"price": {Types: fileTypes.SingleField, Selector: " div:nth-child(2) > div:nth-child(3) > div > div:nth-child(1) > span"},
		//},
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "#shpMain > div.gdColumns.gd3ColumnItem > div.gd3ColumnItem2 > div.mdItemName > p.elCatchCopy"},
			"img":   {Types: fileTypes.SingleImage, Selector: "#itmbasic > div.elMain > ul > li.elPanel.isNew > a > img"},
		},
		DetailMaxCoroutine: 2,
		ProxyAddress:       "socks5://127.0.0.1:4781",
		//ProxyAddress: "socks5://127.0.0.1:4777",
	}

	spider.Start(f)

}

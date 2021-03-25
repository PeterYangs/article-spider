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
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "#shpMain > div.gdColumns.gd3ColumnItem > div.gd3ColumnItem2 > div.mdItemName > p.elCatchCopy"},
			"img":   {Types: fileTypes.SingleImage, Selector: "#itmbasic > div.elMain > ul > li.elPanel.isNew > a > img"},
		},
		DetailMaxCoroutine: 2,
		ProxyAddress:       "socks5://127.0.0.1:4781",
	}

	spider.Start(f)

}

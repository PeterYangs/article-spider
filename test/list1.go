package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.cos-onsen.com",
		Channel:          "/product-list/46?view=new&page=[PAGE]",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "#main_container > article > div > div.page_contents.clearfix.categorylist_contents > div > div.itemlist_box.clearfix > ul > li",
		ListHrefSelector: "div > a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, SingleSelector: "#main_container > article > div > div.page_title > h1 > span.goods_name"},
			//"html":    {Types: fileTypes.OnlyHtml, SingleSelector: "body > section > div > div > div.col-md-7 > div:nth-child(2) > div"},
			//"image":   {Types: fileTypes.SingleImage, SingleSelector: "body > section > div > div > div.col-md-3 > div > div.qrcode-panel.common-panel > div:nth-child(1) > img", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
			//"content": {Types: fileTypes.HtmlWithImage, SingleSelector: "#detail-content", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
		},
	}

	spider.Start(f)

}

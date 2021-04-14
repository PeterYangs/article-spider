package main

import (
	"github.com/PeterYangs/article-spider/fileTypes"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/article-spider/spider"
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
			"title": {Types: fileTypes.SingleField, Selector: "#main_container > article > div > div.page_title > h1 > span.goods_name"},
			//"html":    {Types: fileTypes.OnlyHtml, Selector: "body > section > div > div > div.col-md-7 > div:nth-child(2) > div"},
			"image": {Types: fileTypes.SingleImage, Selector: "body > section > div > div > div.col-md-3 > div > div.qrcode-panel.common-panel > div:nth-child(1) > img", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
			//"content": {Types: fileTypes.HtmlWithImage, Selector: "#detail-content", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
		},
	}

	spider.Start(f)

}

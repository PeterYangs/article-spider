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
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "body > section > div > div > div.col-md-7 > div:nth-child(2) > div > h3"},
			//"image": {Types: fileTypes.SingleImage, Selector: "body > section > div > div > div.col-md-3 > div > div.info-panel > ul > li:nth-child(1) > img.left"},
			//body > section > div > div > div.col-md-3 > div > div.info-panel > ul > li:nth-child(1) > img.left
			//"html":    {Types: fileTypes.OnlyHtml, Selector: "body > section > div > div > div.col-md-7 > div:nth-child(2) > div"},
			//"image":   {Types: fileTypes.SingleImage, Selector: "body > section > div > div > div.col-md-3 > div > div.qrcode-panel.common-panel > div:nth-child(1) > img", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
			//"content": {Types: fileTypes.HtmlWithImage, Selector: "#detail-content", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
		},
		ListFields: map[string]form.Field{
			"image": {Types: fileTypes.SingleImage, Selector: "div > img"},
		},
	}

	spider.Start(f)

}

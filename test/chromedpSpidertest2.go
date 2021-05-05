package main

import (
	"article-spider/chromedpSpider"
	"article-spider/fileTypes"
	"article-spider/form"
)

func main() {

	f := form.Form{

		Host:                "https://down.gamersky.com",
		Channel:             "/Special/bigpc/",
		Limit:               2,
		WaitForListSelector: "body > div.Mid > div.Mid2 > ul > li:nth-child(1)",
		ListSelector:        "body > div.Mid > div.Mid2 > ul > li",
		ListHrefSelector:    "div.tit > a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "body > div.Mid > div.Mid2 > div.Mid2_L > div.Mid2L_ctt.block > div.Mid2L_actdl2 > div.tit"},
			"img":   {Types: fileTypes.SingleImage, Selector: "body > div.Mid > div.Mid2 > div.Mid2_L > div.Mid2L_ctt.block > div.Mid2L_actdl2 > div.game > div.img > img", ImageDir: "demo"},
		},
		//ListFields: map[string]form.Field{
		//	"img":{Types: fileTypes.SingleImage,Selector: " div > div > a > div.img > img"},
		//},
		NextSelector: "#pe100_page_jdgame > a.p1.nexe",
	}

	//chromedpSpider.GetList(f)

	chromedpSpider.Start(f)

}

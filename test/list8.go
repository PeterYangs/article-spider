package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//只爬列表
	f := form.Form{

		Host:             "https://www.doyo.cn",
		Channel:          "/game/2-1-[PAGE].html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.mobile_game_wrap.w1168.clearfix.bg > div > div > div.tab_box > div > div > ul > li",
		ListHrefSelector: " div > a:nth-child(1)",
		DetailFields: map[string]form.Field{
			"img":         {Types: fileTypes.SingleImage, Selector: " body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.img_logo > img"},
			"title":       {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.info > h1"},
			"content":     {Types: fileTypes.HtmlWithImage, Selector: "#hiddenDetail > div"},
			"screenshots": {Types: fileTypes.ListImages, Selector: "#slider3 > ul img"},
			"size":        {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.detail_info > div.info.clearfix > span:nth-child(1) > em"},
		},
		//DetailFields: map[string]form.Field{
		//	"title": {Types: fileTypes.SingleField, Selector: "#shpMain > div.gdColumns.gd3ColumnItem > div.gd3ColumnItem2 > div.mdItemName > p.elCatchCopy"},
		//	"img":   {Types: fileTypes.SingleImage, Selector: "#itmbasic > div.elMain > ul > li.elPanel.isNew > a > img"},
		//},
		DetailMaxCoroutine: 5,
		//ProxyAddress:       "socks5://127.0.0.1:4781",
		//ProxyAddress: "socks5://127.0.0.1:4777",
	}

	spider.Start(f)
}

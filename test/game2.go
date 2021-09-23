package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
	"github.com/PeterYangs/article-spider/v2/tool"
)

func main() {

	f := form.Form{

		Host:             "https://www.93zg.com",
		Channel:          "/game/list-0-[PAGE]-0-0-0-0-0-0.html",
		Limit:            2,
		PageStart:        1,
		ListSelector:     "body > div.yxk_n1.wrap > div.game-main > div.gameblock > ul > li",
		ListHrefSelector: " div > a",
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.SingleField, Selector: "body > div:nth-child(3) > div > div > div.glgamedown > dl > dt > h3", ExcelHeader: "G"},
			// /api/uploads/news20210325/4620210325/88e8a06664b249bf90fe12ccba084f89.jpg
			"content":     {Types: fileTypes.HtmlWithImage, Selector: "body > div:nth-child(3) > div.gp-left > div > div.glgametab.game-tj > div.gtabcont.tjcont > div:nth-child(1) > div.gtabtext", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "game2[date:md]/[random:1-100]"},
			"desc":        {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: tool.GetDescGame},
			"keyword":     {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":         {Types: fileTypes.SingleImage, Selector: "body > div:nth-child(3) > div.gp-left > div > div.glgamedown > div.gldimg > img", ExcelHeader: "F", ImageDir: "game2[date:md]/[random:1-100]"},
			"type":        {Types: fileTypes.Fixed, Selector: "2", ExcelHeader: "L"},
			"size":        {Types: fileTypes.SingleField, Selector: "body > div:nth-child(3) > div.gp-left > div > div.glgamedown > dl > dd:nth-child(2) > p:nth-child(1) > span:nth-child(2) > em", ExcelHeader: "M"},
			"screenshots": {Types: fileTypes.ListImages, Selector: "body > div:nth-child(3) > div.gp-left > div > div.glgametab.game-tj > div.gtabcont.tjcont > div:nth-child(1) > div.gtabrotation > div > ul img", ExcelHeader: "N", ConversionFormatFunc: tool.ConversionGame, ImageDir: "game2[date:md]/[random:1-100]"},
			"category_id": {Types: fileTypes.Fixed, Selector: "10", ExcelHeader: "D"},
		},

		DetailMaxCoroutine: 1,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		CustomExcelHeader:  true,
		//DisableDebug: true,
		//ProxyAddress: "http://127.0.0.1:4780",
	}

	spider.Start(f)

}

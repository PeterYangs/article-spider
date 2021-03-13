package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//爬多图
	f := form.Form{

		Host:             "http://ly.8090.com",
		Channel:          "/gongl/[PAGE].html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "#game_center_right > div > div > ul > li",
		ListHrefSelector: "a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "#game_center_right > div > div.news_con_txt > div.game_read_tit > h1"},
		},
		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}

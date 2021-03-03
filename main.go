package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.weixz.com",
		Channel:          "/gamexz/list_[PAGE]-0.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.GameList.wd1200.mt-20px > ul > li",
		ListHrefSelector: "div.GameListIcon > a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},
			"image": {Types: fileTypes.SingleImage, SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoIcon > img", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
		},
	}

	spider.Start(f)

}

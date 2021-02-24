package main

import (
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.ccaaka.com",
		Channel:          "/cate/32.html?page=[PAGE]",
		Limit:            30,
		PageStart:        2,
		ListSelector:     "body > section > div > div > div.col-md-7.col-xs-12.article-container > div",
		ListHrefSelector: "div.col-md-8.col-xs-8 > a",
	}

	spider.GetList(f)

}

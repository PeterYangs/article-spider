package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.upwork.com",
		Channel:          "/search/jobs/?page=[PAGE]&sort=recency",
		Limit:            3,
		PageStart:        1,
		ListSelector:     "#layout > div > div.air-card.responsive-search-card.p-0-top-bottom.m-0-top-bottom.d-flex-mobile-app.height-100-mobile-app > div > div > div > section",
		ListHrefSelector: " job-tile-responsive > div.p-sm-top-bottom > div.clearfix > h4 > a",
		ListFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "a"},
		},
		DetailMaxCoroutine: 1,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		//CustomExcelHeader:  true,
		//DisableDebug: true,
	}

	spider.Start(f)
}

package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:               "https://www.ahjingcheng.com",
		Channel:            "/show/dongzuo--------[PAGE]---/",
		ListSelector:       "body > div:nth-child(5) > div > div.col-lg-wide-75.col-xs-1.padding-0 > div:nth-child(2) > div > div.stui-pannel_bd > ul > li",
		HrefSelector:       " div > a",
		PageStart:          1,
		Length:             2,
		MiddleHrefSelector: []string{"body > div:nth-child(3) > div > div.col-lg-wide-75.col-xs-1.padding-0 > div:nth-child(1) > div > div:nth-child(2) > div.stui-content__thumb > a"},
		DetailFields: map[string]form.Field{
			"url": {Types: fileTypes.Regular, Selector: `"url":"([0-9A-Za-z/\\._:]+)","url_next"`, RegularIndex: 1},
		},

		//CustomExcelHeader:     true,
		DetailCoroutineNumber: 1,
		HttpHeader: map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		},
	})

	s.Start()

}

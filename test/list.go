package main

import (
	"fmt"
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "https://www.925g.com",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       10,
		ListFields: map[string]form.Field{

			"title": {ExcelHeader: "K", Types: fileTypes.Text, Selector: " a > div > span"},
		},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 2,
		ResultCallback: func(item map[string]string, form *form.Form) {

			for s2, s3 := range item {

				fmt.Println(s2, ":", s3)

			}

		},
	})

	s.Start()

}

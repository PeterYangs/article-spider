package main

import (
	"context"
	"fmt"
	articleSpider "github.com/PeterYangs/article-spider/v4"
	"time"
)

func main() {

	cxt, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	s := articleSpider.NewSpider(articleSpider.Form{

		Host:         "https://www.925g.com",
		Channel:      "/zixun/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: "  a",
		//下一页选择器
		AutoNextSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > ul > li:nth-child(11) > a",
		//列表等待选择器
		//AutoListWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li:nth-child(1)",
		//详情等待选择器
		AutoDetailWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1",
		Length:                 3,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"content": {ExcelHeader: "H", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImageDir: "app", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "app"
			}},
		},

		//cookie
		HttpHeader: map[string]string{
			"cookie": "user_cookie=Vmod7XlkHN; UM_distinctid=17b805b421c1e0-0005d3dc1ac8ea-c343365-1fa400-17b805b421dda7; url_data=https://www.925g.com/zixun/,https://www.925g.com/; PHPSESSID=3m0ee50ba4r40jq3fleob2n71i; CNZZDATA1278942394=1852940385-1600066493-%7C1635143024; Hm_lvt_46233f03c62deb1e98a07bf1e1708415=1634807167,1634887947,1634955841,1635153418; Hm_lpvt_46233f03c62deb1e98a07bf1e1708415=1635153430",
		},
	}, articleSpider.Auto, cxt)

	err := s.Start()

	if err != nil {

		fmt.Println(err)
	}

}

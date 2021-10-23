package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "https://www.925g.com",
		Channel:      "/zixun/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: "  a",
		//下一页选择器
		NextSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > ul > li:nth-child(11) > a",
		//列表等待选择器
		ListWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li:nth-child(1)",
		//详情等待选择器
		DetailWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1",
		Length:             3,
		DetailFields: map[string]form.Field{
			"title": {ExcelHeader: "J", Types: fileTypes.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: fileTypes.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonRightDiv.uk-float-right > div.single-sidebar > div > ul > li > a > img", ImageDir: "app", ImagePrefix: func(form *form.Form, path string) string {

				return "app"
			}},
		},
		ListFields: map[string]form.Field{
			"desc": {Types: fileTypes.Text, Selector: " a > div > p"},
		},
		//cookie
		AutoCookieString: "PHPSESSID=lchdl81cdggfcbp694gf3894lh; user_cookie=9X18yQilnW; url_data=https://www.925g.com/; UM_distinctid=17cab2a6d0a8af-0933cba3984f97-c343365-1fa400-17cab2a6d0be92; CNZZDATA1278942394=499160135-1634952571-%7C1634952571; Hm_lvt_46233f03c62deb1e98a07bf1e1708415=1634959257; Hm_lpvt_46233f03c62deb1e98a07bf1e1708415=1634959383",
	})

	s.StartAuto()

}

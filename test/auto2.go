package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "http://a.4399.cn",
		Channel:      "/game-hot-p-1.html",
		ListSelector: "#j-hotGamelist > li",
		HrefSelector: " a.m_game",
		//下一页选择器
		NextSelector: "body > div.wrapper > div > div > div > div.ks_bd > div > div > a.ks_next",
		//列表等待选择器
		ListWaitSelector: "#j-hotGamelist > li:nth-child(1)",
		//详情等待选择器
		DetailWaitSelector: "body > div.wrapper > div.a_grid_detail.mt > div.m_game_detail.clearfix > div.m_game_name > h1",
		Length:             3,
		DetailFields: map[string]form.Field{
			"title": {ExcelHeader: "J", Types: fileTypes.Text, Selector: "body > div.wrapper > div.a_grid_detail.mt > div.m_game_detail.clearfix > div.m_game_name > h1"},
		},
		//ListFields: map[string]form.Field{
		//	"desc": {Types: fileTypes.Text, Selector: " a > div > p"},
		//},
		//cookie
		//HttpHeader: map[string]string{
		//	"cookie": "CNZZDATA1278942394=499160135-1634952571-%7C1634952571; Hm_lvt_46233f03c62deb1e98a07bf1e1708415=1634959257; Hm_lpvt_46233f03c62deb1e98a07bf1e1708415=1634959383; PHPSESSID=lchdl81cdggfcbp694gf3894lh; user_cookie=9X18yQilnW; url_data=https://www.925g.com/; UM_distinctid=17cab2a6d0a8af-0933cba3984f97-c343365-1fa400-17cab2a6d0be92; UM_distinctid=17cabec696cd0e-07f763e7354ff6-c343365-1fa400-17cabec696ded4; Hm_lvt_46233f03c62deb1e98a07bf1e1708415=1634971970; Hm_lpvt_46233f03c62deb1e98a07bf1e1708415=1634972020",
		//},
	})

	s.StartAuto()

}

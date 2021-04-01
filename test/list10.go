package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
	"encoding/json"
	"fmt"
	"github.com/PeterYangs/tools"
)

func main() {

	f := form.Form{

		Host:             "https://www.doyo.cn",
		Channel:          "/game/2-1-[PAGE].html",
		Limit:            1,
		PageStart:        1,
		ListSelector:     "body > div.mobile_game_wrap.w1168.clearfix.bg > div > div > div.tab_box > div > div > ul > li",
		ListHrefSelector: " div > a:nth-child(1)",
		DetailFields: map[string]form.Field{
			//"img": {Types: fileTypes.SingleImage, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.img_logo > img", ExcelHeader: "A",ImageDir:"[singleField:title]"},
			"title": {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.info > h1"},
			//"content":     {Types: fileTypes.HtmlWithImage, Selector: "#hiddenDetail > div", ExcelHeader: "C"},
			"screenshots": {Types: fileTypes.ListImages, Selector: "#slider3 > ul img", ExcelHeader: "D", ImageDir: "[singleField:title]"},
			//"size":        {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.detail_info > div.info.clearfix > span:nth-child(1) > em", ExcelHeader: "E"},
		},
		DetailMaxCoroutine: 5,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		//CustomExcelHeader:  true,
		//DisableDebug: true,
	}

	spider.Start(f)
}

//转json格式
func Conversion(data string, list map[string]string) string {

	var jsons []map[string]string

	array := tools.Explode(",", data)

	for _, v := range array {

		jsons = append(jsons, map[string]string{"img": v, "name": ""})

	}

	jsonStr, err := json.Marshal(jsons)

	if err != nil {

		fmt.Println(err)

		return ""
	}

	return string(jsonStr)

}

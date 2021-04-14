package main

import (
	"encoding/json"
	"fmt"
	"github.com/PeterYangs/article-spider/fileTypes"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/article-spider/spider"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {

	f := form.Form{

		Host:             "https://www.weixz.com",
		Channel:          "/rjxz/list_[PAGE].html",
		Limit:            60,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.GameList.wd1200.mt-20px > ul > li",
		ListHrefSelector: " div.GameListInfo > div.GameListTitle > a",
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1", ExcelHeader: "G"},
			// /api/uploads/news20210325/4620210325/88e8a06664b249bf90fe12ccba084f89.jpg
			"content":     {Types: fileTypes.HtmlWithImage, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentTexts > div.mobileGamesContentText", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "app[date:md]/[random:1-100]"},
			"desc":        {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: getDescApp},
			"keyword":     {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":         {Types: fileTypes.SingleImage, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoIcon > img", ExcelHeader: "F", ImageDir: "app[date:md]/[random:1-100]"},
			"type":        {Types: fileTypes.Fixed, Selector: "2", ExcelHeader: "L"},
			"size":        {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > p:nth-child(2) > span:nth-child(4)", ExcelHeader: "M"},
			"screenshots": {Types: fileTypes.ListImages, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentImgs > div.mobileGamesContentImg > div > div.swiper-wrapper img", ExcelHeader: "N", ConversionFormatFunc: conversionApp, ImageDir: "app[date:md]/[random:1-100]"},
			"category_id": {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > p:nth-child(2) > span:nth-child(2)", ExcelHeader: "D", ConversionFormatFunc: getCategoryApp},
		},

		DetailMaxCoroutine: 2,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		CustomExcelHeader:  true,
		//DisableDebug: true,
		//ProxyAddress: "http://127.0.0.1:4780",
	}

	spider.Start(f)

}

func getDescApp(data string, resList map[string]string) string {

	if data == "" {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resList["content"]))

		if err != nil {

			return ""
		}

		return doc.Text()

	}

	return data
}

//转json格式
func conversionApp(data string, list map[string]string) string {

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

func getCategoryApp(data string, list map[string]string) string {

	switch data {

	case "影音播放":

		return "20"

	case "社交聊天":

		return "21"

	case "新闻阅读":

		return "23"

	case "摄影美化":

		return "24"

	case "金融理财":

		return "27"

	case "购物支付":

		return "26"

	case "办公学习":

		if tools.Mt_rand(1, 2)%2 == 0 {

			return "25"
		}

		return "31"

	case "生活服务":

		if tools.Mt_rand(1, 2)%2 == 0 {

			return "28"
		}

		return "30"

	case "主题桌面":

		return "22"

	case "系统安全":

		return "22"

	case "地图导航":

		return "29"

	}

	return "30"
}

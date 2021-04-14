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
		Channel:          "/gamexz/list_[PAGE]-0.html",
		Limit:            200,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.GameList.wd1200.mt-20px > ul > li",
		ListHrefSelector: " div.GameListInfo > div.GameListTitle > a",
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1", ExcelHeader: "G"},
			// /api/uploads/news20210325/4620210325/88e8a06664b249bf90fe12ccba084f89.jpg
			"content":     {Types: fileTypes.HtmlWithImage, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentTexts > div.mobileGamesContentText", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "game[date:md]/[random:1-100]"},
			"desc":        {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: getDescGame},
			"keyword":     {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":         {Types: fileTypes.SingleImage, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoIcon > img", ExcelHeader: "F", ImageDir: "game[date:md]/[random:1-100]"},
			"type":        {Types: fileTypes.Fixed, Selector: "2", ExcelHeader: "L"},
			"size":        {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > p:nth-child(2) > span:nth-child(4)", ExcelHeader: "M"},
			"screenshots": {Types: fileTypes.ListImages, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentImgs > div.mobileGamesContentImg > div > div.swiper-wrapper img", ExcelHeader: "N", ConversionFormatFunc: conversionGame, ImageDir: "game[date:md]/[random:1-100]"},
			"category_id": {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > p:nth-child(2) > span:nth-child(2)", ExcelHeader: "D", ConversionFormatFunc: getCategory},
		},

		//DetailMaxCoroutine: 3,
		HttpHeader:        map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		CustomExcelHeader: true,
		//DisableDebug: true,
		//ProxyAddress: "http://127.0.0.1:4780",
	}

	spider.Start(f)

}

func getDescGame(data string, resList map[string]string) string {

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
func conversionGame(data string, list map[string]string) string {

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

func getCategory(data string, list map[string]string) string {

	switch data {

	case "角色扮演":

		return "2"

	case "动作格斗":

		return "3"

	case "休闲益智":

		return "8"

	case "飞行射击":

		return "6"

	case "冒险解密":

		return "10"

	case "策略塔防":

		return "7"

	case "赛车竞速":

		return "9"

	case "棋牌卡牌":

		return "4"

	case "音乐游戏":

		return "12"

	case "模拟经营":

		return "5"

	case "体育竞技":

		return "11"

	case "二次元养成":

		return "13"

	}

	return "2"
}

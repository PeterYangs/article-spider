package main

import (
	"encoding/json"
	"github.com/PeterYangs/article-spider/apiSpider"
	"github.com/PeterYangs/article-spider/fileTypes"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/article-spider/tool"
	"github.com/PeterYangs/tools"
	"log"
	"regexp"
	"strconv"
)

func main() {

	f := form.Form{

		Host:      "http://app.cnfol.com",
		Channel:   "/dataapi/index.php/futuresdata/getCatArticle?catid=133&contid=28828722&page=[PAGE]&limit=23&jsoncallback=jQuery1830549977029189465_1618989236347&_=1618989248636",
		Limit:     40,
		PageStart: 1,
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "body > div.allCnt > div.artMain.mBlock > h3.artTitle", ExcelHeader: "G"},
			// /api/uploads/news20210325/4620210325/88e8a06664b249bf90fe12ccba084f89.jpg
			"content": {Types: fileTypes.HtmlWithImage, Selector: "body > div.allCnt > div.artMain.mBlock > div.Article", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "stock3/[random:1-100]", DefaultImg: getDefaultImg2},
			"desc":    {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: tool.GetDescGame},
			"keyword": {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":     {Types: fileTypes.Fixed, Selector: "1", ExcelHeader: "F"},
		},
		//ListFields: map[string]form.Field{
		//	"image": {Types: fileTypes.SingleImage, Selector: "div > img"},
		//},
		CustomExcelHeader:  true,
		DetailMaxCoroutine: 5,
		ApiConversion: func(result string) []string {

			re1 := regexp.MustCompile(`^jQuery[0-9_]+\((.*?)\)$`).FindStringSubmatch(result)

			var jsons []map[string]interface{}

			err := json.Unmarshal([]byte(re1[1]), &jsons)

			if err != nil {

				log.Print(err)

				return []string{}

			}

			var linkList []string

			for _, m := range jsons {

				linkList = append(linkList, m["Url"].(string))
			}

			return linkList
		},
		HttpHeader: map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
	}

	apiSpider.Start(f)

}

func getDefaultImg2(form form.Form, item form.Field) string {

	number := tools.MtRand(1, 40)

	return "/api/uploads/stock/" + strconv.Itoa((int)(number)) + ".jpg"

}

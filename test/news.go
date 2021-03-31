package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {

	f := form.Form{

		Host:             "http://www.gj078.cn",
		Channel:          "/sports/index_[PAGE].html",
		Limit:            1,
		PageStart:        1,
		ListSelector:     "#recent-content > div",
		ListHrefSelector: " div > a",
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.SingleField, Selector: "#main > article > header > h1", ExcelHeader: "G"},
			// /api/uploads/news20210325/4620210325/88e8a06664b249bf90fe12ccba084f89.jpg
			"content": {Types: fileTypes.HtmlWithImage, Selector: "#main > article > div.entry-content", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "news/[random:1-100]"},
			"desc":    {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: getDesc},
			"keyword": {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
		},
		ListFields: map[string]form.Field{
			"img": {Types: fileTypes.SingleImage, Selector: " div > a > div > img", ExcelHeader: "F", ImageDir: "news/[random:1-100]"},
		},
		//DetailMaxCoroutine: 5,
		HttpHeader:        map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		CustomExcelHeader: true,
		//DisableDebug: true,
	}

	spider.Start(f)

}

func getDesc(data string, resList map[string]string) string {

	if data == "" {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resList["content"]))

		if err != nil {

			return ""
		}

		return doc.Text()

	}

	return data
}

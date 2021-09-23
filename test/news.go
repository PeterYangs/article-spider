package main

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/spider"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {

	f := form.Form{

		Host:             "http://news.4399.com",
		Channel:          "/shouyou/6438_[PAGE].html",
		Limit:            2,
		PageStart:        1,
		ListSelector:     "body > div.main.mb10.clearfix > div.leftbar > div.tabC1 > ul > li",
		ListHrefSelector: " div.top_t > a",
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.SingleField, Selector: "body > div.wp.cf > div.w700.fl > h1", ExcelHeader: "G"},
			// /api/uploads/news20210325/4620210325/88e8a06664b249bf90fe12ccba084f89.jpg
			"content": {Types: fileTypes.HtmlWithImage, Selector: "body > div.wp.cf > div.w700.fl > div.content", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "news/[random:1-100]"},
			"desc":    {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: getDesc},
			"keyword": {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
		},
		ListFields: map[string]form.Field{
			"img": {Types: fileTypes.SingleImage, Selector: " div.l_li_img > a > img", ExcelHeader: "F", ImageDir: "news/[random:1-100]"},
		},
		DetailMaxCoroutine: 2,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		CustomExcelHeader:  true,
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

		return tools.SubStr(doc.Text(), 0, 65)

	}

	return data
}

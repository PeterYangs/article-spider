package main

import (
	"github.com/PeterYangs/article-spider/v3/fileTypes"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/mode"
	"github.com/PeterYangs/article-spider/v3/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "https://www.gameinformer.com/",
		Channel:      "",
		ListSelector: "body > div.dialog-off-canvas-main-canvas > div.layout-container.front-page > main > div.content-wrapper > div.layout-center.load-more-content > div > div",
		HrefSelector: " div > div > div.teaser-left > div > div > a",
		//下一页选择器
		NextSelector: "body > div.dialog-off-canvas-main-canvas > div.layout-container.front-page > main > div.content-wrapper > div.load-more > div.load-more-wrapper",
		//列表等待选择器
		ListWaitSelector: "body > div.dialog-off-canvas-main-canvas > div.layout-container.front-page > main > div.content-wrapper > div.layout-center.load-more-content > div > div:nth-child(1)",
		//详情等待选择器
		DetailWaitSelector: "body > div.dialog-off-canvas-main-canvas > div > main > div > div > div.region.region-content > div.node.node--type-article.node--view-mode-full.ds-standard-article.clearfix > div.ds-header > div.field.field--name-node-title.field--type-ds.field--label-hidden.gi5-node-title.gi5-ds.field__item > h1",
		Length:             3,
		DetailFields: map[string]form.Field{
			"title": {ExcelHeader: "G", Types: fileTypes.Text, Selector: "body > div.dialog-off-canvas-main-canvas > div > main > div > div > div.region.region-content > div.node.node--type-article.node--view-mode-full.ds-standard-article.clearfix > div.ds-header > div.field.field--name-node-title.field--type-ds.field--label-hidden.gi5-node-title.gi5-ds.field__item > h1"},
			"content": {Types: fileTypes.HtmlWithImage, Selector: "body > div.dialog-off-canvas-main-canvas > div > main > div > div > div.region.region-content > div.node.node--type-article.node--view-mode-full.ds-standard-article.clearfix > div.ds-wrapper.ds-content-wrapper > div.ds-main > div.clearfix.text-formatted.field.field--name-body.field--type-text-with-summary.field--label-hidden.gi5-body.gi5-text-with-summary.field__item", ExcelHeader: "E", ImagePrefix: func(form *form.Form, imageName string) string {

				return "/api/uploads"
			}, ImageDir: "game[date:md]/[random:1-100]"},
			"desc":    {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H"},
			"keyword": {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":     {Types: fileTypes.Image, Selector: "body > div.dialog-off-canvas-main-canvas > div > main > div > div > div.region.region-content > div.ds-full-width > div.field.field--name-field-header.field--type-entity-reference-revisions.field--label-hidden.gi5-field-header.gi5-entity-reference-revisions.field__item > div > div > div > div > img", ExcelHeader: "F", ImageDir: "game[date:md]/[random:1-100]"},
			//"type":    {Types: fileTypes.Fixed, Selector: "2", ExcelHeader: "L"},
			//"size":    {Types: fileTypes.SingleField, Selector: "#dinfo > p.base > i:nth-child(3)", ExcelHeader: "M"},
		},

		//cookie
		HttpHeader: map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
			"cookie":     "lang=zh-CN; lang=zh-CN; lang=zh-CN; _ga=GA1.1.1532009431.1641283813; UM_distinctid=17e24238a22739-0fc0995e9cfdad-c343365-1fa400-17e24238a2352e; guid=cff3a072d6ca30b80ee729f0884a8596f65d9a28; CNZZDATA5291371=cnzz_eid%3D1358048227-1641278212-%26ntime%3D1641338428; CNZZDATA1278599438=848177868-1641279863-%7C1641340242; Hm_lvt_ddaa34551214df42d1e5f11974f9f744=1641283822,1641346329; _csrf=3f62bc78510faa5fecfbf404cbee0ec56d1c4f3a; s_a=1; _ga_76F07DJEB4=GS1.1.1641346328.3.1.1641346978.0; Hm_lpvt_ddaa34551214df42d1e5f11974f9f744=1641346980",
		},
		NextPageMode:          mode.LoadMore,
		CustomExcelHeader:     true,
		AutoDetailForceNewTab: true, //详情页强制打开新窗口
	})

	s.StartAuto()

}

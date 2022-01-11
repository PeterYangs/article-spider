package main

import (
	"context"
	"github.com/PeterYangs/article-spider/v3/fileTypes"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/mode"
	"github.com/PeterYangs/article-spider/v3/spider"
	"github.com/chromedp/chromedp"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:         "https://www.btcfans.com",
		Channel:      "/zh-cn/wallet",
		ListSelector: "body > div.page-width.page-content > div.main-content > div > div.module-content > ul > li",
		HrefSelector: " a",
		//下一页选择器
		NextSelector: "body > div.page-width.page-content > div.main-content > div > div.module-content > a",
		//列表等待选择器
		ListWaitSelector: "body > div.page-width.page-content > div.main-content > div > div.module-content > ul > li:nth-child(1)",
		//详情等待选择器
		DetailWaitSelector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.info_1 > div.name > div.name-ch",
		Length:             4,
		DetailFields: map[string]form.Field{
			"title": {ExcelHeader: "G", Types: fileTypes.Text, Selector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.info_1 > div.name > div.name-ch"},
			"content": {Types: fileTypes.HtmlWithImage, Selector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.wallet-des > div > p", ExcelHeader: "E", ImagePrefix: func(form *form.Form, imageName string) string {

				return "/api/uploads"
			}, ImageDir: "game[date:md]/[random:1-100]"},
			"desc":    {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H"},
			"keyword": {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":     {Types: fileTypes.Image, Selector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.info_1 > div.cover > img", ExcelHeader: "F", ImageDir: "game[date:md]/[random:1-100]"},
			"type":    {Types: fileTypes.Fixed, Selector: "2", ExcelHeader: "L"},
			//"size":    {Types: fileTypes.SingleField, Selector: "#dinfo > p.base > i:nth-child(3)", ExcelHeader: "M"},
		},

		//cookie
		HttpHeader: map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
			"cookie":     "lang=zh-CN; lang=zh-CN; lang=zh-CN; _ga=GA1.1.1532009431.1641283813; UM_distinctid=17e24238a22739-0fc0995e9cfdad-c343365-1fa400-17e24238a2352e; guid=cff3a072d6ca30b80ee729f0884a8596f65d9a28; CNZZDATA5291371=cnzz_eid%3D1358048227-1641278212-%26ntime%3D1641338428; CNZZDATA1278599438=848177868-1641279863-%7C1641340242; Hm_lvt_ddaa34551214df42d1e5f11974f9f744=1641283822,1641346329; _csrf=3f62bc78510faa5fecfbf404cbee0ec56d1c4f3a; s_a=1; _ga_76F07DJEB4=GS1.1.1641346328.3.1.1641346978.0; Hm_lpvt_ddaa34551214df42d1e5f11974f9f744=1641346980",
		},
		NextPageMode:      mode.LoadMore,
		CustomExcelHeader: true,
		AutoPrefixEvent: func(chromedpCtx context.Context) {

			chromedp.Run(
				chromedpCtx,

				chromedp.Click("#Alert > div > div.sure_btn", chromedp.ByQuery),
			)

		},
	})

	s.StartAuto()

}

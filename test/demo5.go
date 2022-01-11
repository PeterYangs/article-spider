package main

import (
	"github.com/PeterYangs/article-spider/v3/fileTypes"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/spider"
)

func main() {

	s := spider.NewSpider()
	s.LoadForm(form.CustomForm{
		Host:         "https://www.cgcosplay.jp",
		Channel:      "/product-list?page=[PAGE]",
		ListSelector: "#inner_main_container > section > div > div.page_contents.clearfix.alllist_contents > div > div.itemlist_box.tiled_list_box.layout_photo > div > ul > li",
		HrefSelector: " div > a",
		PageStart:    1,
		Length:       10,
		ListFields: map[string]form.Field{
			"title": {ExcelHeader: "A", Types: fileTypes.Text, Selector: "div > a > div > div.list_item_data > p.item_name > span.goods_name"},
			"price": {ExcelHeader: "B", Types: fileTypes.Text, Selector: "div > a > div > div.list_item_data > div > div > p.selling_price > span.figure"},
			"img": {ExcelHeader: "C", Types: fileTypes.Image, Selector: "  div > a > div > div.list_item_photo > div > div", ImageDir: "cgcosplay_image", ImagePrefix: func(form *form.Form, path string) string {

				return "cgcosplay_image"
			}},
		},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 10,
		LazyImageAttrName:     "data-src",
		HttpProxy:             "http://127.0.0.1:4780",
	})

	s.Start()

}

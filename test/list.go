package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.cosonsen.co.jp",
		Channel:          "/products_new.html?page=[PAGE]&cPath=69&disp_order=6",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "#newProductsDefault > div.products_new > div",
		ListHrefSelector: "div.discount_img > a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, SingleSelector: "#contentMainWrapper > div.con_right > div.banner > form > div.product_detail > div.product_detail_param > div:nth-child(2) > h2"},
			"img":   {Types: fileTypes.SingleImage, SingleSelector: "#showleft > div > img"},
		},

		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}

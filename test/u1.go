package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host: "https://www.925g.com/",

		DetailUrls: []string{

			"https://www.925g.com/gonglue/138499.html",
			"https://www.925g.com/gonglue/138498.html",
			"https://www.925g.com/gonglue/138497.html",
			"https://www.925g.com/gonglue/138496.html",
			"https://www.925g.com/gonglue/138495.html",
			"https://www.925g.com/gonglue/138494.html",
		},
		DetailFields: map[string]articleSpider.Field{
			"title": {Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img":   {Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "[date:md]/[random:1-100]"},
			"content": {Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "/api"
			}, ImageDir: "[date:md]/[random:1-100]"},
		},
		DetailCoroutineNumber: 3,
		FilterError:           true,
	}

	s := articleSpider.NewSpider(f, articleSpider.Url, context.Background())

	s.Start()

}

package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
	"strings"
	"time"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.xyzs.com",
		Channel:      "/app/soft/index_[PAGE].html",
		ListSelector: "body > div.wrapper > section.aplist > ul > li",

		PageStart: 51,
		Length:    100,

		ListFields: map[string]articleSpider.Field{
			"title": {Types: articleSpider.Text, Selector: " a > p.name"},
		},

		DetailCoroutineNumber: 1,
		FilterError:           true,
		Filter: func(m map[string]string) bool {

			defer time.Sleep(100 * time.Millisecond)

			if strings.Contains(m["title"], "直播") {

				return true
			}

			return false

		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal, context.Background())

	s.Start()

}

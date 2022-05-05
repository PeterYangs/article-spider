package main

import (
	"context"
	"encoding/json"
	articleSpider "github.com/PeterYangs/article-spider/v3"
	"time"
)

func main() {

	cxt, _ := context.WithTimeout(context.Background(), 10000*time.Second)

	f := articleSpider.Form{
		Host:      "http://www.tiyuxiu.com",
		Channel:   "/data/list_0_[PAGE].json?__t=16339338",
		PageStart: 1,
		Length:    10,
		DetailFields: map[string]articleSpider.Field{

			"title":   {Types: articleSpider.Text, Selector: "h1"},
			"content": {Types: articleSpider.HtmlWithImage, Selector: "#main-content"},
		},
		ApiConversion: func(html string, form *articleSpider.Form) []string {

			type list struct {
				Url string
			}

			var l []list

			json.Unmarshal([]byte(html), &l)

			var temp []string

			for _, l2 := range l {

				temp = append(temp, l2.Url)

			}

			return temp

		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Api, cxt)

	s.Start()
}

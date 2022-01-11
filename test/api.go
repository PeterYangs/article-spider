package main

import (
	"encoding/json"
	"github.com/PeterYangs/article-spider/v3/fileTypes"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/spider"
)

func main() {

	s := spider.NewSpider()

	s.LoadForm(form.CustomForm{
		Host:      "http://www.tiyuxiu.com",
		Channel:   "/data/list_0_[PAGE].json?__t=16339338",
		PageStart: 1,
		Length:    10,
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.Text, Selector: "h1"},
		},
		//CustomExcelHeader:     true,
		DetailCoroutineNumber: 2,
		ApiConversion: func(html string, form2 *form.Form) []string {

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
	})

	s.StartApi()

}

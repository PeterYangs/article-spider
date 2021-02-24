package spider

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
)

//爬取详情
func GetDetail(form form.Form, detailUrl string, wait *sync.WaitGroup) {

	defer wait.Done()

	html, err := tools.GetWithString(detailUrl)

	if err != nil {

		fmt.Println(err)

		return

	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		//log.Fatal(err)

		fmt.Println(err)

		return

	}

	for field, item := range form.DetailFields {

		switch item.Types {

		case fileTypes.SingleField:

			v := doc.Find(item.SingleSelector).Text()

			fmt.Println(field + "----" + v)

		}

	}

}

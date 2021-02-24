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

	//获取详情页面html
	html, err := tools.GetWithString(detailUrl)

	if err != nil {

		fmt.Println(err)

		return

	}

	//加载
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		//log.Fatal(err)

		fmt.Println(err)

		return

	}

	var res = make(map[string]string)

	//解析详情页面选择器
	for field, item := range form.DetailFields {

		switch item.Types {

		case fileTypes.SingleField:

			v := doc.Find(item.SingleSelector).Text()

			fmt.Println(v)

			res[field] = v

		}

	}

	//写入管道
	form.Storage <- res

}

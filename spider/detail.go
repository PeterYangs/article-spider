package spider

import (
	"article-spider/common"
	"article-spider/form"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
)

//爬取详情
func GetDetail(form form.Form, detailUrl string, wait *sync.WaitGroup, detailMaxChan chan int) {

	//form.s

	defer func(detailMaxChan chan int, max int) {

		if max != 0 {

			<-detailMaxChan

		}

		wait.Done()

	}(detailMaxChan, form.DetailMaxCoroutine)

	//获取详情页面html
	html, err := tools.GetToString(detailUrl, tools.HttpSetting{})

	if err != nil {

		fmt.Println(err)

		return

	}

	//自动转码
	if form.DisableAutoCoding == false {

		html, err = common.DealCoding(html)

		if err != nil {

			fmt.Println(err)

			return

		}

	}

	//panic(html)

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

	//解析选择器返回map
	res := common.ResolveSelector(form, doc, form.DetailFields)

	//合并列表中数据
	for i, v := range form.StorageTemp {

		res[i] = v

	}

	//写入管道
	form.Storage <- res

}

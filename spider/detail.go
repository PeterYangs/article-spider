package spider

import (
	"github.com/PeterYangs/article-spider/v2/common"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
)

// GetDetail 爬取详情
func GetDetail(form form.Form, detailUrl string, wait *sync.WaitGroup, detailMaxChan chan int) {

	defer func(detailMaxChan chan int, max int) {

		if max != 0 {

			<-detailMaxChan

		}

		wait.Done()

	}(detailMaxChan, form.DetailMaxCoroutine)

	//获取详情页面html
	html, header, err := form.Client.Request().GetToStringWithHeader(detailUrl)

	if err != nil {

		common.ErrorLine(form, err.Error())

		return

	}

	//自动转码
	if form.DisableAutoCoding == false {

		html, err = common.DealCoding(html, header)

		if err != nil {

			common.ErrorLine(form, err.Error())

			return

		}

	}

	//加载
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		//log.Fatal(err)

		common.ErrorLine(form, err.Error())

		return

	}

	//解析选择器返回map
	res := common.ResolveSelector(form, doc, form.DetailFields)

	//t:=doc.Find(form.NextSelector).Text()

	//f:=form.DetailFields["title"].

	//fmt.Println(doc.Find(form.DetailFields.NextSelector).Text())
	//fmt.Println(form.NextSelector)

	//合并列表中数据
	for i, v := range form.StorageTemp {

		res[i] = v

	}

	//处理格式转换
	res = common.ConversionFormat(form, res)

	//写入管道
	form.Storage <- res

}

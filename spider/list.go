package spider

import (
	"fmt"
	"github.com/PeterYangs/article-spider/common"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
)

func GetList(form form.Form) {

	hostLast := tools.SubStr(form.Host, len(form.Host)-1, 1)

	if hostLast == "/" {

		form.Host = tools.SubStr(form.Host, 0, len(form.Host)-1)
	}

	ChannelFirst := tools.SubStr(form.Channel, 0, 1)

	if ChannelFirst != "/" {

		form.Channel = "/" + form.Channel
	}

	common.GetChannelList(form, func(listUrl string) {

		html, header, err := form.Client.Request().GetToStringWithHeader(listUrl)

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

		//goquery加载html
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {

			common.ErrorLine(form, err.Error())

			return

		}

		//详情页面并发同步锁
		var wait sync.WaitGroup

		//控制详情页面最大并发数管道
		detailMaxChan := make(chan int, form.DetailMaxCoroutine)

		findOne := false

		//查找列表中的a链接
		doc.Find(form.ListSelector).Each(func(i int, s *goquery.Selection) {

			findOne = true

			//只爬列表
			isReturn := common.OnlyList(form, s)

			if isReturn {

				return
			}

			href := ""

			isFind := false

			//a链接是列表的情况
			if form.ListHrefSelector == "" {

				href, isFind = s.Attr("href")

			} else {

				href, isFind = s.Find(form.ListHrefSelector).Attr("href")

			}

			if href == "" {

				fmt.Println("a链接为空")

				return
			}

			if isFind == true {

				href = common.GetHref(href, form.Host)

				//列表选择器不为空时
				if len(form.ListFields) > 0 {

					t, err := s.Html()

					if err != nil {

						common.ErrorLine(form, err.Error())

						return

					}

					tempDoc, err := goquery.NewDocumentFromReader(strings.NewReader(t))

					res := common.ResolveSelector(form, tempDoc, form.ListFields)

					if len(res) != 0 {

						form.StorageTemp = res
					}

				}

				wait.Add(1)

				//控制最大并发
				if form.DetailMaxCoroutine != 0 {

					detailMaxChan <- 1

				}

				//根据列表的长度开启协程爬取详情页
				go GetDetail(form, href, &wait, detailMaxChan)

			}

		})

		if !findOne {

			common.ErrorLine(form, "未找到任何a链接")
		}

		wait.Wait()

		//}

	})

	//通知excel已完成
	form.IsFinish <- true

}

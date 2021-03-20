package spider

import (
	"article-spider/common"
	"article-spider/form"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"sync"
)

func GetList(form form.Form) {

	//当前页码
	var pageCurrent int

	for pageCurrent = form.PageStart; pageCurrent <= form.Limit; pageCurrent++ {

		//当前列表url
		listUrl := form.Host + strings.Replace(form.Channel, "[PAGE]", strconv.Itoa(pageCurrent), -1)

		//获取html页面
		html, err := tools.GetToString(listUrl, form.HttpSetting)

		if err != nil {

			fmt.Println(err)

			continue

		}

		html, err = common.DealCoding(html)

		if err != nil {

			fmt.Println(err)

			continue

		}

		//自动转码
		if form.DisableAutoCoding == false {

			html, err = common.DealCoding(html)

			if err != nil {

				fmt.Println(err)

				continue

			}

		}

		//if err != nil {
		//
		//	fmt.Println(err)
		//
		//	continue
		//
		//}

		//goquery加载html
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {

			fmt.Println(err)

			continue

		}

		//详情页面并发同步锁
		var wait sync.WaitGroup

		//控制详情页面最大并发数管道
		detailMaxChan := make(chan int, form.DetailMaxCoroutine)

		//查找列表中的a链接
		doc.Find(form.ListSelector).Each(func(i int, s *goquery.Selection) {

			//只爬列表
			if len(form.DetailFields) <= 0 && len(form.ListFields) > 0 {

				ts, err := s.Html()

				if err != nil {

					log.Println(err)

					return

				}

				tempDoc, err := goquery.NewDocumentFromReader(strings.NewReader(ts))

				if err != nil {

					fmt.Println(err)

					return
				}

				res := common.ResolveSelector(form, tempDoc, form.ListFields)

				form.Storage <- res

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

				//doc, err := goquery.NewDocumentFromReader(s)

				//列表选择器不为空时
				if len(form.ListFields) > 0 {

					t, err := s.Html()

					if err != nil {

						log.Println(err)

						return

					}

					tempDoc, err := goquery.NewDocumentFromReader(strings.NewReader(t))

					res := common.ResolveSelector(form, tempDoc, form.ListFields)

					if len(res) != 0 {

						form.StorageTemp = res
					}

				}

				//if len(form.DetailFields) > 0 {

				wait.Add(1)

				//控制最大并发
				if form.DetailMaxCoroutine != 0 {

					detailMaxChan <- 1

				}

				//根据列表的长度开启协程爬取详情页
				go GetDetail(form, href, &wait, detailMaxChan)

				//} else {
				//
				//	//panic("")
				//
				//	//只爬列表
				//
				//	form.Storage <- res
				//
				//}

			}

		})

		wait.Wait()

	}

	//b:=true
	//
	//b_:=&b
	//
	//form.IsFinish=b_
	//fmt.Println("gg")
	//
	//fmt.Println(len(form.IsFinish))

	form.IsFinish <- true

	//fmt.Println("gg2")

}

package spider

import (
	"article-spider/common"
	"article-spider/form"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
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
		html, err := tools.GetWithString(listUrl)

		if err != nil {

			fmt.Println(err)

			continue

		}

		//goquery加载html
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {

			fmt.Println(err)

			continue

		}

		//详情页面并发同步锁
		var wait sync.WaitGroup

		//查找列表中的a链接
		doc.Find(form.ListSelector).Each(func(i int, s *goquery.Selection) {

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

				wait.Add(1)

				//根据列表的长度开启协程爬取详情页
				go GetDetail(form, href, &wait)

			}

		})

		wait.Wait()

	}

	println("执行完毕")

}

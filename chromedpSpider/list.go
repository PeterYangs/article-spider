package chromedpSpider

import (
	"article-spider/chromedpForm"
	"article-spider/common"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

func GetList(form chromedpForm.Form) {

	//浏览器设置
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), //关闭无头
	)

	allocCtx, AllocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	defer AllocatorCancel()

	//创建一个浏览器实例
	ctx, chromedpCancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer chromedpCancel()

	// 设置超时时间
	ctx, TimeoutCancel := context.WithTimeout(ctx, 60*time.Second)

	defer TimeoutCancel()

	//当前页码
	var pageCurrent int

	for pageCurrent = 0; pageCurrent <= form.Limit; pageCurrent++ {

		//html, err := tools.GetToString(listUrl, form.HttpSetting)

		html := ""
		err := chromedp.Run(ctx,
			chromedp.Navigate(form.Host+form.Channel),
			chromedp.WaitVisible(form.WaitForListSelector),
			chromedp.OuterHTML("html", &html, chromedp.ByQuery),
		)

		if err != nil {

			fmt.Println(err)

			return

		}

		fmt.Println(html)

		//goquery加载html
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {

			//common.ErrorLine(form, err.Error())

			fmt.Println(err)

			continue

		}

		return

	}

}

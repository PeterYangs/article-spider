package chromedpSpider

import (
	"article-spider/chromedpForm"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
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

		var list []*cdp.Node

		err := chromedp.Run(ctx,
			chromedp.Navigate(form.Host+form.Channel),
			chromedp.WaitVisible(form.WaitForListSelector),
			chromedp.Nodes(form.ListSelector, &list, chromedp.ByQueryAll),
		)

		if err != nil {

			fmt.Println(err)

			return

		}

		//fmt.Println(list)

		for _, v := range list {

			//_=v

			//fmt.Println(v.FullXPath())

			///html/body/div[10]/div[3]/ul/li[1]/div[2]/a

			txt := ""

			err := chromedp.Run(ctx,

				//chromedp.Navigate("https://www.baidu.com"),

				chromedp.WaitVisible(v.FullXPath()+"/div[2]/a"),
				chromedp.Click(v.FullXPath()+"/div[2]/a"),
				chromedp.WaitVisible("/html/body/div[10]/div[2]/div[1]/div[1]/div[2]/div[1]"),
				chromedp.Text("/html/body/div[10]/div[2]/div[1]/div[1]/div[2]/div[1]", &txt),
			)

			//err=chromedp.Run(ctx,
			//	//chromedp.Navigate(form.Host+form.Channel),
			//	//chromedp.WaitVisible(form.WaitForListSelector),
			//	chromedp.Nodes(v.,&list,chromedp.ByQuery),
			//)
			//
			if err != nil {

				fmt.Println(err)

				return

			}
			//
			//

			fmt.Println(txt)

			time.Sleep(5 * time.Second)

			return

		}

		//err := chromedp.Run(ctx,
		//	chromedp.Navigate(form.Host+form.Channel),
		//	chromedp.WaitVisible(form.WaitForListSelector),
		//	chromedp.OuterHTML("html", &html, chromedp.ByQuery),
		//)

		////fmt.Println(html)
		//
		////goquery加载html
		//doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		//if err != nil {
		//
		//	//common.ErrorLine(form, err.Error())
		//
		//	fmt.Println(err)
		//
		//	continue
		//
		//}
		//
		//
		//doc.Find(form.ListSelector).Each(func(i int, s *goquery.Selection) {
		//
		//	//href := ""
		//	//
		//	//isFind := false
		//
		//	//a链接是列表的情况
		//	if form.ListClickSelector == "" {
		//
		//		//href, isFind = s.Attr("href")
		//
		//
		//
		//	} else {
		//
		//		gg:=s.Find(form.ListClickSelector).First().Nodes
		//
		//		//fmt.Println(gg)
		//
		//		for _,v:=range gg{
		//
		//			fmt.Println(*v)
		//
		//		}
		//
		//	}
		//})

		return

	}

}

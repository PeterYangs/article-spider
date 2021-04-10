package chromedpSpider

import (
	"article-spider/common"
	"article-spider/form"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

func GetList(form form.Form) {

	ctx, TimeoutCancel := createContext(-1)

	defer TimeoutCancel()

	// 监听得到第二个tab页的target ID
	ch := make(chan target.ID, 2)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*target.EventTargetCreated); ok &&
			// if OpenerID == "", this is the first tab.
			ev.TargetInfo.OpenerID != "" {
			ch <- ev.TargetInfo.TargetID
		}
	})

	//当前页码
	var pageCurrent int

	for pageCurrent = 0; pageCurrent <= form.Limit; pageCurrent++ {

		//html, err := tools.GetToString(listUrl, form.HttpSetting)

		if pageCurrent == 0 {

			chromedp.Run(ctx,
				chromedp.Navigate(form.Host+form.Channel),
			)

		}

		var list []*cdp.Node

		err := chromedp.Run(ctx,
			chromedp.WaitVisible(form.WaitForListSelector, chromedp.ByQuery),
			chromedp.Nodes(form.ListPath, &list),
		)

		if err != nil {

			fmt.Println(err)

			return

		}

		for _, v := range list {

			err := chromedp.Run(ctx,

				//chromedp.Navigate("https://www.baidu.com"),

				chromedp.WaitVisible(v.FullXPath()+form.ListClickPath),
				chromedp.Click(v.FullXPath()+form.ListClickPath),
			)

			if err != nil {

				fmt.Println(err)

				return

			}

			waitNewTab := time.After(1 * time.Second)

			select {

			case TargetID := <-ch:

				ctx2, newTabCancle := chromedp.NewContext(ctx, chromedp.WithTargetID(TargetID))

				html := ""

				chromedp.Run(ctx2,
					chromedp.Sleep(1*time.Second),
					chromedp.OuterHTML("html", &html, chromedp.ByQuery),
				)

				//加载
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

				if err != nil {
					//log.Fatal(err)

					common.ErrorLine(form, err.Error())

					return

				}

				//解析选择器返回map
				res := common.ResolveSelector(form, doc, form.DetailFields)

				//合并列表中数据
				for i, v := range form.StorageTemp {

					res[i] = v

				}

				//处理格式转换
				res = common.ConversionFormat(form, res)

				//写入管道
				form.Storage <- res

				newTabCancle()

			case <-waitNewTab:

				fmt.Println("nothing")

			}

			time.Sleep(1 * time.Second)

		}

		//点击下一页
		chromedp.Run(ctx, chromedp.Click(form.NextSelector, chromedp.ByQuery))

		//return

	}

	//通知excel已完成
	form.IsFinish <- true

}

func createContext(timeout int) (context.Context, context.CancelFunc) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	if timeout != -1 {

		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	}

	return ctx, cancel
}

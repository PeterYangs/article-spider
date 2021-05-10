package chromedpSpider

import (
	"context"
	"fmt"
	"github.com/PeterYangs/article-spider/common"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
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

	common.GetChannelList(form, func(pageCurrent string) {

		if pageCurrent == "0" {

			chromedp.Run(ctx,
				chromedp.Navigate(form.Host+form.Channel),
			)

		}

		var html string

		err := chromedp.Run(ctx,
			chromedp.WaitVisible(form.WaitForListSelector, chromedp.ByQuery),
			chromedp.OuterHTML("html", &html, chromedp.ByQuery),
		)

		if err != nil {

			fmt.Println(err)

			return

		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

		if err != nil {

			log.Fatal(err)
		}

		doc.Find(form.ListSelector).Each(func(index int, selection *goquery.Selection) {

			isReturn := common.OnlyList(form, selection)

			if isReturn {

				return
			}

			isFind := false

			//a链接是列表的情况
			if form.ListHrefSelector == "" {

				isFind = true

			} else {

				size := selection.Find(form.ListHrefSelector).Size()

				if size >= 1 {

					isFind = true

				}

			}

			if isFind {

				//需要点击的选择器
				currentSelector := form.ListSelector + ":nth-child(" + strconv.Itoa(index+1) + ")" + " > " + form.ListHrefSelector

				chromedp.Run(
					ctx,
					chromedp.WaitVisible(currentSelector, chromedp.ByQuery),
					chromedp.Click(currentSelector, chromedp.ByQuery),
				)

				//列表选择器不为空时(放在这里是因为有些网站的图片是懒加载，点击后再获取列表信息)
				if len(form.ListFields) > 0 {

					t, err := selection.Html()

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

					//fmt.Println("nothing")

					html := ""

					chromedp.Run(ctx,
						chromedp.Sleep(1*time.Second),
						chromedp.OuterHTML("html", &html, chromedp.ByQuery),
					)

					//加载
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

					if err != nil {

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

					//返回上一页
					chromedp.Run(ctx,
						chromedp.NavigateBack(),
						chromedp.Sleep(1*time.Second),
					)

				}

			}

		})

		//点击下一页
		chromedp.Run(ctx, chromedp.Click(form.NextSelector, chromedp.ByQuery))

		//return

	})

	//通知excel已完成
	form.IsFinish <- true

}

func createContext(timeout int) (context.Context, context.CancelFunc) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	if timeout != -1 {

		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	}

	return ctx, cancel
}

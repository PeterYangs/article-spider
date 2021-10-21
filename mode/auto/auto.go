package auto

import (
	"context"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/notice"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"strings"
)

type auto struct {
	form *form.Form
}

func NewAuto(form *form.Form) *auto {

	return &auto{form: form}
}

func (a *auto) GetList() {

	cxt, _ := a.createContext()

	var html string

	chromedp.Run(
		cxt,
		chromedp.Navigate(a.form.Host+a.form.Channel),
		chromedp.WaitVisible("body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li:nth-child(1)", chromedp.ByQuery),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	//panic(html)

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {

		a.form.Notice.PushMessage(notice.NewError(err.Error()))

		return

	}

	//查找列表中的a链接
	size := doc.Find(a.form.ListSelector).Each(func(i int, s *goquery.Selection) {

		href := ""

		isFind := false

		storage := make(map[string]string)

		//a链接是列表的情况
		if n.form.HrefSelector == "" {

			href, isFind = s.Attr("href")

		} else {

			href, isFind = s.Find(n.form.HrefSelector).Attr("href")

		}

		if href == "" || isFind == false {

			n.form.Notice.PushMessage(notice.NewError("a链接为空"))

			return
		}

		//列表选择器不为空时
		if len(n.form.ListFields) > 0 {

			t, err := s.Html()

			if err != nil {

				n.form.Notice.PushMessage(notice.NewError(err.Error()))

				return

			}

			//解析列表选择器
			storage, err = n.form.ResolveSelector(t, n.form.ListFields, listUrl)

			if err != nil {

				n.form.Notice.PushMessage(notice.NewError(err.Error()))

				return
			}

		}

		//如果详情选择器为空就跳过
		if len(n.form.DetailFields) <= 0 {

			n.form.Storage <- storage

			return

		}

		//控制协程并发数
		n.form.DetailCoroutineChan <- true

		n.form.DetailWait.Add(1)

		go n.GetDetail(n.form.GetHref(href), storage)

	}).Size()

}

func (a *auto) createContext() (context.Context, context.CancelFunc) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		// chromedp.WithDebugf(log.Printf),
	)
	//defer cancel()

	return ctx, cancel

}

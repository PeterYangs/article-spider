package auto

import (
	"context"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/notice"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"strconv"
	"strings"
	"time"
)

type auto struct {
	form *form.Form
}

func NewAuto(form *form.Form) *auto {

	return &auto{form: form}
}

func (a *auto) GetList() {

	cxt, cancel := a.createContext()

	// 监听得到第二个tab页的target ID
	ch := make(chan target.ID, 2)
	chromedp.ListenTarget(cxt, func(ev interface{}) {
		if ev, ok := ev.(*target.EventTargetCreated); ok &&
			// if OpenerID == "", this is the first tab.
			ev.TargetInfo.OpenerID != "" {
			ch <- ev.TargetInfo.TargetID
		}
	})

	if a.form.HttpHeader["cookie"] != "" {

		chromedp.Run(
			cxt,
			a.setcookies(a.getCookieMap(a.form.HttpHeader["cookie"])),
			chromedp.Navigate(a.form.Host+a.form.Channel),
		)

		//panic("xx")
	} else {

		chromedp.Run(
			cxt,
			//a.setcookies(a.getCookieMap(a.form.HttpHeader["cookie"])),
			chromedp.Navigate(a.form.Host+a.form.Channel),
		)
	}

	//chromedp.Run(
	//	cxt,
	//	//a.setcookies(a.getCookieMap(a.form.AutoCookieString)),
	//	chromedp.Navigate(a.form.Host+a.form.Channel),
	//)

	a.dealList(cxt, cancel, ch)

}

func (a *auto) dealList(cxt context.Context, cancel context.CancelFunc, ch chan target.ID) {

	var html string

	listCxt, _ := context.WithTimeout(cxt, 3*time.Second)

	chromedp.Run(
		listCxt,
		chromedp.WaitVisible(a.form.ListWaitSelector, chromedp.ByQuery),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {

		a.form.Notice.PushMessage(notice.NewError(err.Error()))

		return

	}

	//查找列表中的a链接
	size := doc.Find(a.form.ListSelector).Each(func(i int, s *goquery.Selection) {

		storage := make(map[string]string)

		//列表选择器不为空时
		if len(a.form.ListFields) > 0 {

			t, err := s.Html()

			if err != nil {

				a.form.Notice.PushMessage(notice.NewError(err.Error()))

				return

			}

			//解析列表选择器
			storage, err = a.form.ResolveSelector(t, a.form.ListFields, a.form.Host)

			if err != nil {

				a.form.Notice.PushMessage(notice.NewError(err.Error()))

				return
			}

		}

		//如果详情选择器为空就跳过
		if len(a.form.DetailFields) <= 0 {

			a.form.Storage <- storage

			//相当于详情完成一个
			a.form.CurrentIndex++

			return

		}

		clickSelector := a.form.ListSelector + ":nth-child(" + strconv.Itoa(i+1) + ")" + " > " + a.form.HrefSelector

		clickLength := doc.Find(clickSelector).Size()

		if clickLength <= 0 {

			a.form.Notice.PushMessage(notice.NewError("未找到详情选择器"))

			return
		}

		//点击详情选择器
		e := chromedp.Run(
			cxt,
			chromedp.WaitVisible(clickSelector, chromedp.ByQuery),
			chromedp.Click(clickSelector, chromedp.NodeVisible),
		)

		if e != nil {

			a.form.Notice.PushMessage(notice.NewError(e.Error()))

		}

		waitNewTab := time.After(2 * time.Second)

		select {

		case TargetID := <-ch:

			detailCxt, NewTabCancel := chromedp.NewContext(cxt, chromedp.WithTargetID(TargetID))

			a.GetDetail(detailCxt, storage, true, NewTabCancel)

		case <-waitNewTab:

			a.GetDetail(cxt, storage, false, cancel)

		}

	}).Size()

	a.form.PageCurrent++

	if a.form.DetailSize == 0 && size > 0 {

		a.form.DetailSize = size

		//计算大概爬取总数量
		a.form.Total = a.form.Length * size

	}

	if size <= 0 {

		a.form.Notice.PushMessage(notice.NewInfo("a链接未发现"))
	}

	cxt, cancel = a.clickNext(cxt, cancel, ch)

	if a.form.PageCurrent >= a.form.Length {

		a.form.Notice.PushMessage(notice.NewError("完成"))

		return

	} else {

		a.dealList(cxt, cancel, ch)

	}

}

func (a *auto) GetDetail(detailCxt context.Context, storage map[string]string, isNewTab bool, cancel context.CancelFunc) {

	defer func() {

		a.form.CurrentIndex++

	}()

	html := ""

	tempCxt, _ := context.WithTimeout(detailCxt, 5*time.Second)

	e := chromedp.Run(tempCxt,
		chromedp.WaitVisible(a.form.DetailWaitSelector, chromedp.ByQuery),
		//chromedp.Wait
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	//fmt.Println(html)

	if e != nil {

		a.form.Notice.PushMessage(notice.NewError(e.Error()))
	}

	res, err := a.form.ResolveSelector(html, a.form.DetailFields, a.form.Host)

	if err != nil {

		a.form.Notice.PushMessage(notice.NewError(err.Error()))

		return
	}

	//合并列表结果
	for s, s2 := range res {

		storage[s] = s2

	}

	for s, s2 := range storage {

		storage[s] = strings.TrimSpace(s2)
	}

	a.form.Storage <- storage

	if isNewTab {

		cancel()

	} else {

		backCxt, _ := context.WithTimeout(detailCxt, 3*time.Second)

		//返回上一页
		chromedp.Run(backCxt,
			chromedp.NavigateBack(),
			chromedp.Sleep(1*time.Second),
		)

	}

}

//创建chromedp的context
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

	return ctx, cancel

}

func (a *auto) clickNext(cxt context.Context, cancel context.CancelFunc, ch chan target.ID) (context.Context, context.CancelFunc) {

	clickCxt, _ := context.WithTimeout(cxt, 3*time.Second)

	//点击下一页
	chromedp.Run(clickCxt, chromedp.Click(a.form.NextSelector, chromedp.ByQuery))

	waitNewTab := time.After(1 * time.Second)

	select {

	case TargetID := <-ch:

		detailCxt, NewTabCancel := chromedp.NewContext(cxt, chromedp.WithTargetID(TargetID))

		cancel()

		return detailCxt, NewTabCancel

	case <-waitNewTab:

		return cxt, cancel

	}

}

func (a *auto) setcookies(cookies map[string]string) chromedp.Tasks {
	//if len(cookies)%2 != 0 {
	//	panic("length of cookies must be divisible by 2")
	//}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			//for i := 0; i < len(cookies); i += 2 {
			//	err := network.SetCookie(cookies[i], cookies[i+1]).
			//		WithExpires(&expr).
			//		WithDomain("localhost").
			//		WithHTTPOnly(true).
			//		Do(ctx)
			//	if err != nil {
			//		return err
			//	}
			//
			//	//network.SetCookie()
			//
			//
			//}

			for s, s2 := range cookies {

				err := network.SetCookie(s, s2).
					WithExpires(&expr).
					WithDomain("www.925g.com").
					WithHTTPOnly(true).
					Do(ctx)

				if err != nil {
					return err
				}

			}

			return nil
		}),
		// navigate to site
		//chromedp.Navigate(host),
		//// read the returned values
		//chromedp.Text(`#result`, res, chromedp.ByID, chromedp.NodeVisible),
		// read network values
		//chromedp.ActionFunc(func(ctx context.Context) error {
		//	cookies, err := network.GetAllCookies().Do(ctx)
		//	if err != nil {
		//		return err
		//	}
		//
		//	for i, cookie := range cookies {
		//		log.Printf("chrome cookie %d: %+v", i, cookie)
		//	}
		//
		//	return nil
		//}),
	}
}

func (a *auto) getCookieMap(cookie string) map[string]string {

	cookieMap := make(map[string]string)

	arr := tools.Explode("; ", cookie)

	for _, s := range arr {

		index := strings.Index(s, "=")

		cookieMap[s[:index]] = s[index+1:]

	}

	return cookieMap

}

package article_spider

import (
	"context"
	"errors"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type auto struct {
	s         *Spider
	listIndex int //当前列表进行个数
	page      int //当前页码
	size      int //第一页列表长度
}

func NewAuto(s *Spider) *auto {

	return &auto{s: s}
}

func (a *auto) Start() error {

	if a.s.form.AutoDetailWaitSelector == "" {

		return errors.New("详情等待选择器未配置")
	}

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

	if a.s.form.HttpHeader["cookie"] != "" {

		task, err := a.setcookies(a.getCookieMap(a.s.form.HttpHeader["cookie"]))

		if err != nil {

			a.s.notice.Error(err.Error())

			return nil
		}

		chromedp.Run(
			cxt,
			task,
			chromedp.Navigate(a.s.form.Host+a.s.form.Channel),
		)

	} else {

		chromedp.Run(
			cxt,
			chromedp.Navigate(a.s.form.Host+a.s.form.Channel),
		)
	}

	//执行前置事件
	if a.s.form.AutoPrefixEvent != nil {

		a.s.form.AutoPrefixEvent(cxt)

	}

	a.GetList(cxt, cancel, ch)

	return nil

}

func (a *auto) GetList(cxt context.Context, cancel context.CancelFunc, ch chan target.ID) {

	for {

		var ListErr error

		cxt, cancel, ListErr = a.dealList(cxt, cancel, ch)

		if ListErr != nil {

			a.s.notice.Error(ListErr.Error())
		}

		if a.page >= a.s.form.Length {

			a.s.cancel()

			return

		}

	}

}

func (a *auto) dealList(cxt context.Context, cancel context.CancelFunc, ch chan target.ID) (context.Context, context.CancelFunc, error) {

	defer func() {

		a.page++

		a.s.autoPage = a.page

	}()

	if a.s.form.AutoListWaitSelector != "" {

		listCxt, _ := context.WithTimeout(cxt, 6*time.Second)

		//等待选择器
		chromedp.Run(
			listCxt,
			chromedp.WaitVisible(a.s.form.AutoListWaitSelector, chromedp.ByQuery),
		)

	} else {

		listCxt, _ := context.WithTimeout(cxt, 6*time.Second)

		//等待选择器
		chromedp.Run(
			listCxt,
			chromedp.WaitVisible(a.s.form.ListSelector+":nth-child(1)", chromedp.ByQuery),
		)

	}

	var html string

	htmlCxt, _ := context.WithTimeout(cxt, 6*time.Second)

	e := chromedp.Run(htmlCxt, chromedp.OuterHTML("html", &html, chromedp.ByQuery))

	if e != nil {

		return cxt, cancel, errors.New(fmt.Sprint("获取html失败", e))
	}

	if html == "" {

		return cxt, cancel, errors.New(fmt.Sprint("html为空"))
	}

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {

		return cxt, cancel, err

	}

	size := doc.Find(a.s.form.ListSelector).Size()

	if a.size == 0 {

		a.size = size
	}

	//查找列表中的a链接
	doc.Find(a.s.form.ListSelector).Each(func(i int, s *goquery.Selection) {

		//不要把上一页的长度也计算进去
		if a.s.form.AutoNextPageMode == LoadMore {

			if i+1 > a.size {

				return
			}
		}

		storage := make(map[string]string)

		//列表选择器不为空时
		if len(a.s.form.ListFields) > 0 {

			t, err := s.Html()

			if err != nil {

				a.s.notice.Error(err.Error())

				return

			}

			//解析列表选择器
			storage, err = a.s.form.ResolveSelector(t, a.s.form.ListFields, a.s.form.Host)

			if err != nil {

				a.s.notice.Error(err.Error())

				return
			}

		}

		//如果详情选择器为空就跳过
		if len(a.s.form.DetailFields) <= 0 {

			a.s.result.Push(storage)

			//相当于详情完成一个
			a.s.currentIndex++

			return

		}

		//详情点击选择器
		var clickSelector string

		switch a.s.form.AutoNextPageMode {

		case Pagination:

			//获取列表中的a链接选择器
			clickSelector = a.s.form.ListSelector + ":nth-child(" + strconv.Itoa(i+1) + ")" + " > " + a.s.form.HrefSelector

			break

		case LoadMore:

			clickSelector = a.s.form.ListSelector + ":nth-child(" + strconv.Itoa(a.listIndex+1) + ")" + " > " + a.s.form.HrefSelector

			a.listIndex++

		default:

			//获取列表中的a链接选择器
			clickSelector = a.s.form.ListSelector + ":nth-child(" + strconv.Itoa(i+1) + ")" + " > " + a.s.form.HrefSelector

		}

		cxtW, _ := context.WithTimeout(cxt, 6*time.Second)

		//点击详情选择器
		e := chromedp.Run(
			cxtW,
			chromedp.WaitVisible(clickSelector, chromedp.ByQuery),
			//chromedp.Click(clickSelector, chromedp.NodeVisible),
		)

		if e != nil {

			a.s.notice.Error("未找到详情选择器:", clickSelector, e)

			if a.s.form.AutoNextPageMode == LoadMore {

				//点击下一页
				cxt, cancel = a.clickNext(cxt, cancel, ch)

			}

			return
		}

		////等待详情选择器
		//chromedp.Run(
		//	cxt,
		//	chromedp.WaitVisible(clickSelector, chromedp.ByQuery),
		//)

		href := ""

		isFind := false

		if a.s.form.AutoNextPageMode == LoadMore {

			href, isFind = doc.Find(clickSelector).Attr("href")

			if !isFind {

				return
			}

		} else {

			//a链接是列表的情况
			if a.s.form.HrefSelector == "" {

				href, isFind = s.Attr("href")

			} else {

				href, isFind = s.Find(a.s.form.HrefSelector).Attr("href")

			}

		}

		//fmt.Println(a.s.form.AutoDetailForceNewTab, isFind, href)

		//点击详情页强制打开新窗口
		if a.s.form.AutoDetailForceNewTab && isFind && href != "" {

			tag := target.CreateTarget("")

			detailCxt, NewTabCancel := chromedp.NewContext(cxt, chromedp.WithTargetID(target.ID(tag.BrowserContextID)))

			e = chromedp.Run(detailCxt, chromedp.Navigate(a.s.form.GetHref(href)))

			a.GetDetail(detailCxt, storage, true, NewTabCancel)

		} else {

			//点击详情选择器
			e = chromedp.Run(
				cxt,
				chromedp.WaitVisible(clickSelector, chromedp.ByQuery),
				chromedp.Click(clickSelector, chromedp.NodeVisible),
			)

			waitNewTab := time.After(6 * time.Second)

			select {

			case TargetID := <-ch:

				detailCxt, NewTabCancel := chromedp.NewContext(cxt, chromedp.WithTargetID(TargetID))

				a.GetDetail(detailCxt, storage, true, NewTabCancel)

			case <-waitNewTab:

				a.GetDetail(cxt, storage, false, cancel)

			}

		}

		if e != nil {

			a.s.notice.Error("点击详情选择器错误", e.Error())

			return

		}

	}).Size()

	//点击下一页
	cxt, cancel = a.clickNext(cxt, cancel, ch)

	return cxt, cancel, nil

}

func (a *auto) GetDetail(detailCxt context.Context, storage map[string]string, isNewTab bool, cancel context.CancelFunc) {

	defer func() {

		a.s.currentIndex++

		if isNewTab {

			cancel()

		} else {

			backCxt, _ := context.WithTimeout(detailCxt, 6*time.Second)

			//返回上一页
			chromedp.Run(backCxt,
				chromedp.NavigateBack(),
				chromedp.Sleep(1*time.Second),
			)

		}

	}()

	if a.s.form.AutoDetailWaitSelector != "" {

		detailWaitCxt, _ := context.WithTimeout(detailCxt, 6*time.Second)

		e := chromedp.Run(detailWaitCxt,
			chromedp.WaitVisible(a.s.form.AutoDetailWaitSelector, chromedp.ByQuery),
			//chromedp.Wait
			//chromedp.OuterHTML("html", &html, chromedp.ByQuery),
		)

		if e != nil {

			a.s.notice.Error(e.Error())
		}

	}

	html := ""

	htmlCxt, _ := context.WithTimeout(detailCxt, 6*time.Second)

	e := chromedp.Run(htmlCxt, chromedp.OuterHTML("html", &html, chromedp.ByQuery))

	if e != nil {

		a.s.notice.Error(fmt.Sprint("获取html失败", e))

		//return cxt, cancel, errors.New()

		return
	}

	if html == "" {

		a.s.notice.Error(fmt.Sprint("html为空"))

		//return cxt, cancel, errors.New()

		return
	}

	res, err := a.s.form.ResolveSelector(html, a.s.form.DetailFields, a.s.form.Host)

	if err != nil {

		//a.form.Notice.PushMessage(notice.NewError(err.Error()))

		a.s.notice.Error(err.Error())

		return
	}

	//合并列表结果
	for s, s2 := range res {

		storage[s] = s2

	}

	for s, s2 := range storage {

		storage[s] = strings.TrimSpace(s2)
	}

	//a.form.Storage <- storage

	a.s.result.Push(storage)

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

	//设置代理
	if a.s.form.HttpProxy != "" {

		opts = append(opts,
			chromedp.ProxyServer(a.s.form.HttpProxy),
		)

	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		// chromedp.WithDebugf(log.Printf),
	)

	return ctx, cancel

}

//点击下一页
func (a *auto) clickNext(cxt context.Context, cancel context.CancelFunc, ch chan target.ID) (context.Context, context.CancelFunc) {

	clickCxt, _ := context.WithTimeout(cxt, 6*time.Second)

	//点击下一页
	chromedp.Run(clickCxt, chromedp.Click(a.s.form.AutoNextSelector, chromedp.ByQuery))

	waitNewTab := time.After(3 * time.Second)

	select {

	case TargetID := <-ch:

		detailCxt, NewTabCancel := chromedp.NewContext(cxt, chromedp.WithTargetID(TargetID))

		//新开标签关闭上一页标签
		cancel()

		return detailCxt, NewTabCancel

	case <-waitNewTab:

		return cxt, cancel

	}

}

func (a *auto) setcookies(cookies map[string]string) (chromedp.Tasks, error) {

	re1 := regexp.MustCompile("^(http|https)://([^/]+).*$").FindStringSubmatch(a.s.form.Host)

	if len(re1) <= 0 {

		return nil, errors.New("或者domain失败，请检查host设置")
	}

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
					WithDomain(re1[2]).
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
	}, nil
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

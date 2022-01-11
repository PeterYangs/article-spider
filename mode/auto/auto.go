package auto

import (
	"context"
	"errors"
	"fmt"
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PeterYangs/article-spider/v3/mode"
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
	form      *form.Form
	listIndex int //当前列表进行个数
	page      int //当前页码
	size      int //第一页列表长度
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

		task, err := a.setcookies(a.getCookieMap(a.form.HttpHeader["cookie"]))

		if err != nil {

			a.form.Notice.Error(err.Error())

			return
		}

		chromedp.Run(
			cxt,
			task,
			chromedp.Navigate(a.form.Host+a.form.Channel),
		)

	} else {

		chromedp.Run(
			cxt,
			chromedp.Navigate(a.form.Host+a.form.Channel),
		)
	}

	//执行前置事件
	if a.form.AutoPrefixEvent != nil {

		a.form.AutoPrefixEvent(cxt)

	}

	for {

		var ListErr error

		cxt, cancel, ListErr = a.dealList(cxt, cancel, ch)

		if ListErr != nil {

			a.form.Notice.Error(ListErr.Error())
		}

		if a.page >= a.form.Length {

			a.form.Notice.Error("完成")

			return

		}

	}

}

func (a *auto) dealList(cxt context.Context, cancel context.CancelFunc, ch chan target.ID) (context.Context, context.CancelFunc, error) {

	defer func() {

		a.page++

		a.form.AutoPage = a.page

	}()

	var html string

	listCxt, _ := context.WithTimeout(cxt, 6*time.Second)

	//获取列表页面的html
	chromedp.Run(
		listCxt,
		chromedp.WaitVisible(a.form.ListWaitSelector, chromedp.ByQuery),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {

		return cxt, cancel, err

	}

	size := doc.Find(a.form.ListSelector).Size()

	if a.size == 0 {

		a.size = size
	}

	//查找列表中的a链接
	doc.Find(a.form.ListSelector).Each(func(i int, s *goquery.Selection) {

		//不要把上一页的长度也计算进去
		if a.form.NextPageMode == mode.LoadMore {

			if i+1 > a.size {

				return
			}
		}

		storage := make(map[string]string)

		//列表选择器不为空时
		if len(a.form.ListFields) > 0 {

			t, err := s.Html()

			if err != nil {

				a.form.Notice.Error(err.Error())

				return

			}

			//解析列表选择器
			storage, err = a.form.ResolveSelector(t, a.form.ListFields, a.form.Host)

			if err != nil {

				a.form.Notice.Error(err.Error())

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

		//详情点击选择器
		var clickSelector string

		switch a.form.NextPageMode {

		case mode.Pagination:

			//获取列表中的a链接选择器
			clickSelector = a.form.ListSelector + ":nth-child(" + strconv.Itoa(i+1) + ")" + " > " + a.form.HrefSelector

			break

		case mode.LoadMore:

			clickSelector = a.form.ListSelector + ":nth-child(" + strconv.Itoa(a.listIndex+1) + ")" + " > " + a.form.HrefSelector

			a.listIndex++

		default:

			//获取列表中的a链接选择器
			clickSelector = a.form.ListSelector + ":nth-child(" + strconv.Itoa(i+1) + ")" + " > " + a.form.HrefSelector

		}

		cxtW, _ := context.WithTimeout(cxt, 6*time.Second)

		//点击详情选择器
		e := chromedp.Run(
			cxtW,
			chromedp.WaitVisible(clickSelector, chromedp.ByQuery),
			//chromedp.Click(clickSelector, chromedp.NodeVisible),
		)

		if e != nil {

			a.form.Notice.Error("未找到详情选择器:", clickSelector, e)

			if a.form.NextPageMode == mode.LoadMore {

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

		if a.form.NextPageMode == mode.LoadMore {

			href, isFind = doc.Find(clickSelector).Attr("href")

			if !isFind {

				return
			}

		} else {

			//a链接是列表的情况
			if a.form.HrefSelector == "" {

				href, isFind = s.Attr("href")

			} else {

				href, isFind = s.Find(a.form.HrefSelector).Attr("href")

			}

		}

		fmt.Println(a.form.AutoDetailForceNewTab, isFind, href)

		//点击详情页强制打开新窗口
		if a.form.AutoDetailForceNewTab && isFind && href != "" {

			tag := target.CreateTarget("")

			detailCxt, NewTabCancel := chromedp.NewContext(cxt, chromedp.WithTargetID(target.ID(tag.BrowserContextID)))

			e = chromedp.Run(detailCxt, chromedp.Navigate(a.form.GetHref(href)))

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

			//a.form.Notice.PushMessage(notice.NewError(e.Error()))

			a.form.Notice.Error("点击详情选择器错误", e.Error())

			return

		}

	}).Size()

	//点击下一页
	cxt, cancel = a.clickNext(cxt, cancel, ch)

	return cxt, cancel, nil

}

func (a *auto) GetDetail(detailCxt context.Context, storage map[string]string, isNewTab bool, cancel context.CancelFunc) {

	defer func() {

		a.form.CurrentIndex++

	}()

	html := ""

	tempCxt, _ := context.WithTimeout(detailCxt, 6*time.Second)

	e := chromedp.Run(tempCxt,
		chromedp.WaitVisible(a.form.DetailWaitSelector, chromedp.ByQuery),
		//chromedp.Wait
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	//fmt.Println(html)

	if e != nil {

		//a.form.Notice.PushMessage(notice.NewError(e.Error()))

		a.form.Notice.Error(e.Error())
	}

	res, err := a.form.ResolveSelector(html, a.form.DetailFields, a.form.Host)

	if err != nil {

		//a.form.Notice.PushMessage(notice.NewError(err.Error()))

		a.form.Notice.Error(err.Error())

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

		backCxt, _ := context.WithTimeout(detailCxt, 6*time.Second)

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

	//设置代理
	if a.form.HttpProxy != "" {

		opts = append(opts,
			chromedp.ProxyServer(a.form.HttpProxy),
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
	chromedp.Run(clickCxt, chromedp.Click(a.form.NextSelector, chromedp.ByQuery))

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

	re1 := regexp.MustCompile("^(http|https)://([^/]+).*$").FindStringSubmatch(a.form.Host)

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

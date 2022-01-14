package article_spider

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type api struct {
	s *Spider
}

func NewApi(s *Spider) *api {

	return &api{s: s}
}

func (a *api) Start() {

	defer func() {

		a.s.detailWait.Wait()

		a.s.cancel()

	}()

	if a.s.form.ApiConversion == nil {

		a.s.notice.Error("api转换函数未配置")

		return
	}

	//列表回调
	a.s.getChannelList(func(listUrl string) {

		a.GetList(listUrl)

	})

}

func (a *api) GetList(listUrl string) {

	html, err := a.s.form.GetHtml(listUrl)

	if err != nil {

		a.s.notice.Error(err.Error())

		return

	}

	hrefList := a.s.form.ApiConversion(html, &a.s.form)

	if len(hrefList) <= 0 {

		a.s.notice.Error("api解析链接长度为0")

		return
	}

	if a.s.detailSize == 0 && len(hrefList) > 0 {

		a.s.detailSize = len(hrefList)

		//计算大概爬取总数量
		a.s.total = a.s.form.Length * len(hrefList)

	}

	for _, s := range hrefList {

		//控制协程并发数
		a.s.detailCoroutineChan <- true

		a.s.detailWait.Add(1)

		go a.GetDetail(a.s.form.GetHref(s), map[string]string{})

	}

}

func (a *api) GetDetail(detailUrl string, storage map[string]string) {

	defer func() {

		<-a.s.detailCoroutineChan

		a.s.detailWait.Done()

		a.s.currentIndex++

	}()

	html, err := a.s.form.GetHtml(detailUrl)

	if err != nil {

		a.s.notice.Error(err.Error())

		return

	}

	//中间链接（中间页面）
	if len(a.s.form.MiddleSelector) > 0 {

		for _, s := range a.s.form.MiddleSelector {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

			if err != nil {

				a.s.notice.Error(err.Error())

				return

			}

			href, b := doc.Find(s).Attr("href")

			if !b {

				return
			}

			href = a.s.form.GetHref(href)

			html, err = a.s.form.GetHtml(href)

			if err != nil {

				a.s.notice.Error(err.Error())

				return

			}

		}

	}

	res, err := a.s.form.ResolveSelector(html, a.s.form.DetailFields, detailUrl)

	if err != nil {

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

	a.s.result.Push(storage)

}

package api

import (
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type api struct {
	form *form.Form
}

func NewApi(form *form.Form) *api {

	return &api{form: form}
}

func (a *api) GetList(listUrl string) {

	html, err := a.form.GetHtml(listUrl)

	if err != nil {

		//a.form.Notice.PushMessage(notice.NewError(err.Error()))

		a.form.Notice.Error(err.Error())

		return

	}

	if a.form.ApiConversion == nil {

		//a.form.Notice.PushMessage(notice.NewError("api转换函数未配置"))

		a.form.Notice.Error("api转换函数未配置")

		return
	}

	hrefList := a.form.ApiConversion(html, a.form)

	if len(hrefList) <= 0 {

		//a.form.Notice.PushMessage(notice.NewError("api解析链接长度为0"))

		a.form.Notice.Error("api解析链接长度为0")

		return
	}

	if a.form.DetailSize == 0 && len(hrefList) > 0 {

		a.form.DetailSize = len(hrefList)

		//计算大概爬取总数量
		a.form.Total = a.form.Length * len(hrefList)

	}

	for _, s := range hrefList {

		//控制协程并发数
		a.form.DetailCoroutineChan <- true

		a.form.DetailWait.Add(1)

		go a.GetDetail(a.form.GetHref(s), map[string]string{})

	}

}

func (a *api) GetDetail(detailUrl string, storage map[string]string) {

	defer func() {

		<-a.form.DetailCoroutineChan

		a.form.DetailWait.Done()

		a.form.CurrentIndex++

	}()

	html, err := a.form.GetHtml(detailUrl)

	if err != nil {

		//a.form.Notice.PushMessage(notice.NewError(err.Error()))

		a.form.Notice.Error(err.Error())

		return

	}

	//中间链接（中间页面）
	if len(a.form.MiddleSelector) > 0 {

		for _, s := range a.form.MiddleSelector {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

			if err != nil {

				//a.form.Notice.PushMessage(notice.NewError(err.Error()))

				a.form.Notice.Error(err.Error())

				return

			}

			href, b := doc.Find(s).Attr("href")

			if !b {

				return
			}

			href = a.form.GetHref(href)

			html, err = a.form.GetHtml(href)

			if err != nil {

				a.form.Notice.Error(err.Error())

				return

			}

		}

	}

	res, err := a.form.ResolveSelector(html, a.form.DetailFields, detailUrl)

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

}

package normal

import (
	"github.com/PeterYangs/article-spider/v3/form"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type normal struct {
	form *form.Form
}

func NewNormal(form *form.Form) *normal {

	return &normal{form: form}
}

func (n *normal) GetList(listUrl string) {

	content, header, err := n.form.Client.R().GetToContentWithHeader(listUrl)

	if err != nil {

		//n.form.Notice.PushMessage(notice.NewError(err.Error()))

		n.form.Notice.Error(err.Error())

		return

	}

	html := content.ToString()

	//自动转码
	if n.form.DisableAutoCoding == false {

		html, err = n.form.DealCoding(html, header)

		if err != nil {

			//n.form.Notice.PushMessage(notice.NewError(err.Error()))

			n.form.Notice.Error(err.Error())

			return

		}

	}

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {

		//n.form.Notice.PushMessage(notice.NewError(err.Error()))

		n.form.Notice.Error(err.Error())

		return

	}

	//查找列表中的a链接
	size := doc.Find(n.form.ListSelector).Each(func(i int, s *goquery.Selection) {

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

			//n.form.Notice.PushMessage(notice.NewError("a链接为空,当前链接为:", listUrl))

			n.form.Notice.Error("a链接为空,当前链接为:", listUrl)

			return
		}

		//列表选择器不为空时
		if len(n.form.ListFields) > 0 {

			t, err := s.Html()

			if err != nil {

				//n.form.Notice.PushMessage(notice.NewError(err.Error()))

				n.form.Notice.Error(err.Error())

				return

			}

			//解析列表选择器
			storage, err = n.form.ResolveSelector(t, n.form.ListFields, listUrl)

			if err != nil {

				//n.form.Notice.PushMessage(notice.NewError(err.Error()))

				n.form.Notice.Error(err.Error())

				return
			}

		}

		//如果详情选择器为空就跳过
		if len(n.form.DetailFields) <= 0 {

			n.form.Storage <- storage

			//相当于详情完成一个
			n.form.CurrentIndex++

			return

		}

		//控制协程并发数
		n.form.DetailCoroutineChan <- true

		n.form.DetailWait.Add(1)

		go n.GetDetail(n.form.GetHref(href), storage)

	}).Size()

	//n.form.Notice.PushMessage(notice.NewError(size))

	if n.form.DetailSize == 0 && size > 0 {

		n.form.DetailSize = size

		//计算大概爬取总数量
		n.form.Total = n.form.Length * size

	}

	if size <= 0 {

		//n.form.Notice.PushMessage(notice.NewInfo("a链接未发现"))

		n.form.Notice.Error("a链接未发现")

	}

}

func (n *normal) GetDetail(detailUrl string, storage map[string]string) {

	defer func() {

		<-n.form.DetailCoroutineChan

		n.form.DetailWait.Done()

		n.form.CurrentIndex++

	}()

	html, err := n.form.GetHtml(detailUrl)

	if err != nil {

		//n.form.Notice.PushMessage(notice.NewError(err.Error()))

		n.form.Notice.Error(err.Error())

		return

	}

	//中间链接（中间页面）
	if len(n.form.MiddleSelector) > 0 {

		for _, s := range n.form.MiddleSelector {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

			if err != nil {

				//n.form.Notice.PushMessage(notice.NewError(err.Error()))

				n.form.Notice.Error(err.Error())

				return

			}

			href, b := doc.Find(s).Attr("href")

			if !b {

				return
			}

			href = n.form.GetHref(href)

			html, err = n.form.GetHtml(href)

			if err != nil {

				//n.form.Notice.PushMessage(notice.NewError(err.Error()))

				n.form.Notice.Error(err.Error())

				return

			}

		}

	}

	res, err := n.form.ResolveSelector(html, n.form.DetailFields, detailUrl)

	if err != nil {

		//n.form.Notice.PushMessage(notice.NewError(err.Error()))

		n.form.Notice.Error(err.Error())

		return
	}

	//合并列表结果
	for s, s2 := range res {

		storage[s] = s2

	}

	for s, s2 := range storage {

		storage[s] = strings.TrimSpace(s2)
	}

	n.form.Storage <- storage

}

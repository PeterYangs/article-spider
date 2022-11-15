package article_spider

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type normal struct {
	s *Spider
}

func NewNormal(s *Spider) *normal {

	return &normal{s: s}
}

func (n normal) Start() {
	//TODO implement me

	n.s.getChannelList(func(listUrl string) {

		n.GetList(listUrl)

	})

	n.s.detailWait.Wait()

	n.s.cancel()

}

func (n normal) GetList(listUrl string) {

	select {

	case <-n.s.cxt.Done():

		return

	default:

	}

	//content, header, err := n.s.client.R().GetToContentWithHeader(listUrl)

	content, err := n.s.client.R().GetToContent(listUrl)

	if err != nil {

		n.s.notice.Error(err.Error())

		return

	}

	html := content.ToString()

	//自动转码
	if n.s.form.DisableAutoCoding == false {

		html, err = n.s.form.DealCoding(html, content.Header())

		if err != nil {

			//n.form.Notice.PushMessage(notice.NewError(err.Error()))

			n.s.notice.Error(err.Error())

			return

		}

	}

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {

		n.s.notice.Error(err.Error())

		return

	}

	//查找列表中的a链接
	size := doc.Find(n.s.form.ListSelector).Each(func(i int, s *goquery.Selection) {

		href := ""

		isFind := false

		var rows *Rows = NewRows(map[string]string{})

		//a链接是列表的情况
		if strings.TrimSpace(n.s.form.HrefSelector) == "" || n.s.form.HrefSelector == "myself" {

			href, isFind = s.Attr("href")

		} else {

			href, isFind = s.Find(n.s.form.HrefSelector).Attr("href")

		}

		//列表选择器不为空时
		if len(n.s.form.ListFields) > 0 {

			t, err := s.Html()

			if err != nil {

				n.s.notice.Error(err.Error())

				return

			}

			//解析列表选择器
			rows, err = n.s.form.ResolveSelector(t, n.s.form.ListFields, listUrl)

			if err != nil {

				n.s.notice.Error(err.Error())

				return
			}

		}

		//如果详情选择器为空就跳过
		if len(n.s.form.DetailFields) <= 0 {

			n.s.result.Push(rows)

			n.s.currentIndex++

			return

		}

		if strings.TrimSpace(href) == "" || isFind == false {

			n.s.notice.Error("a链接为空,当前链接为:", listUrl)

			return
		}

		n.s.detailCoroutineChan <- true

		n.s.detailWait.Add(1)

		go n.GetDetail(n.s.form.GetHref(href), rows)

	}).Size()

	if n.s.detailSize == 0 && size > 0 {

		n.s.detailSize = size

		//计算大概爬取总数量
		n.s.total = n.s.form.Length * size

	}

	if size <= 0 {

		n.s.notice.Error("a链接未发现")

	}

}

func (n normal) GetDetail(detailUrl string, rows *Rows) {

	defer func() {

		<-n.s.detailCoroutineChan

		n.s.detailWait.Done()

		n.s.currentIndex++

	}()

	select {

	case <-n.s.cxt.Done():

		return

	default:

	}

	html, err := n.s.form.GetHtml(detailUrl)

	if err != nil {

		//n.form.Notice.PushMessage(notice.NewError(err.Error()))

		//n.form.Notice.Error(err.Error())

		return

	}

	//中间链接（中间页面）
	if len(n.s.form.MiddleSelector) > 0 {

		for _, s := range n.s.form.MiddleSelector {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

			if err != nil {

				n.s.notice.Error(err.Error())

				return

			}

			href, b := doc.Find(s).Attr("href")

			if !b {

				return
			}

			href = n.s.form.GetHref(href)

			html, err = n.s.form.GetHtml(href)

			if err != nil {

				n.s.notice.Error(err.Error())

				return

			}

		}

	}

	res, err := n.s.form.ResolveSelector(html, n.s.form.DetailFields, detailUrl)

	if err != nil {

		n.s.notice.Error(err.Error())

		return
	}

	//合并列表结果
	rows.Add(res)

	n.s.result.Push(rows)

}

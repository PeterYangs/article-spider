package normal

import (
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/notice"
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

	//html, header, err := form.Client.Request().GetToStringWithHeader(listUrl)

	html, header, err := n.form.Client.Request().GetToStringWithHeader(listUrl)

	if err != nil {

		n.form.Notice.PushMessage(notice.NewError(err.Error()))

		return

	}

	//自动转码
	if n.form.DisableAutoCoding == false {

		html, err = n.form.DealCoding(html, header)

		if err != nil {

			n.form.Notice.PushMessage(notice.NewError(err.Error()))

			return

		}

	}

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {

		n.form.Notice.PushMessage(notice.NewError(err.Error()))

		return

	}

	//查找列表中的a链接
	size := doc.Find(n.form.ListSelector).Each(func(i int, s *goquery.Selection) {

		href := ""

		isFind := false

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

		//n.form.Notice.PushMessage(notice.NewInfo(href))

		n.GetDetail(n.form.GetHref(href))

	}).Size()

	if size <= 0 {

		n.form.Notice.PushMessage(notice.NewInfo("a链接未发现"))
	}

	n.form.Notice.Close()

}

func (n *normal) GetDetail(detailUrl string) {

	html, header, err := n.form.Client.Request().GetToStringWithHeader(detailUrl)

	if err != nil {

		n.form.Notice.PushMessage(notice.NewError(err.Error()))

		return

	}

	//自动转码
	if n.form.DisableAutoCoding == false {

		html, err = n.form.DealCoding(html, header)

		if err != nil {

			n.form.Notice.PushMessage(notice.NewError(err.Error()))

			return

		}

	}

	////goquery加载html
	//doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	//if err != nil {
	//
	//	n.form.Notice.PushMessage(notice.NewError(err.Error()))
	//
	//	return
	//
	//}

}

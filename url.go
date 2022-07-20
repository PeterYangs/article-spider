package article_spider

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type url struct {
	s *Spider
}

func NewUrl(s *Spider) *url {

	return &url{s: s}
}

func (u url) Start() {

	u.s.total = len(u.s.form.DetailUrls)

	for _, uu := range u.s.form.DetailUrls {

		u.s.detailCoroutineChan <- true

		u.s.detailWait.Add(1)

		go u.GetDetail(uu)

	}

	u.s.detailWait.Wait()

	u.s.cancel()

}

func (u url) GetDetail(detailUrl string) {

	defer func() {

		<-u.s.detailCoroutineChan

		u.s.detailWait.Done()

		u.s.currentIndex++

	}()

	select {

	case <-u.s.cxt.Done():

		return

	default:

	}

	html, err := u.s.form.GetHtml(detailUrl)

	if err != nil {

		//n.form.Notice.PushMessage(notice.NewError(err.Error()))

		//n.form.Notice.Error(err.Error())

		return

	}

	//中间链接（中间页面）
	if len(u.s.form.MiddleSelector) > 0 {

		for _, s := range u.s.form.MiddleSelector {

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

			if err != nil {

				u.s.notice.Error(err.Error())

				return

			}

			href, b := doc.Find(s).Attr("href")

			if !b {

				return
			}

			href = u.s.form.GetHref(href)

			html, err = u.s.form.GetHtml(href)

			if err != nil {

				u.s.notice.Error(err.Error())

				return

			}

		}

	}

	res, err := u.s.form.ResolveSelector(html, u.s.form.DetailFields, detailUrl)

	if err != nil {

		u.s.notice.Error(err.Error())

		return
	}

	u.s.result.Push(res)

}
